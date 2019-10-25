package keeper

import (
	"encoding/binary"
	"os"
	"path/filepath"

	wasm "github.com/confio/go-cosmwasm"
	wasmTypes "github.com/confio/go-cosmwasm/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmwasm/modules/incubator/contract/internal/types"
	"github.com/davecgh/go-spew/spew"
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

// Create uploads and compiles a WASM contract, returning a contract id
func (k Keeper) Create(ctx sdk.Context, creator sdk.AccAddress, wasmCode []byte, deposit sdk.Coins) (sdk.AccAddress, sdk.Error) {
	// TODO: implement wasm hash function to check code before trying to create?

	codeID, err := k.wasmer.Create(wasmCode)
	if err != nil {
		spew.Dump(err)
		return nil, types.ErrCreateFailed()
	}

	spew.Dump(codeID)

	// Create a contract address
	contractAddr := k.newContractAddress(ctx)
	existingAcc := k.accountKeeper.GetAccount(ctx, contractAddr)
	if existingAcc != nil {
		return nil, types.ErrAccountExists(existingAcc.GetAddress())
	}

	// Deposit initial contract funds
	k.accountKeeper.SetAccount(ctx, &auth.BaseAccount{Address: contractAddr})
	sdkErr := k.bankKeeper.SendCoins(ctx, creator, contractAddr, deposit)
	if sdkErr != nil {
		return nil, sdkErr
	}

	// Store n/[ContractAddress] -> [CodeID]
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyContractCode(contractAddr), k.cdc.MustMarshalBinaryBare(codeID))

	// 32 bytes is overkill... or int
	// a lot of these...
	// instantiating should be cheap...
	// auto-increment
	// instance 7 can use code blob 2

	// wasm/int = hash this... -> code id, creator address, contract address.., prefix for the KVStore.. each contract needs, make sure overlap..
	// 8 byte big endian... for prefix
	// address based on hash of id

	// key name: wasm/id/creator?

	return contractAddr, nil
}

// Instantiate creates an instance of a WASM contract
// TODO: combine msgs and params into a struct
func (k Keeper) Instantiate(ctx sdk.Context, contractAddress sdk.AccAddress, msgs []byte) sdk.Result {
	params := wasmTypes.Params{
		Block: wasmTypes.BlockInfo{},
		Message: wasmTypes.MessageInfo{
			Signer: "signer",
			SentFunds: []wasmTypes.Coin{{
				Denom:  "ATOM",
				Amount: "100",
			}},
		},
		Contract: wasmTypes.ContractInfo{
			Address: "contract",
			Balance: []wasmTypes.Coin{{
				Denom:  "ATOM",
				Amount: "100",
			}},
		},
	}

	gasLimit := int64(100000000)
	// works per wasm instruction..
	// 1 gas < 1 nanosecond....
	// hardcode constants.. n = 100
	// if you know you have x gas in cosmos... contract = x * n

	store := ctx.KVStore(k.storeKey)
	codeIDBin := store.Get(types.KeyContractCode(contractAddress))
	var codeID []byte
	if codeIDBin != nil {
		k.cdc.MustUnmarshalBinaryBare(codeIDBin, &codeID)
	}

	res, err := k.wasmer.Instantiate(codeID, params, msgs, store, gasLimit)
	if err != nil {
		spew.Dump(err)
		return sdk.Result{
			Codespace: types.DefaultCodespace,
			Data:      nil,
			// GasUsed:   uint64(res.GasUsed),
		}
	}

	spew.Dump(res)

	return sdk.Result{
		Codespace: types.DefaultCodespace,
		Data:      []byte(res.Data),
		GasUsed:   uint64(res.GasUsed),
	}
}

// TODO: create using creator + code id + nonce?
func (k Keeper) newContractAddress(ctx sdk.Context) sdk.AccAddress {
	id := k.autoIncrementID(ctx, types.KeyNextContractID)
	return addrFromUint64(id)
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
