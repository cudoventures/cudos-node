package cli

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/CudoVentures/cudos-node/x/token/types"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
)

const (
	flagPacketTimeoutTimestamp = "packet-timeout-timestamp"
	listSeparator              = ","
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdCreateToken())

	return cmd
}

func CmdCreateToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [denom] [name] [decimals] [initial-balances] [max-supply]",
		Short: "Create a new fungible token",
		Long: `
Parameters:
  denom:            Unique token symbol, eg. ETH
  name:             Token name, eg. Ethereum
  decimals:         Number of digits after the decimal point, eg. 1.00 has two decimals
  initial-balances: Mint tokens to specified addresses, eg. --initial-balances='[{"address":"cudos12mly447tcat35rs6ltzj8j8ez6ul8yv6dxh3u8","amount":123}]'
  max-supply:       Maximum supply that cannot be exceeded, eg. 21000000
Example:
  cudos-noded tx token create-token TOK Tokemania 2 '[{"address":"cudos1wfmq7vdd7sw6v648z0ts5ftfl8ptp9hru9sjsm","amount":123}]' 100000000 --from=cudos1j34vv7gmvl6cqg90my8mgzqpfgak6xzrg5l6p7 --chain-id=cudos-network --gas=auto --fees=100000000000acudos 
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
