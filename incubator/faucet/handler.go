package faucet

import (
	"fmt"

	"github.com/cosmos/modules/incubator/faucet/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns a handler for "faucet" type messages.
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

// Handle a message to Mint
func handleMsgMint(ctx sdk.Context, keeper Keeper, msg types.MsgMint) (*sdk.Result, error) {

	keeper.Logger(ctx).Info("received mint message: %s", msg)
	err := keeper.MintAndSend(ctx, msg.Minter)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{}, nil // return
}
