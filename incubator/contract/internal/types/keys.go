package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the contract module
	ModuleName = "contract"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// TStoreKey is the string transient store representation
	TStoreKey = "transient_" + ModuleName

	// QuerierRoute is the querier route for the staking module
	QuerierRoute = ModuleName

	// RouterKey is the msg router key for the staking module
	RouterKey = ModuleName
)

var (
	// KeyNextCodeID     = []byte("nextCodeId")
	KeyNextContractID = []byte("nextContractId")
)

// KeyContractCode returns the key for retrieving a contract code by address
func KeyContractCode(id sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("n/%x", id))
}
