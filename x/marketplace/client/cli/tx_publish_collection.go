package cli

import (
	"fmt"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
)

func CmdPublishCollection() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "publish-collection [denom-id] [first-sale-royalties] [resale-royalties]",
		Short: "Publish collection for sale",
		Long:  "Publish collection for sale",
		Example: fmt.Sprintf(
			"$ %s tx marketplace publish-collection <denom-id> " +
				"--first-sale-royalties=<first-sale-royalties> " +
				"--resale-royalties=<resale-royalties> " +
				version.AppName,
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			denomID := args[0]

			firstSaleRoyalties, err := cmd.Flags().GetString(FlagFirstSaleRoyalties)
			if err != nil {
				return err
			}

			resaleRoyalties, err := cmd.Flags().GetString(FlagResaleRoyalties)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgPublishCollection(
				clientCtx.GetFromAddress().String(),
				denomID,
				firstSaleRoyalties,
				resaleRoyalties,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsPublishCollection)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
