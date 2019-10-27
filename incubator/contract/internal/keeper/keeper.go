package keeper

import (
	"encoding/binary"
	"encoding/json"
	"os"
	"path/filepath"

	wasm "github.com/confio/go-cosmwasm"
	wasmTypes "github.com/confio/go-cosmwasm/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmwasm/modules/incubator/contract/internal/types"
)

// Keeper will have a reference to Wasmer with it's own data directory.
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
	accountKeeper auth.AccountKeeper
	bankKeeper    bank.Keeper

	wasmer wasm.Wasmer
}

// NewKeeper creates a new contract Keeper instance
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, accountKeeper auth.AccountKeeper, bankKeeper bank.Keeper) Keeper {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	wasmer, err := wasm.NewWasmer(filepath.Join(home, ".wasmer"), 3)
	if err != nil {
		panic(err)
	}

	return Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		wasmer:        *wasmer,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
	}
}

// Create uploads and compiles a WASM contract, returning a short identifier for the contract
func (k Keeper) Create(ctx sdk.Context, creator sdk.AccAddress, wasmCode []byte) (contractID uint64, sdkErr sdk.Error) {
	codeID, err := k.wasmer.Create(wasmCode)
	if err != nil {
		return contractID, types.ErrCreateFailed(err)
	}

	store := ctx.KVStore(k.storeKey)
	contractID = k.autoIncrementID(ctx, types.KeyNextContractID)
	contractInfo := types.NewContractInfo(codeID, creator)
	// 0x01 | ContractID (uint64) -> ContractInfo
	store.Set(types.GetCodeKey(contractID), k.cdc.MustMarshalBinaryLengthPrefixed(contractInfo))

	return contractID, nil
}

// Instantiate creates an instance of a WASM contract
func (k Keeper) Instantiate(ctx sdk.Context, creator sdk.AccAddress, contractID uint64, initMsg interface{}, deposit sdk.Coins) (sdk.AccAddress, sdk.Error) {
	// create contract address
	contractAddress := addrFromUint64(contractID)
	existingAccnt := k.accountKeeper.GetAccount(ctx, contractAddress)
	if existingAccnt != nil {
		return nil, types.ErrAccountExists(existingAccnt.GetAddress())
	}

	// deposit initial contract funds
	contractAccount := k.accountKeeper.NewAccountWithAddress(ctx, contractAddress)
	contractAccount.SetCoins(deposit)
	k.accountKeeper.SetAccount(ctx, contractAccount)

	// get contact info
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetCodeKey(contractID))
	var contractInfo types.ContractInfo
	if bz != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &contractInfo)
	}

	// prepare params for contract instantiate call
	params := types.NewInstanceParams(ctx, creator, deposit, contractAccount)
	initMsgBz, err := json.Marshal(initMsg)
	if err != nil {
		return nil, sdk.ErrUnknownRequest("error encoding init message")
	}

	// create prefixed data store
	// 0x03 | contractAddress (sdk.AccAddress)
	prefixStoreKey := types.GetInstanceStorePrefixKey(contractAddress)
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)

	// instantiate wasm contract
	_, err = k.wasmer.Instantiate(contractInfo.CodeID, params, initMsgBz, prefixStore, 100000000)
	if err != nil {
		return contractAddress, types.ErrInstantiateFailed(err)
	}

	// persist instance
	instance := types.NewInstance(contractID, creator, initMsgBz, prefixStore)
	// 0x02 | contractAddress (sdk.AccAddress) -> Instance
	store.Set(types.GetContractAddressKey(contractAddress), k.cdc.MustMarshalBinaryLengthPrefixed(instance))

	return contractAddress, nil
}

// Execute executes the contract instance (STUB)
func (k Keeper) Execute(ctx sdk.Context, contractAddress sdk.AccAddress, params wasmTypes.Params, msgs interface{}) sdk.Result {
	// get contractID, store from contractAddress
	// get codeID from contractID
	// res, err := k.wasmer.Execute(codeID, params, msgs, store, gasLimit)
	return sdk.Result{}
}

func (k Keeper) autoIncrementID(ctx sdk.Context, nextIDKey []byte) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(nextIDKey)
	var id uint64 = 0
	if bz != nil {
		k.cdc.MustUnmarshalBinaryBare(bz, &id)
	}
	bz = k.cdc.MustMarshalBinaryBare(id + 1)
	store.Set(nextIDKey, bz)
	return id
}

func addrFromUint64(id uint64) sdk.AccAddress {
	addr := make([]byte, 20)
	addr[0] = 'C'
	binary.PutUvarint(addr[1:], id)
	return addr
}
