package cli

import (
	"context"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
)

// Query all collections listed for sale in the marketplace
func CmdListCollection() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-collections",
		Short: "list all Collections",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllCollectionRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.CollectionAll(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// Show details about collection listed for sale in the marketplace
func CmdShowCollection() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-collection [id]",
		Short: "shows a Collection",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			params := &types.QueryGetCollectionRequest{
				Id: id,
			}

			res, err := queryClient.Collection(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
