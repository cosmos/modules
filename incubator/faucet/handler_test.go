package faucet

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	emoji "github.com/tmdvs/Go-Emoji-Utils"

	"github.com/cosmos/modules/incubator/faucet/internal/types"
	"github.com/tendermint/tendermint/crypto"
)

func TestEmoji(t *testing.T) {

	moduleAcct := sdk.AccAddress(crypto.AddressHash([]byte("foobar")))
	denom := "🥵"
	msg := types.NewMsgMint(moduleAcct, moduleAcct, time.Now().Unix(), denom)
	err := msg.ValidateBasic()
	require.NoError(t, err)

	results := emoji.FindAll(msg.Denom)
	if len(results) != 1 {
		fmt.Println("results did not equal 1")
		require.True(t, false)
	}
	emo, ok := results[0].Match.(emoji.Emoji)
	if !ok {
		fmt.Println("Not correct interface for Emoji")
		require.True(t, false)
	}
	fmt.Println(emo.Value)
}
