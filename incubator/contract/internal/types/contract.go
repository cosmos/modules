package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// ContractInfo is data for the uploaded contract WASM code
type ContractInfo struct {
	CodeID  []byte         `json:"code_id"`
	Creator sdk.AccAddress `json:"creator"`
}

// NewContractInfo fills a new Contract struct
func NewContractInfo(codeID []byte, creator sdk.AccAddress) ContractInfo {
	return ContractInfo{
		CodeID:  codeID,
		Creator: creator,
	}
}
