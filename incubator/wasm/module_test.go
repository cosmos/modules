package wasm

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

func TestHandleCreate(t *testing.T) {
	cases := map[string]struct {
		msg     sdk.Msg
		isValid bool
	}{
		"empty": {
			msg:     MsgStoreCode{},
			isValid: false,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			data, cleanup := setupTest(t)
			defer cleanup()

			h := data.module.NewHandler()
			res := h(data.ctx, tc.msg)
			if !tc.isValid {
				require.False(t, res.IsOK(), "%#v", res)
				return
			}
			require.True(t, res.IsOK(), "%#v", res)
			assert.Equal(t, 1, 1)
		})
	}
}
