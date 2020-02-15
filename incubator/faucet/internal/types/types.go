package types

import (
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Mining is a struct that contains all the metadata of a mint
type Mining struct {
	minter   sdk.AccAddress `json:"Minter"`
	lastTime time.Time      `json:"LastTime"`
	total    sdk.Coin       `json:"Total"`
}

// NewMining returns a new Mining
func NewMining(minter sdk.AccAddress, coin sdk.Coin) Mining {
	return Mining{
		minter:   minter,
		lastTime: time.Now(),
		total:    coin,
	}
}

// GetMinter get minter of mining
func (w Mining) GetMinter() sdk.AccAddress {
	return w.minter
}

// implement fmt.Stringer
func (w Mining) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Minter: %s, Time: %s, Total: %s`, w.minter, w.lastTime, w.total))
}
