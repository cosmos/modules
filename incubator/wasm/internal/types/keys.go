package types

import (
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

// nolint
var (
	KeyLastCodeID     = []byte("lastCodeId")
	KeyLastInstanceID = []byte("lastInstanceId")

	CodeKeyPrefix       = []byte{0x01}
	InstanceKeyPrefix   = []byte{0x02}
	InstanceStorePrefix = []byte{0x03}
)

// GetCodeKey constructs the key for retreiving the ID for the WASM code
func GetCodeKey(contractID uint64) []byte {
	contractIDBz := sdk.Uint64ToBigEndian(contractID)
	return append(CodeKeyPrefix, contractIDBz...)
}

// GetContractAddressKey returns the key for the WASM contract instance
func GetContractAddressKey(addr sdk.AccAddress) []byte {
	return append(InstanceKeyPrefix, addr...)
}

// GetInstanceStorePrefixKey returns the store prefix for the WASM contract instance
func GetInstanceStorePrefixKey(addr sdk.AccAddress) []byte {
	return append(InstanceStorePrefix, addr...)
}
