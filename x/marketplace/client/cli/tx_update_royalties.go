package cli

import (
	"fmt"
	"strconv"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdUpdateRoyalties() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-royalties [id] [mint-royalties] [resale-royalties]",
		Short: "Update collection royalties",
		Long:  "Update collection royalties",
		Args:  cobra.ExactArgs(1),
		Example: fmt.Sprintf(
			"$ %s tx marketplace update-royalties <collection-id> "+
				"--mint-royalties=<mint-royalties> "+
				"--resale-royalties=<resale-royalties> ",
			version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argId := args[0]

			collectionID, err := strconv.ParseUint(argId, 10, 64)
			if err != nil {
				return err
			}

			flagMintRoyalties, err := cmd.Flags().GetString(FlagUpdateMintRoyalties)
			if err != nil {
				return err
			}

			mintRoyalties, err := parseRoyalties(flagMintRoyalties)
			if err != nil {
				return err
			}

			flagResaleRoyalties, err := cmd.Flags().GetString(FlagUpdateResaleRoyalties)
			if err != nil {
				return err
			}

			resaleRoyalties, err := parseRoyalties(flagResaleRoyalties)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateRoyalties(
				clientCtx.GetFromAddress().String(),
				collectionID,
				mintRoyalties,
				resaleRoyalties,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsUpdateRoyalties)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
