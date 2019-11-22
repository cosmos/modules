package wasm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

type testData struct {
	module     module.AppModule
	ctx        sdk.Context
	acctKeeper auth.AccountKeeper
}

// returns a cleanup function, which must be defered on
func setupTest(t *testing.T) (testData, func()) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)

	ctx, acctKeeper, keeper := CreateTestInput(t, false, tempDir)
	data := testData{
		module:     NewAppModule(keeper),
		ctx:        ctx,
		acctKeeper: acctKeeper,
	}
	cleanup := func() { os.RemoveAll(tempDir) }
	return data, cleanup
}

func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := ed25519.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

func mustLoad(path string) []byte {
	bz, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return bz
}

var (
	key1, pub1, addr1 = keyPubAddr()
	testContract      = mustLoad("./internal/keeper/testdata/contract.wasm")
)

func TestHandleCreate(t *testing.T) {
	cases := map[string]struct {
		msg     sdk.Msg
		isValid bool
	}{
		"empty": {
			msg:     MsgStoreCode{},
			isValid: false,
		},
		"invalid wasm": {
			msg: MsgStoreCode{
				Sender:       addr1,
				WASMByteCode: []byte("foobar"),
			},
			isValid: false,
		},
		"valid wasm": {
			msg: MsgStoreCode{
				Sender:       addr1,
				WASMByteCode: testContract,
			},
			isValid: true,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			data, cleanup := setupTest(t)
			defer cleanup()

			h := data.module.NewHandler()
			q := data.module.NewQuerierHandler()

			res := h(data.ctx, tc.msg)
			if !tc.isValid {
				require.False(t, res.IsOK(), "%#v", res)
				assertCodeList(t, q, data.ctx, 0)
				assertCodeBytes(t, q, data.ctx, 1, nil)
				return
			}
			require.True(t, res.IsOK(), "%#v", res)
			assertCodeList(t, q, data.ctx, 1)
			assertCodeBytes(t, q, data.ctx, 1, testContract)
		})
	}
}

func assertCodeList(t *testing.T, q sdk.Querier, ctx sdk.Context, expectedNum int) {
	bz, sdkerr := q(ctx, []string{QueryListCode}, abci.RequestQuery{})
	require.NoError(t, sdkerr)

	if len(bz) == 0 {
		require.Equal(t, expectedNum, 0)
		return
	}

	var res []CodeInfo
	err := json.Unmarshal(bz, &res)
	require.NoError(t, err)

	assert.Equal(t, expectedNum, len(res))
}

type wasmCode struct {
	Code []byte `json:"code", yaml:"code"`
}

func assertCodeBytes(t *testing.T, q sdk.Querier, ctx sdk.Context, codeID uint64, expectedBytes []byte) {
	path := []string{QueryGetCode, fmt.Sprintf("%d", codeID)}
	bz, sdkerr := q(ctx, path, abci.RequestQuery{})
	require.NoError(t, sdkerr)

	if len(bz) == 0 {
		require.Equal(t, len(expectedBytes), 0)
		return
	}

	var res wasmCode
	err := json.Unmarshal(bz, &res)
	require.NoError(t, err)

	assert.Equal(t, expectedBytes, res.Code)
}
