package cli

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
)

var _ = strconv.Itoa(0)

func CmdCollectionByDenomId() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection-by-denom-id [denom-id]",
		Short: "Query collection-by-denom-id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqDenomId := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryCollectionByDenomIdRequest{
				DenomId: reqDenomId,
			}

			res, err := queryClient.CollectionByDenomId(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
