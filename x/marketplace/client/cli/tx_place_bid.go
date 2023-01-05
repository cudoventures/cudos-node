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

func CmdPlaceBid() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "place-bid [auction-id] [amount]",
		Short: "Place a bid on an auction",
		// todo example
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			auctionId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return types.ErrInvalidAuctionId
			}

			amount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return types.ErrInvalidPrice
			}

			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgPlaceBid(ctx.GetFromAddress().String(), auctionId, amount)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
