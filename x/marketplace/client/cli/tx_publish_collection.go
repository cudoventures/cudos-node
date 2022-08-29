package cli

import (
	"fmt"
	"strings"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
)

func CmdPublishCollection() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "publish-collection [denom-id] [mint-royalties] [resale-royalties]",
		Short: "Publish collection for sale",
		Long:  "Publish collection for sale",
		Example: fmt.Sprintf(
			"$ %s tx marketplace publish-collection <denom-id> "+
				"--mint-royalties=<mint-royalties> "+
				"--resale-royalties=<resale-royalties> ",
			version.AppName,
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			denomID := args[0]

			flagMintRoyalties, err := cmd.Flags().GetString(FlagMintRoyalties)
			if err != nil {
				return err
			}

			mintRoyalties, err := parseRoyalties(flagMintRoyalties)
			if err != nil {
				return err
			}

			flagResaleRoyalties, err := cmd.Flags().GetString(FlagResaleRoyalties)
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

			msg := types.NewMsgPublishCollection(
				clientCtx.GetFromAddress().String(),
				denomID,
				mintRoyalties,
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

func parseRoyalties(royaltiesStr string) ([]types.Royalty, error) {
	royaltiesStr = strings.TrimSpace(royaltiesStr)
	royaltiesStrList := strings.Split(royaltiesStr, ",")

	var royalties []types.Royalty

	for _, royaltyStr := range royaltiesStrList {
		split := strings.Split(strings.TrimSpace(royaltyStr), ":")

		address, err := sdk.AccAddressFromBech32(strings.TrimSpace(split[0]))
		if err != nil {
			return nil, err
		}

		percent, err := sdk.NewDecFromStr(strings.TrimSpace(split[1]))
		if err != nil {
			return nil, err
		}

		royalty := types.Royalty{
			Address: address.String(),
			Percent: percent,
		}

		royalties = append(royalties, royalty)
	}

	return royalties, nil
}
