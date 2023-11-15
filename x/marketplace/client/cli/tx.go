package cli

import (
	"fmt"
	"time"

	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
)

var DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())

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

	cmd.AddCommand(CmdPublishCollection())
	cmd.AddCommand(CmdPublishNft())
	cmd.AddCommand(CmdBuyNft())
	cmd.AddCommand(CmdMintNft())
	cmd.AddCommand(CmdRemoveNft())
	cmd.AddCommand(CmdVerifyCollection())
	cmd.AddCommand(CmdUnverifyCollection())
	cmd.AddCommand(CmdCreateCollection())
	cmd.AddCommand(CmdUpdateRoyalties())
	cmd.AddCommand(CmdUpdatePrice())
	cmd.AddCommand(CmdAddAdmin())
	cmd.AddCommand(CmdRemoveAdmin())
	// this line is used by starport scaffolding # 1

	return cmd
}
