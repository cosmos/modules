package cli

import (
	"encoding/json"
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
	poatypes "github.com/cosmos/modules/incubator/poa/internal/types"
)

// SetGenTxsInAppGenesisState - sets the genesis transactions in the app genesis state
func SetGenTxsInAppGenesisState(cdc *codec.Codec, appGenesisState map[string]json.RawMessage,
	genTxs []authtypes.StdTx) (map[string]json.RawMessage, error) {

	genesisState := GetGenesisStateFromAppState(cdc, appGenesisState)
	// convert all the GenTxs to JSON
	genTxsBz := make([]json.RawMessage, 0, len(genTxs))
	for _, genTx := range genTxs {
		txBz, err := cdc.MarshalJSON(genTx)
		if err != nil {
			return appGenesisState, err
		}
		genTxsBz = append(genTxsBz, txBz)
	}

	genesisState.GenTxs = genTxsBz
	return SetGenesisStateInAppState(cdc, appGenesisState, genesisState), nil
}

// ValidateAccountInGenesis checks that the provided key has sufficient
// coins in the genesis accounts
func ValidateAccountInGenesis(appGenesisState map[string]json.RawMessage,
	genAccIterator types.GenesisAccountsIterator,
	key sdk.AccAddress, cdc *codec.Codec) error {

	accountIsInGenesis := false

	poaDataBz := appGenesisState[poatypes.ModuleName]
	var poaData poatypes.GenesisState
	cdc.MustUnmarshalJSON(poaDataBz, &poaData)

	genUtilDataBz := appGenesisState[poatypes.ModuleName]
	var genesisState GenesisState
	cdc.MustUnmarshalJSON(genUtilDataBz, &genesisState)

	if !accountIsInGenesis {
		return fmt.Errorf("account %s in not in the app_state.accounts array of genesis.json", key)
	}

	return nil
}

type deliverTxfn func(abci.RequestDeliverTx) abci.ResponseDeliverTx

// DeliverGenTxs - deliver a genesis transaction
func DeliverGenTxs(ctx sdk.Context, cdc *codec.Codec, genTxs []json.RawMessage,
	stakingKeeper types.StakingKeeper, deliverTx deliverTxfn) []abci.ValidatorUpdate {

	for _, genTx := range genTxs {
		var tx authtypes.StdTx
		cdc.MustUnmarshalJSON(genTx, &tx)
		bz := cdc.MustMarshalBinaryLengthPrefixed(tx)
		res := deliverTx(abci.RequestDeliverTx{Tx: bz})
		if !res.IsOK() {
			panic(res.Log)
		}
	}
	return stakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
}

// ---

// GenesisState defines the raw genesis transaction in JSON
type GenesisState struct {
	GenTxs []json.RawMessage `json:"gentxs" yaml:"gentxs"`
}

// GetGenesisStateFromAppState gets the genutil genesis state from the expected app state
func GetGenesisStateFromAppState(cdc *codec.Codec, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[types.ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[types.ModuleName], &genesisState)
	}
	return genesisState
}

// SetGenesisStateInAppState sets the genutil genesis state within the expected app state
func SetGenesisStateInAppState(cdc *codec.Codec,
	appState map[string]json.RawMessage, genesisState GenesisState) map[string]json.RawMessage {

	genesisStateBz := cdc.MustMarshalJSON(genesisState)
	appState[types.ModuleName] = genesisStateBz
	return appState
}
