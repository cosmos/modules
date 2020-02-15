package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
)

/* When a module wishes to interact with an other module it is good practice to define what it will use
// as an interface so the module can not use things that are not permitted. */

// SupplyKeeper is required for mining coin
type SupplyKeeper interface {
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(
		ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins,
	) error
	GetSupply(ctx sdk.Context) (supply exported.SupplyI)
}
