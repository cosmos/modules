package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/modules/incubator/poa/internal/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/staking/validators/{validatorAddr}",
		createValidatorHandlerFn(cliCtx),
	).Methods("POST")
	r.HandleFunc(
		"/staking/validators/{validatorAddr}",
		editValidatorHandlerFn(cliCtx),
	).Methods("PUT")
}

type (
	// CreateValidatorRequest defines the properties of a create validator request's body.
	CreateValidatorRequest struct {
		BaseReq          rest.BaseReq             `json:"base_req" yaml:"base_req"`
		ValidatorAddress sdk.ValAddress           `json:"validator_address" yaml:"validator_address"` // in bech32
		Pubkey           string                   `json:"pubkey" yaml:"pubkey"`
		Description      stakingtypes.Description `json:"description" yaml:"description"`
	}

	// EditValidatorRequest defines the properties of a edit validator request's body.
	EditValidatorRequest struct {
		BaseReq          rest.BaseReq             `json:"base_req" yaml:"base_req"`
		ValidatorAddress sdk.ValAddress           `json:"address" yaml:"address"`
		Description      stakingtypes.Description `json:"description" yaml:"description"`
	}
)

// createValidatorHandlerFn implements a create validator handler that is responsible
// for constructing a properly formatted create validator transaction for signing.

func createValidatorHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateValidatorRequest

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		pk, err := sdk.GetConsPubKeyBech32(req.Pubkey)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgCreateValidator(req.ValidatorAddress, pk, req.Description)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// editValidatorHandlerFn implements a edit validator handler that is responsible
// for constructing a properly formatted edit validator transaction for signing.
func editValidatorHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req EditValidatorRequest

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.NewMsgEditValidator(req.ValidatorAddress, req.Description)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
