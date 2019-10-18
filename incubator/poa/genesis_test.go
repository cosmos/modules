package poa_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"

	sdk "github.com/cosmos/cosmos-sdk/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/modules/incubator/poa"
	keep "github.com/cosmos/modules/incubator/poa/internal/keeper"
	"github.com/cosmos/modules/incubator/poa/internal/types"
)

func TestInitGenesis(t *testing.T) {
	ctx, accKeeper, keeper, supplyKeeper := keep.CreateTestInput(t, false)

	weight := sdk.NewInt(10)

	params := keeper.GetParams(ctx)
	validators := make([]types.Validator, 2)

	// initialize the validators
	validators[0].OperatorAddress = sdk.ValAddress(keep.Addrs[0])
	validators[0].ConsPubKey = keep.PKs[0]
	validators[0].Description = stakingtypes.NewDescription("hoop", "", "", "", "")
	validators[0].Status = sdk.Bonded
	validators[0].Weight = weight
	validators[1].OperatorAddress = sdk.ValAddress(keep.Addrs[1])
	validators[1].ConsPubKey = keep.PKs[1]
	validators[1].Description = stakingtypes.NewDescription("bloop", "", "", "", "")
	validators[1].Status = sdk.Bonded
	validators[1].Weight = weight

	genesisState := types.NewGenesisState(params, validators)
	vals := poa.InitGenesis(ctx, keeper, accKeeper, supplyKeeper, genesisState)

	actualGenesis := poa.ExportGenesis(ctx, keeper)
	require.Equal(t, genesisState.Params, actualGenesis.Params)
	require.EqualValues(t, keeper.GetAllValidators(ctx), actualGenesis.Validators)

	// now make sure the validators are bonded and intra-tx counters are correct
	resVal, found := keeper.GetValidator(ctx, sdk.ValAddress(keep.Addrs[0]))
	require.True(t, found)
	require.Equal(t, sdk.Bonded, resVal.Status)

	resVal, found = keeper.GetValidator(ctx, sdk.ValAddress(keep.Addrs[1]))
	require.True(t, found)
	require.Equal(t, sdk.Bonded, resVal.Status)

	abcivals := make([]abci.ValidatorUpdate, len(vals))
	for i, val := range validators {
		abcivals[i] = val.ABCIValidatorUpdate()
	}

	require.Equal(t, abcivals, vals)
}

func TestValidateGenesis(t *testing.T) {
	genValidators1 := make([]types.Validator, 1, 5)
	pk := ed25519.GenPrivKey().PubKey()
	genValidators1[0] = types.NewValidator(sdk.ValAddress(pk.Address()), pk, stakingtypes.NewDescription("", "", "", "", ""))

	tests := []struct {
		name    string
		mutate  func(*types.GenesisState)
		wantErr bool
	}{
		{"default", func(*types.GenesisState) {}, false},
		// validate genesis validators
		{"duplicate validator", func(data *types.GenesisState) {
			data.Validators = genValidators1
			data.Validators = append(data.Validators, genValidators1[0])
		}, true},
		{"jailed and bonded validator", func(data *types.GenesisState) {
			data.Validators = genValidators1
			data.Validators[0].Jailed = true
			data.Validators[0].Status = sdk.Bonded
		}, true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			genesisState := types.DefaultGenesisState()
			tt.mutate(&genesisState)
			if tt.wantErr {
				assert.Error(t, types.ValidateGenesis(genesisState))
			} else {
				assert.NoError(t, types.ValidateGenesis(genesisState))
			}
		})
	}
}
