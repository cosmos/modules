package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Codes for governance errors
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeCreatedFailed sdk.CodeType = 1
	CodeAccountExists sdk.CodeType = 2
)

// ErrCreateFailed error for wasm code that has already been uploaded or failed
func ErrCreateFailed() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeCreatedFailed, fmt.Sprintf("created wasm contract failed"))
}

// ErrAccountExists error for a contract account that already exists
func ErrAccountExists(addr sdk.AccAddress) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeAccountExists, fmt.Sprintf("contract account %s already exists", addr.String()))
}
