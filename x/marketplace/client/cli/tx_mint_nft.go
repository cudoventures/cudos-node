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

func CmdMintNft() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint-nft [denom-id] [recipient] [price] [name] [uri] [data]",
		Short: "Mint NFT via marketplace",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argDenomId := args[0]
			argRecipient := args[1]
			argPrice := args[2]
			argName := args[3]

			price, err := sdk.ParseCoinNormalized(argPrice)
			if err != nil {
				return err
			}

			uri, err := cmd.Flags().GetString(FlagMintNftUri)
			if err != nil {
				return err
			}

			data, err := cmd.Flags().GetString(FlagMintNftData)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgMintNft(
				clientCtx.GetFromAddress().String(),
				argDenomId,
				argRecipient,
				argName,
				uri,
				data,
				price,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsMintNFT)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
