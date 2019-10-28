package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Codes for wasm contract errors
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeCreatedFailed     sdk.CodeType = 1
	CodeAccountExists     sdk.CodeType = 2
	CodeInstantiateFailed sdk.CodeType = 3
)

// ErrCreateFailed error for wasm code that has already been uploaded or failed
func ErrCreateFailed(err error) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeCreatedFailed, fmt.Sprintf("create wasm contract failed: %s", err.Error()))
}

// ErrAccountExists error for a contract account that already exists
func ErrAccountExists(addr sdk.AccAddress) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeAccountExists, fmt.Sprintf("contract account %s already exists", addr.String()))
}

// ErrInstantiateFailed error for rust instantiate contract failure
func ErrInstantiateFailed(err error) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInstantiateFailed, fmt.Sprintf("instantiate wasm contract failed: %s", err.Error()))
}
