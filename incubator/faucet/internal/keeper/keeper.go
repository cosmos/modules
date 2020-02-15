package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/modules/incubator/faucet/internal/types"
)

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	SupplyKeeper types.SupplyKeeper
	storeKey     sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc          *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the nameservice Keeper
func NewKeeper(supplyKeeper types.SupplyKeeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		SupplyKeeper: supplyKeeper,
		storeKey:     storeKey,
		cdc:          cdc,
	}
}

// IsPresent check if the name is present in the store or not
func (k Keeper) IsPresent(ctx sdk.Context, minter sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(minter))
}
