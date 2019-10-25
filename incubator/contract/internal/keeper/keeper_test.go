package keeper

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func TestNewKeeper(t *testing.T) {
	_, _, keeper := CreateTestInput(t, false)
	require.NotNil(t, keeper)
}

func TestCreate(t *testing.T) {
	// remove existing wasmer directory
	home, err := os.UserHomeDir()
	require.NoError(t, err)
	os.RemoveAll(path.Join(home, ".wasmer"))

	ctx, accKeeper, keeper := CreateTestInput(t, false)
	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, deposit)

	wasmCode, err := ioutil.ReadFile("./testdata/contract.wasm")
	require.NoError(t, err)

	contractAddr, err := keeper.Create(ctx, creator, wasmCode, deposit)
	require.NoError(t, err)
	require.Equal(t, contractAddr.Empty(), false)

	require.Fail(t, "failureMessage string")
}

func TestInstantiate(t *testing.T) {
	// remove existing wasmer directory
	home, err := os.UserHomeDir()
	require.NoError(t, err)
	os.RemoveAll(path.Join(home, ".wasmer"))

	ctx, accKeeper, keeper := CreateTestInput(t, false)
	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, deposit)

	wasmCode, err := ioutil.ReadFile("./testdata/contract.wasm")
	require.NoError(t, err)

	contractAddr, err := keeper.Create(ctx, creator, wasmCode, deposit)
	require.NoError(t, err)
	require.Equal(t, contractAddr.Empty(), false)

	msg := []byte(`{"verifier": "fred", "beneficiary": "bob"}`)
	res := keeper.Instantiate(ctx, contractAddr, msg)
	// require.NotNil(t, res.Data)
	require.Nil(t, res.Data)
}

func createFakeFundedAccount(ctx sdk.Context, am auth.AccountKeeper, coins sdk.Coins) sdk.AccAddress {
	_, _, addr := keyPubAddr()
	baseAcct := auth.NewBaseAccountWithAddress(addr)
	_ = baseAcct.SetCoins(coins)
	am.SetAccount(ctx, &baseAcct)

	return addr
}

func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := ed25519.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}
