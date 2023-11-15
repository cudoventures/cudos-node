package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/CudoVentures/cudos-node/x/addressbook/types"
)

func CmdListAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-addresses",
		Short: "list all addresses",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllAddressRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.AddressAll(context.Background(), params)
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

func CmdShowAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-address [creator] [network] [label]",
		Short: "shows a address",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argCreator := args[0]
			argNetwork := args[1]
			argLabel := args[2]

			params := &types.QueryGetAddressRequest{
				Creator: argCreator,
				Network: argNetwork,
				Label:   argLabel,
			}

			res, err := queryClient.Address(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
