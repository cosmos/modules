package keeper

import (
	"os"
	"path/filepath"

	wasm "github.com/confio/go-cosmwasm"
	"github.com/cosmos/cosmos-sdk/codec"
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
