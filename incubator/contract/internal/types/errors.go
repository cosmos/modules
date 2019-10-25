package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Codes for wasm contract errors
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeCreatedFailed sdk.CodeType = 1
)

// ErrCreateFailed error for wasm code that has already been uploaded or failed
func ErrCreateFailed(err error) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeCreatedFailed, fmt.Sprintf("create wasm contract failed: %s", err.Error()))
}
