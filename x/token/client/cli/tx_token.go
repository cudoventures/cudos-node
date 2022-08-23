package cli

import (
	"encoding/json"

	"github.com/CudoVentures/cudos-node/x/token/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdCreateToken() *cobra.Command {
	cmd := &cobra.Command{
		// todo add example + initial-balances and max-supply should be optional flags
		Use:   "create [denom] [name] [decimals] [initial-balances] [max-supply]",
		Short: "Create a new fungible token",
		Long: `
Parameters:
  denom:            Unique token symbol, eg. ETH
  name:             Token name, eg. Ethereum
  decimals:         Number of digits after the decimal point, eg. 1.00 has two decimals
  initial-balances: Mint tokens to specified addresses, eg. --initial-balances='[{"address":"cudos12mly447tcat35rs6ltzj8j8ez6ul8yv6dxh3u8","amount":"123"}]'
  max-supply:       Maximum supply that cannot be exceeded, eg. 21000000
Example:
  cudos-noded tx token create-token TOK Tokemania 2 '[{"address":"cudos1wfmq7vdd7sw6v648z0ts5ftfl8ptp9hru9sjsm","amount":"123"}]' 100000000 --from=cudos1j34vv7gmvl6cqg90my8mgzqpfgak6xzrg5l6p7 --chain-id=cudos-network --gas=auto --fees=100000000000acudos 
`,
		Args: cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			denom := args[0]
			name := args[1]

			decimals, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}

			initialBalances := []*types.Balance{}
			err = json.Unmarshal([]byte(args[3]), &initialBalances)
			if err != nil {
				return err
			}

			maxSupply, err := cast.ToUint64E(args[4])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateToken(
				clientCtx.GetFromAddress().String(),
				denom,
				name,
				decimals,
				initialBalances,
				maxSupply,
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
			argInitialBalances := new(types.Balance)
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
