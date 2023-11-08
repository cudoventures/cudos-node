package cli

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
)

var _ = strconv.Itoa(0)

func CmdUpdatePrice() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-price [id] [price]",
		Short: "Update NFT price",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argId := args[0]

			nftID, err := strconv.ParseUint(argId, 10, 64)
			if err != nil {
				return err
			}

			argPrice := args[1]

			price, err := sdk.ParseCoinNormalized(argPrice)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdatePrice(
				clientCtx.GetFromAddress().String(),
				nftID,
				price,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
