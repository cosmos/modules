package keeper

import (
	"encoding/binary"
	"path/filepath"

	wasm "github.com/confio/go-cosmwasm"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmwasm/modules/incubator/contract/internal/types"
	"github.com/tendermint/tendermint/crypto"
)

// GasMultiplier is how many cosmwasm gas points = 1 sdk gas point
// SDK reference costs can be found here: https://github.com/cosmos/cosmos-sdk/blob/02c6c9fafd58da88550ab4d7d494724a477c8a68/store/types/gas.go#L153-L164
// A write at ~3000 gas and ~200us = 10 gas per us (microsecond) cpu/io
// Rough timing have 88k gas at 90us, which is equal to 1k sdk gas... (one read)
const GasMultiplier = 100

// MaxGas for a contract is 900 million (enforced in rust)
const MaxGas = 900_000_000

// Keeper will have a reference to Wasmer with it's own data directory.
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
	accountKeeper auth.AccountKeeper
	bankKeeper    bank.Keeper

	wasmer wasm.Wasmer
}

// NewKeeper creates a new contract Keeper instance
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, accountKeeper auth.AccountKeeper, bankKeeper bank.Keeper, homeDir string) Keeper {
	wasmer, err := wasm.NewWasmer(filepath.Join(homeDir, "wasm"), 3)
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
func (k Keeper) Create(ctx sdk.Context, creator sdk.AccAddress, wasmCode []byte) (codeID uint64, sdkErr sdk.Error) {
	codeHash, err := k.wasmer.Create(wasmCode)
	if err != nil {
		return 0, types.ErrCreateFailed(err)
	}

	store := ctx.KVStore(k.storeKey)
	codeID = k.autoIncrementID(ctx, types.KeyLastCodeID)
	contractInfo := types.NewCodeInfo(codeHash, creator)
	// 0x01 | codeID (uint64) -> ContractInfo
	store.Set(types.GetCodeKey(codeID), k.cdc.MustMarshalBinaryBare(contractInfo))

	return codeID, nil
}

// Instantiate creates an instance of a WASM contract
func (k Keeper) Instantiate(ctx sdk.Context, creator sdk.AccAddress, codeID uint64, initMsg []byte, deposit sdk.Coins) (sdk.AccAddress, sdk.Error) {
	// create contract address
	contractAddress := k.generateContractAddress(ctx, codeID)
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
	bz := store.Get(types.GetCodeKey(codeID))
	var codeInfo types.CodeInfo
	if bz != nil {
		k.cdc.MustUnmarshalBinaryBare(bz, &codeInfo)
	}

	// prepare params for contract instantiate call
	params := types.NewParams(ctx, creator, deposit, contractAccount)

	// create prefixed data store
	// 0x03 | contractAddress (sdk.AccAddress)
	prefixStoreKey := types.GetContractStorePrefixKey(contractAddress)
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)

	// instantiate wasm contract
	gas := gasForContract(ctx)
	res, err := k.wasmer.Instantiate(codeInfo.CodeHash, params, initMsg, prefixStore, gas)
	if err != nil {
		return contractAddress, types.ErrInstantiateFailed(err)
	}
	consumeGas(ctx, res.GasUsed)

	// persist instance
	instance := types.NewContract(codeID, creator, initMsg, prefixStore)
	// 0x02 | contractAddress (sdk.AccAddress) -> Instance
	store.Set(types.GetContractAddressKey(contractAddress), k.cdc.MustMarshalBinaryBare(instance))

	return contractAddress, nil
}

// Execute executes the contract instance
func (k Keeper) Execute(ctx sdk.Context, contractAddress sdk.AccAddress, creator sdk.AccAddress, coins sdk.Coins, msgs []byte) (sdk.Result, sdk.Error) {
	store := ctx.KVStore(k.storeKey)

	var contract types.Contract
	contractBz := store.Get(types.GetContractAddressKey(contractAddress))
	if contractBz != nil {
		k.cdc.MustUnmarshalBinaryBare(contractBz, &contract)
	}

	var codeInfo types.CodeInfo
	contractInfoBz := store.Get(types.GetCodeKey(contract.CodeID))
	if contractInfoBz != nil {
		k.cdc.MustUnmarshalBinaryBare(contractInfoBz, &codeInfo)
	}

	contractAccount := k.accountKeeper.GetAccount(ctx, contractAddress)
	params := types.NewParams(ctx, creator, coins, contractAccount)

	prefixStoreKey := types.GetContractStorePrefixKey(contractAddress)
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)

	gas := gasForContract(ctx)
	res, err := k.wasmer.Execute(codeInfo.CodeHash, params, msgs, prefixStore, gas)
	if err != nil {
		return sdk.Result{}, types.ErrExecuteFailed(err)
	}
	consumeGas(ctx, res.GasUsed)

	// TODO: this needs to dispatch all the messages returned from the Execute function
	// this is how we can send the tokens out of the contract

	return types.CosmosResult(*res), nil
}

func gasForContract(ctx sdk.Context) uint64 {
	meter := ctx.GasMeter()
	remaining := (meter.Limit() - meter.GasConsumed()) * GasMultiplier
	if remaining > MaxGas {
		return MaxGas
	}
	return remaining
}

func consumeGas(ctx sdk.Context, gas uint64) {
	consumed := gas / GasMultiplier
	ctx.GasMeter().ConsumeGas(consumed, "wasm contract")
}

// generates a contract address from codeID + instanceID
func (k Keeper) generateContractAddress(ctx sdk.Context, codeID uint64) sdk.AccAddress {
	instanceID := k.autoIncrementID(ctx, types.KeyLastInstanceID)
	// NOTE: It is possible to get a duplicate address if either codeID or instanceID
	// overflow 32 bits. This is highly improbable, but something that could be refactored.
	contractID := codeID<<32 + instanceID
	return addrFromUint64(contractID)
}

func (k Keeper) autoIncrementID(ctx sdk.Context, lastIDKey []byte) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(lastIDKey)
	id := uint64(1)
	if bz != nil {
		id = binary.BigEndian.Uint64(bz)
	}
	bz = sdk.Uint64ToBigEndian(id + 1)
	store.Set(lastIDKey, bz)
	return id
}

func addrFromUint64(id uint64) sdk.AccAddress {
	addr := make([]byte, 20)
	addr[0] = 'C'
	binary.PutUvarint(addr[1:], id)
	return sdk.AccAddress(crypto.AddressHash(addr))
}
