package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	// ErrWithdrawTooOften withdraw too often
	ErrWithdrawTooOften  = sdkerrors.Register(ModuleName, 100, "Each address can withdraw only once")
	ErrFaucetKeyEmpty    = sdkerrors.Register(ModuleName, 101, "Armor should Not be empty.")
	ErrFaucetKeyExisted  = sdkerrors.Register(ModuleName, 102, "Faucet key existed")
	ErrCantWithdrawStake = sdkerrors.Register(ModuleName, 103, "Can't withdraw staking token")
	ErrNoEmoji           = sdkerrors.Register(ModuleName, 104, "No emoji present in denom")
)
