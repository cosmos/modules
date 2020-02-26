package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	// ErrWithdrawTooOften withdraw too often
	ErrWithdrawTooOften = sdkerrors.Register(ModuleName, 1, "You can only withdraw once in 24 hours")
)
