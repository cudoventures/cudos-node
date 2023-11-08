package cli

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
)

var _ = strconv.Itoa(0)

func CmdCreateCollection() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-collection [id] [mint-royalties] [resale-royalties] [verified]",
		Short: "Create denom and publish it into the marketplace",
		Long:  "Create denom and publish it into the marketplace",
		Example: "--name=<denom-name> " +
			"--symbol=<symbol-name> " +
			"--schema=<schema-content or path to schema.json> " +
			"--traits=<traits>" +
			"--description=<description> " +
			"--minter=<minter> " +
			"--data=<data> " +
			"--mint-royalties=<mint-royalties> " +
			"--resale-royalties=<resale-royalties> " +
			"--verified=<verified> ",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argDenomId := args[0]

			flagDenomName, err := cmd.Flags().GetString(FlagCreateCollectionName)
			if err != nil {
				return err
			}

			flagSchema, err := cmd.Flags().GetString(FlagCreateCollectionSchema)
			if err != nil {
				return err
			}

			flagSymbol, err := cmd.Flags().GetString(FlagCreateCollectionSymbol)
			if err != nil {
				return err
			}

			flagTraits, err := cmd.Flags().GetString(FlagCreateCollectionTraits)
			if err != nil {
				return err
			}

			flagDescription, err := cmd.Flags().GetString(FlagCreateCollectionDescription)
			if err != nil {
				return err
			}

			flagMinter, err := cmd.Flags().GetString(FlagCreateCollectionMinter)
			if err != nil {
				return err
			}

			flagData, err := cmd.Flags().GetString(FlagCreateCollectionData)
			if err != nil {
				return err
			}

			flagMintRoyalties, err := cmd.Flags().GetString(FlagCreateCollectionMintRoyalties)
			if err != nil {
				return err
			}

			mintRoyalties, err := parseRoyalties(flagMintRoyalties)
			if err != nil {
				return err
			}

			flagResaleRoyalties, err := cmd.Flags().GetString(FlagCreateCollectionResaleRoyalties)
			if err != nil {
				return err
			}

			resaleRoyalties, err := parseRoyalties(flagResaleRoyalties)
			if err != nil {
				return err
			}

			verified := false

			flagVerified, err := cmd.Flags().GetString(FlagCreateCollectionVerified)
			if err != nil {
				return err
			}

			if flagVerified != "" {
				verified, err = strconv.ParseBool(flagVerified)
				if err != nil {
					return err
				}
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateCollection(
				clientCtx.GetFromAddress().String(),
				argDenomId,
				flagDenomName,
				flagSchema,
				flagSymbol,
				flagTraits,
				flagDescription,
				flagMinter,
				flagData,
				mintRoyalties,
				resaleRoyalties,
				verified,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsCreateCollection)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
