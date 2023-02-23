package cli

import (
	"github.com/CudoVentures/cudos-node/x/addressbook/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdCreateAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-address [network] [label] [value]",
		Short: "Create a new address",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			// Get indexes
			indexNetwork := args[0]
			indexLabel := args[1]

			// Get value arguments
			argValue := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateAddress(
				clientCtx.GetFromAddress().String(),
				indexNetwork,
				indexLabel,
				argValue,
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

func CmdUpdateAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-address [network] [label] [value]",
		Short: "Update a address",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get indexes
			indexNetwork := args[0]
			indexLabel := args[1]

			// Get value arguments
			argValue := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateAddress(
				clientCtx.GetFromAddress().String(),
				indexNetwork,
				indexLabel,
				argValue,
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

func CmdDeleteAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-address [network] [label]",
		Short: "Delete a address",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			indexNetwork := args[0]
			indexLabel := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteAddress(
				clientCtx.GetFromAddress().String(),
				indexNetwork,
				indexLabel,
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
