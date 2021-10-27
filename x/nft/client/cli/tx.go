package cli

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"cudos.org/cudos-node/x/nft/types"
)

// NewTxCmd returns the transaction commands for this module
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "NFT transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetCmdIssueDenom(),
		GetCmdMintNFT(),
		GetCmdEditNFT(),
		GetCmdTransferNft(),
		GetCmdBurnNFT(),
		GetCmdApproveNft(),
		GetCmdApproveAllNFT(),
		GetCmdRevokeNft(),
	)

	return txCmd
}

// GetCmdIssueDenom is the CLI command for an IssueDenom transaction
func GetCmdIssueDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "issue [denom-id]",
		Long: "Issue a new denom.",
		Example: fmt.Sprintf(
			"$ %s tx nft issue <denom-id> "+
				"--name=<denom-name> "+
				"--from=<key-name> "+
				"--schema=<schema-content or path to schema.json> "+
				"--chain-id=<chain-id> "+
				"--fees=<fee>",
			version.AppName,
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			denomName, err := cmd.Flags().GetString(FlagDenomName)
			if err != nil {
				return err
			}
			schema, err := cmd.Flags().GetString(FlagSchema)
			if err != nil {
				return err
			}
			optionsContent, err := ioutil.ReadFile(schema)
			if err == nil {
				schema = string(optionsContent)
			}

			msg := types.NewMsgIssueDenom(
				args[0],
				denomName,
				schema,
				clientCtx.GetFromAddress().String(),
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsIssueDenom)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdMintNFT is the CLI command for a MintNFT transaction
func GetCmdMintNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "mint [denom-id] [token-id]",
		Long: "Mint an NFT and set the owner to the recipient. Only the denom creator can mint a new NFT.",
		Example: fmt.Sprintf(
			"$ %s tx nft mint <denom-id> <token-id> "+
				"--recipient=<recipient> "+
				"--from=<key-name> "+
				"--uri=<uri> "+
				"--chain-id=<chain-id> "+
				"--fees=<fee>",
			version.AppName,
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var sender = clientCtx.GetFromAddress().String()

			recipient, err := cmd.Flags().GetString(FlagRecipient)
			if err != nil {
				return err
			}

			recipientStr := strings.TrimSpace(recipient)
			if len(recipientStr) > 0 {
				if _, err = sdk.AccAddressFromBech32(recipientStr); err != nil {
					return err
				}
			} else {
				recipient = sender
			}

			tokenName, err := cmd.Flags().GetString(FlagTokenName)
			if err != nil {
				return err
			}
			tokenURI, err := cmd.Flags().GetString(FlagTokenURI)
			if err != nil {
				return err
			}
			tokenData, err := cmd.Flags().GetString(FlagTokenData)
			if err != nil {
				return err
			}

			msg := types.NewMsgMintNFT(
				args[1],
				args[0],
				tokenName,
				tokenURI,
				tokenData,
				sender,
				recipient,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsMintNFT)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdEditNFT is the CLI command for sending an MsgEditNFT transaction
func GetCmdEditNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "edit [denom-id] [token-id]",
		Long: "Edit the token data of an NFT.",
		Example: fmt.Sprintf(
			"$ %s tx nft edit <denom-id> <token-id> "+
				"--from=<key-name> "+
				"--uri=<uri> "+
				"--chain-id=<chain-id> "+
				"--fees=<fee>",
			version.AppName,
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			tokenName, err := cmd.Flags().GetString(FlagTokenName)
			if err != nil {
				return err
			}
			tokenURI, err := cmd.Flags().GetString(FlagTokenURI)
			if err != nil {
				return err
			}
			tokenData, err := cmd.Flags().GetString(FlagTokenData)
			if err != nil {
				return err
			}
			msg := types.NewMsgEditNFT(
				args[1],
				args[0],
				tokenName,
				tokenURI,
				tokenData,
				clientCtx.GetFromAddress().String(),
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsEditNFT)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdTransferNft is the CLI command for sending a TransferNft transaction
func GetCmdTransferNft() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "transfer [from] [to] [denom-id] [token-id]",
		Long: "Transfer an NFT to a recipient.",
		Example: fmt.Sprintf(
			"$ %s tx nft transfer <from> <to> <denom-id> <token-id> "+
				"--from=<key-name> "+
				"--uri=<uri> "+
				"--chain-id=<chain-id> "+
				"--fees=<fee>",
			version.AppName,
		),
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			from := args[0]
			to := args[1]
			denomId := args[2]
			tokenId := args[3]
			msgSender := clientCtx.GetFromAddress().String()

			msg := types.NewMsgTransferNft(
				tokenId,
				denomId,
				from,
				to,
				msgSender,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsTransferNft)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdApproveNft  is the CLI command for grants permission to spender to transfer or send the given token
func GetCmdApproveNft() *cobra.Command {
	cmd := &cobra.Command{
		// no expire in ERC721
		Use:  "approve [approvedAddress][denom-id] [token-id] ",
		Long: "Adds the to address to the approved list.",
		Example: fmt.Sprintf(
			"$ %s tx nft approve <approvedAddress> <denom-id> <token-id> "+
				"--from=<key-name> "+
				"--chain-id=<chain-id> "+
				"--fees=<fee>",
			version.AppName,
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var sender = clientCtx.GetFromAddress().String()
			approvedAddress := args[0]
			denomId := args[1]
			tokenId := args[2]

			// nolint: govet
			if _, err := sdk.AccAddressFromBech32(approvedAddress); err != nil {
				return err
			}

			msg := types.NewMsgApproveNft(
				tokenId,
				denomId,
				sender,
				approvedAddress,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsApproveNft)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdApproveAllNFT is the CLI command to add a valid address to the users approved list
func GetCmdApproveAllNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "approveAll [operator] [approved]",
		Long: "Adds operatorToBeApproved address to the globally approved list of sender.",
		Example: fmt.Sprintf(
			"$ %s tx nft approveAll <operator> <true/false> "+
				"--from=<key-name> "+
				"--chain-id=<chain-id> "+
				"--fees=<fee>",
			version.AppName,
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sender := clientCtx.GetFromAddress().String()
			operator := args[0]
			approved, err := strconv.ParseBool(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgApproveAllNft(
				operator,
				sender,
				approved,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsApproveAllNft)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdRevokeNft is the CLI command for ownership transfer of the token to contract account
func GetCmdRevokeNft() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "revoke [addressToRevoke] [denom-id] [token-id]",
		Long: "Revokes a previously granted permission to transfer the given an NFT.",
		Example: fmt.Sprintf(
			"$ %s tx nft revoke <addressToRevoke> <denom-id> <token-id>"+
				"--from=<key-name> "+
				"--chain-id=<chain-id> "+
				"--fees=<fee>",
			version.AppName,
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			addressToRevoke := args[0]
			denomId := args[1]
			tokenId := args[2]
			sender := clientCtx.GetFromAddress().String()

			// nolint: govet
			if _, err := sdk.AccAddressFromBech32(addressToRevoke); err != nil {
				return err
			}

			msg := types.NewMsgRevokeNft(
				addressToRevoke,
				sender,
				denomId,
				tokenId,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsRevokeNft)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdBurnNFT is the CLI command for sending a BurnNFT transaction
func GetCmdBurnNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "burn [denom-id] [token-id]",
		Long: "Burn an NFT.",
		Example: fmt.Sprintf(
			"$ %s tx nft burn <denom-id> <token-id> "+
				"--from=<key-name> "+
				"--chain-id=<chain-id> "+
				"--fees=<fee>",
			version.AppName,
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgBurnNFT(
				clientCtx.GetFromAddress().String(),
				args[1],
				args[0],
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
