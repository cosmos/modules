package keeper

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func TestNewKeeper(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	_, _, keeper := CreateTestInput(t, false, tempDir)
	require.NotNil(t, keeper)
}

func TestCreate(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	ctx, accKeeper, keeper := CreateTestInput(t, false, tempDir)

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, deposit)

	wasmCode, err := ioutil.ReadFile("./testdata/contract.wasm")
	require.NoError(t, err)

	contractID, err := keeper.Create(ctx, creator, wasmCode)
	require.NoError(t, err)
	require.Equal(t, uint64(1), contractID)
}

func TestInstantiate(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	ctx, accKeeper, keeper := CreateTestInput(t, false, tempDir)

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, deposit)

	wasmCode, err := ioutil.ReadFile("./testdata/contract.wasm")
	require.NoError(t, err)

	contractID, err := keeper.Create(ctx, creator, wasmCode)
	require.NoError(t, err)

	initMsg := InitMsg{
		Verifier:    "fred",
		Beneficiary: "bob",
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	gasBefore := ctx.GasMeter().GasConsumed()

	addr, err := keeper.Instantiate(ctx, creator, contractID, initMsgBz, deposit)
	require.NoError(t, err)
	require.Equal(t, "cosmos18vd8fpwxzck93qlwghaj6arh4p7c5n89uzcee5", addr.String())

	gasAfter := ctx.GasMeter().GasConsumed()
	kvStoreGas := uint64(20588) // calculated by disabling contract gas reduction and running test
	require.Equal(t, kvStoreGas+433, gasAfter-gasBefore)
}

func TestExecute(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	ctx, accKeeper, keeper := CreateTestInput(t, false, tempDir)

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, deposit)

	wasmCode, err := ioutil.ReadFile("./testdata/contract.wasm")
	require.NoError(t, err)

	contractID, err := keeper.Create(ctx, creator, wasmCode)
	require.NoError(t, err)

	_, _, fred := keyPubAddr()
	_, _, bob := keyPubAddr()
	initMsg := InitMsg{
		Verifier:    fred.String(),
		Beneficiary: bob.String(),
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	addr, err := keeper.Instantiate(ctx, creator, contractID, initMsgBz, deposit)
	require.NoError(t, err)
	require.Equal(t, "cosmos18vd8fpwxzck93qlwghaj6arh4p7c5n89uzcee5", addr.String())

	// unauthorized
	res, err := keeper.Execute(ctx, addr, creator, deposit, []byte(`{}`))
	require.Error(t, err)
	require.Contains(t, err.Error(), "Unauthorized")

	// verifier can execute, and get proper gas amount
	start := time.Now()
	gasBefore := ctx.GasMeter().GasConsumed()

	res, err = keeper.Execute(ctx, addr, fred, deposit, []byte(`{}`))
	diff := time.Now().Sub(start)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, uint64(81488), res.GasUsed)

	// make sure gas is properly deducted from ctx
	gasAfter := ctx.GasMeter().GasConsumed()
	kvStoreGas := uint64(5278) // calculated by disabling contract gas reduction and running test
	require.Equal(t, kvStoreGas+814, gasAfter-gasBefore)

	t.Logf("Duration: %v (81488 gas)\n", diff)
}

type InitMsg struct {
	Verifier    string `json:"verifier"`
	Beneficiary string `json:"beneficiary"`
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
