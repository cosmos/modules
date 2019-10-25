package types

import (
	wasmTypes "github.com/confio/go-cosmwasm/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

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

// Instance stores a WASM contract instance
type Instance struct {
	ContractID  uint64       `json:"contract_id"`
	PrefixStore prefix.Store `json:"prefix_store"`
}

// NewInstanceParams initializes params for a contract instance
func NewInstanceParams(ctx sdk.Context, creator sdk.AccAddress, deposit sdk.Coins, contractAcct auth.Account) wasmTypes.Params {
	return wasmTypes.Params{
		Block: wasmTypes.BlockInfo{
			Height:  ctx.BlockHeight(),
			Time:    ctx.BlockTime().Unix(),
			ChainID: ctx.ChainID(),
		},
		Message: wasmTypes.MessageInfo{
			Signer:    creator.String(),
			SentFunds: NewWasmCoins(deposit),
		},
		Contract: wasmTypes.ContractInfo{
			Address: contractAcct.GetAddress().String(),
			Balance: NewWasmCoins(contractAcct.GetCoins()),
		},
	}
}

// NewWasmCoins translates between Cosmos SDK coins and Wasm coins
func NewWasmCoins(cosmosCoins sdk.Coins) (wasmCoins []wasmTypes.Coin) {
	for _, coin := range cosmosCoins {
		wasmCoin := wasmTypes.Coin{
			Denom:  coin.Denom,
			Amount: coin.Amount.String(),
		}
		wasmCoins = append(wasmCoins, wasmCoin)
	}
	return wasmCoins
}

// NewInstance creates a new instance of a given WASM contract
func NewInstance(contractID uint64, prefixStore prefix.Store) Instance {
	return Instance{
		ContractID:  contractID,
		PrefixStore: prefixStore,
	}
}
