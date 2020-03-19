package cli

import (
	"bufio"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/modules/incubator/faucet/internal/types"
)

// GetTxCmd return faucet sub-command for tx
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	faucetTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "faucet transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	faucetTxCmd.AddCommand(flags.PostCommands(
		GetCmdMint(cdc),
		GetCmdMintFor(cdc),
	)...)

	return faucetTxCmd
}

// GetCmdWithdraw is the CLI command for mining coin
func GetCmdMint(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "mint",
		Short: "mint coin to sender address",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			msg := types.NewMsgMint(cliCtx.GetFromAddress(), cliCtx.GetFromAddress(), sdk.NewInt(time.Now().Unix()))
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdWithdraw is the CLI command for mining coin
func GetCmdMintFor(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "mintfor [address]",
		Short: "mint coin for new address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			address, _ := sdk.AccAddressFromBech32(args[0])

			msg := types.NewMsgMint(cliCtx.GetFromAddress(), address, sdk.NewInt(time.Now().Unix()))
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
