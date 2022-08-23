package cli

import (
	"encoding/json"
	"errors"

	"github.com/CudoVentures/cudos-node/x/token/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdCreateToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-token [denom] [name] [decimals] [initial-balances] [max-supply] [allowances]",
		Short: "Create a new token",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get indexes
			indexDenom := args[0]

			// Get value arguments
			argName := args[1]
			argDecimals, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}

			// todo currently it's single initialBalance, figure out array as single args
			argInitialBalances := make([]*types.InitialBalance, 0)
			err = json.Unmarshal([]byte(args[3]), &argInitialBalances)
			if err != nil {
				return err
			}

			argMaxSupply, ok := sdk.NewIntFromString(args[4])
			if !ok {
				return errors.New("invalid max supply")
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateToken(
				clientCtx.GetFromAddress().String(),
				indexDenom,
				argName,
				argDecimals,
				argInitialBalances,
				&argMaxSupply,
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

func CmdUpdateToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-token [denom] [name] [decimals] [initial-balances] [max-supply] [allowances]",
		Short: "Update a token",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get indexes
			indexDenom := args[0]

			// Get value arguments
			argName := args[1]
			argDecimals, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}
			argInitialBalances := new(types.InitialBalance)
			err = json.Unmarshal([]byte(args[3]), argInitialBalances)
			if err != nil {
				return err
			}
			argMaxSupply := args[4]
			argAllowances := new(types.Allowances)
			err = json.Unmarshal([]byte(args[5]), argAllowances)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateToken(
				clientCtx.GetFromAddress().String(),
				indexDenom,
				argName,
				argDecimals,
				argInitialBalances,
				argMaxSupply,
				argAllowances,
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

func CmdDeleteToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-token [denom]",
		Short: "Delete a token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			indexDenom := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteToken(
				clientCtx.GetFromAddress().String(),
				indexDenom,
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
