package cli

import (
	"fmt"
	"strconv"
	"time"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdPublishAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use: "publish-auction [denom-id] [token-id] [duration] [auction]",
		Example: fmt.Sprintf(
			`$ %s tx marketplace publish-auction "xyz" "1" "25h" "{"@type":"/cudoventures.cudosnode.marketplace.EnglishAuction","minPrice":{"denom":"acudos","amount":"1"}}"`,
			version.AppName,
		),
		Short: "List NFT for an auction",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			denomId := args[0]
			tokenId := args[1]

			duration, err := time.ParseDuration(args[2])
			if err != nil {
				return types.ErrInvalidAuctionDuration
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var a types.Auction
			err = clientCtx.Codec.UnmarshalInterfaceJSON([]byte(args[3]), &a)
			if err != nil {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "%s", err)
			}

			sender := clientCtx.GetFromAddress().String()
			msg, err := types.NewMsgPublishAuction(sender, denomId, tokenId, duration, a)
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
