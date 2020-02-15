package faucet

import (
	"fmt"

	"github.com/cosmos/modules/incubator/faucet/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns a handler for "nameservice" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgMint:
			return handleMsgMint(ctx, keeper, msg)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized faucet Msg type: %v", msg.Type()))
		}
	}
}

// Handle a message to set name
func handleMsgMint(ctx sdk.Context, keeper Keeper, msg types.MsgMint) (*sdk.Result, error) {

	fmt.Println("hand messages")
	// to implement later.
	// if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Name)) { // Checks if the the msg sender is the same as the current owner
	// 	return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner") // If not, throw an error
	// }
	// keeper.SetName(ctx, msg.Name, msg.Value) // If so, set the name to the value specified in the msg.
	keeper.IsPresent(ctx, msg.GetSigners()[0])
	return &sdk.Result{}, nil // return
}
