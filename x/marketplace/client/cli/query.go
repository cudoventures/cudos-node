package cli

import (
	"fmt"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group marketplace queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryParams())
	cmd.AddCommand(CmdListCollection())
	cmd.AddCommand(CmdShowCollection())
	cmd.AddCommand(CmdListNft())
	cmd.AddCommand(CmdShowNft())
	cmd.AddCommand(CmdCollectionByDenomId())

	cmd.AddCommand(CmdListAdmins())

	// this line is used by starport scaffolding # 1

	return cmd
}
