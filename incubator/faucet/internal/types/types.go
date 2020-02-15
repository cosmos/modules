package types

import (
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Mining is a struct that contains all the metadata of a mint
type Mining struct {
	Minter   sdk.AccAddress `json:"Minter"`
	LastTime time.Time      `json:"LastTime"`
	Total    sdk.Coin       `json:"Total"`
}

// NewMining returns a new Mining
func NewMining(minter sdk.AccAddress, coin sdk.Coin) Mining {
	return Mining{
		Minter:   minter,
		LastTime: time.Now(),
		Total:    coin,
	}
}

// GetMinter get minter of mining
func (w Mining) GetMinter() sdk.AccAddress {
	return w.Minter
}

// implement fmt.Stringer
func (w Mining) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Minter: %s, Time: %s, Total: %s`, w.Minter, w.LastTime, w.Total))
}
