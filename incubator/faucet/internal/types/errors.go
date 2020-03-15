package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	// ErrWithdrawTooOften withdraw too often
	ErrWithdrawTooOften = sdkerrors.Register(ModuleName, 100, "Each address can withdraw only once")
)
