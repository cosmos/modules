package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Mining is a struct that contains all the metadata of a mint
type Mining struct {
	Minter   sdk.AccAddress `json:"Minter"`
	LastTime int64          `json:"LastTime"`
	Tally    int64          `json:"Tally"`
}

// NewMining returns a new Mining
func NewMining(minter sdk.AccAddress, tally int64) Mining {
	return Mining{
		Minter:   minter,
		LastTime: 0,
		Tally:    tally,
	}
}

// GetMinter get minter of mining
func (w Mining) GetMinter() sdk.AccAddress {
	return w.Minter
}

// implement fmt.Stringer
func (w Mining) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Minter: %s, Time: %s, Tally: %s`, w.Minter, w.LastTime, w.Tally))
}

type FaucetKey struct {
	Armor string `json:" armor"`
}

// NewFaucetKey create a instance
func NewFaucetKey(armor string) FaucetKey {
	return FaucetKey{
		Armor: armor,
	}
}

// implement fmt.Stringer
func (f FaucetKey) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Armor: %s`, f.Armor))
}
