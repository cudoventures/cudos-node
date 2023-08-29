package cli

import (
	"strconv"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdPublishNft() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "publish-nft [token-id] [denom-id] [price]",
		Short: "Publish NFT for sale",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			tokenId := args[0]
			denomId := args[1]
			argPrice := args[2]

			price, err := sdk.ParseCoinNormalized(argPrice)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgPublishNft(
				clientCtx.GetFromAddress().String(),
				tokenId,
				denomId,
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
