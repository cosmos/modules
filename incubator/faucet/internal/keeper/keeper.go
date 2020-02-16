package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/modules/incubator/faucet/internal/types"
	"github.com/tendermint/tendermint/libs/log"
	"time"
)

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	SupplyKeeper  types.SupplyKeeper
	StakingKeeper types.StakingKeeper
	amount        int64         // set default amount for each mint.
	limit         time.Duration // rate limiting for mint, etc 24 * time.Hours
	storeKey      sdk.StoreKey  // Unexposed key to access store from sdk.Context
	cdc           *codec.Codec  // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the Faucet Keeper
func NewKeeper(
	supplyKeeper types.SupplyKeeper,
	stakingKeeper types.StakingKeeper,
	amount int64,
	rateLimit time.Duration,
	storeKey sdk.StoreKey,
	cdc *codec.Codec) Keeper {
	return Keeper{
		SupplyKeeper:  supplyKeeper,
		StakingKeeper: stakingKeeper,
		amount:        amount,
		limit:         rateLimit,
		storeKey:      storeKey,
		cdc:           cdc,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// MintAndSend mint coins and send to minter.
func (k Keeper) MintAndSend(ctx sdk.Context, minter sdk.AccAddress) error {

	mining := k.getMining(ctx, minter)

	// refuse mint in 24 hours
	if k.isPresent(ctx, minter) && mining.LastTime.Add(k.limit).After(time.Now()) {
		return types.ErrWithdrawTooOften
	}

	denom := k.StakingKeeper.BondDenom(ctx)
	newCoin := sdk.NewCoin(denom, sdk.NewInt(k.amount))
	mining.Total = mining.Total.Add(newCoin)
	k.setMining(ctx, minter, mining)

	k.Logger(ctx).Info("Mint coin: %s", newCoin)

	err := k.SupplyKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(newCoin))
	if err != nil {
		return err
	}
	err = k.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, minter, sdk.NewCoins(newCoin))
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) getMining(ctx sdk.Context, minter sdk.AccAddress) types.Mining {
	store := ctx.KVStore(k.storeKey)
	if !k.isPresent(ctx, minter) {
		denom := k.StakingKeeper.BondDenom(ctx)
		return types.NewMining(minter, sdk.NewCoin(denom, sdk.NewInt(0)))
	}
	bz := store.Get(minter.Bytes())
	var mining types.Mining
	k.cdc.MustUnmarshalBinaryBare(bz, &mining)
	return mining
}

func (k Keeper) setMining(ctx sdk.Context, minter sdk.AccAddress, mining types.Mining) {
	if mining.Minter.Empty() {
		return
	}
	if !mining.Total.IsPositive() {
		return
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(minter.Bytes(), k.cdc.MustMarshalBinaryBare(mining))
}

// IsPresent check if the name is present in the store or not
func (k Keeper) isPresent(ctx sdk.Context, minter sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(minter.Bytes())
}
