package cli

import (
	"strconv"
	"time"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdPublishAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use: "publish-auction [token-id] [denom-id] [duration] [auction-type]",
		// todo example
		Short: "List NFT for an auction",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			tokenId := args[0]
			denomId := args[1]

			duration, err := time.ParseDuration(args[2])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var auctionType types.AuctionType
			if err := clientCtx.Codec.UnmarshalInterfaceJSON([]byte(args[3]), &auctionType); err != nil {
				return err
			}

			msg, err := types.NewMsgPublishAuction(clientCtx.GetFromAddress().String(), denomId, tokenId, duration, auctionType)
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
