package cli

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/CudoVentures/cudos-node/x/nft/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/version"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                types.ModuleName,
		Short:              "Querying commands for the NFT module",
		DisableFlagParsing: true,
	}

	queryCmd.AddCommand(
		GetCmdQueryDenom(),
		GetCmdQueryDenomByName(),
		GetCmdQueryDenomBySymbol(),
		GetCmdQueryDenoms(),
		GetCmdQueryCollection(),
		GetCmdQueryCollectionsByDenomIds(),
		GetCmdQuerySupply(),
		GetCmdQueryOwner(),
		GetCmdQueryNFT(),
		GetCmdQueryApprovedNFT(),
		GetCmdQueryIsApprovedForAll(),
	)

	return queryCmd
}

// GetCmdQuerySupply queries the supply of a nft collection
func GetCmdQuerySupply() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "supply [denom-id]",
		Short:   "Query total supply of denom",
		Long:    "Query total supply of a collection or owner of NFTs.",
		Example: fmt.Sprintf("$ %s query nft supply <denom-id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var owner sdk.AccAddress
			ownerStr, err := cmd.Flags().GetString(FlagOwner)
			if err != nil {
				return err
			}

			if len(ownerStr) > 0 {
				owner, err = sdk.AccAddressFromBech32(ownerStr)
				if err != nil {
					return err
				}
			}

			if err := types.ValidateDenomID(args[0]); err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Supply(context.Background(), &types.QuerySupplyRequest{
				DenomId: args[0],
				Owner:   owner.String(),
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	cmd.Flags().AddFlagSet(FsQuerySupply)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryOwner queries all the NFTs owned by an account
// todo: change the name of this to something like QueryAllNFTsOfOwner..
func GetCmdQueryOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "owner [address]",
		Short:   "Query address NFT holdings",
		Long:    "Get the NFTs owned by an account address.",
		Example: fmt.Sprintf("$ %s query nft owner <address> --denom-id=<denom-id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if _, err := sdk.AccAddressFromBech32(args[0]); err != nil {
				return err
			}
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			denomID, err := cmd.Flags().GetString(FlagDenomID)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Owner(context.Background(), &types.QueryOwnerRequest{
				DenomId:    denomID,
				Owner:      args[0],
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	cmd.Flags().AddFlagSet(FsQueryOwner)
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "nfts")

	return cmd
}

// GetCmdQueryCollection queries all the NFTs from a collection
func GetCmdQueryCollection() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "collection [denom-id]",
		Short:   "Query NFTs in a collection",
		Long:    "Get all the NFTs from a given collection.",
		Example: fmt.Sprintf("$ %s query nft collection <denom-id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if err := types.ValidateDenomID(args[0]); err != nil {
				return err
			}
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Collection(
				context.Background(),
				&types.QueryCollectionRequest{
					DenomId:    args[0],
					Pagination: pageReq,
				},
			)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "nfts")

	return cmd
}

// GetCmdQueryCollectionsByDenomIds queries for all the collections matching a set of denom ids
func GetCmdQueryCollectionsByDenomIds() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "collection [denom-id1,denom-id2,denom-id3..]",
		Short:   "Query Collections by denom ids",
		Long:    "Get all the collections for given denom ids.",
		Example: fmt.Sprintf("$ %s query nft collectionByDenomIds <denom-id,denom-id2, denom-id3>", version.AppName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			denomIds := strings.Split(args[1], ",")

			if len(denomIds) == 0 {
				err := errors.New("denomIds array is empty")
				return err
			}

			for _, denomId := range denomIds {
				if err := types.ValidateDenomID(denomId); err != nil {
					return err
				}
			}
			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.CollectionsByDenomIds(context.Background(), &types.QueryCollectionsByIdsRequest{
				DenomIds: denomIds,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryDenoms queries all denoms
func GetCmdQueryDenoms() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "denoms",
		Short:   "Query all denoms",
		Long:    "Query all denominations of all collections of NFTs.",
		Example: fmt.Sprintf("$ %s query nft denoms", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Denoms(context.Background(), &types.QueryDenomsRequest{Pagination: pageReq})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "all denoms")
	return cmd
}

// GetCmdQueryDenom queries the specified denom
func GetCmdQueryDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "denom [denom-id]",
		Short:   "Query a denom by Id",
		Long:    "Query the denom by the specified denom id.",
		Example: fmt.Sprintf("$ %s query nft denom <denom-id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if err := types.ValidateDenomID(args[0]); err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Denom(
				context.Background(),
				&types.QueryDenomRequest{DenomId: args[0]},
			)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp.Denom)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryDenomByName queries the specified denom by name
func GetCmdQueryDenomByName() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "denom-by-name [denom-name]",
		Short:   "Query denom by name",
		Long:    "Query the denom by the specified denom name.",
		Example: fmt.Sprintf("$ %s query nft denom <denom-name>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if err := types.ValidateDenomName(args[0]); err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.DenomByName(
				context.Background(),
				&types.QueryDenomByNameRequest{DenomName: args[0]},
			)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp.Denom)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryDenomBySymbol queries the specified denom by symbol
func GetCmdQueryDenomBySymbol() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "denom-by-symbol [symbol]",
		Short:   "Query denom by symbol",
		Long:    "Query the denom by the specified symbol.",
		Example: fmt.Sprintf("$ %s query nft denom <symbol>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if err := types.ValidateDenomSymbol(args[0]); err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.DenomBySymbol(
				context.Background(),
				&types.QueryDenomBySymbolRequest{Symbol: args[0]},
			)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp.Denom)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryNFT queries a single NFTs from a collection
// todo: rename this to QueryNFT in the Use:
func GetCmdQueryNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "token [denom-id] [token-id]",
		Short:   "Query single NFT",
		Long:    "Query a single NFT from a collection by denom id and token id.",
		Example: fmt.Sprintf("$ %s query nft token <denom-id> <token-id>", version.AppName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if err := types.ValidateDenomID(args[0]); err != nil {
				return err
			}

			if err := types.ValidateTokenID(args[1]); err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.NFT(context.Background(), &types.QueryNFTRequest{
				DenomId: args[0],
				TokenId: args[1],
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp.NFT)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryApprovedNFT queries the NFT and returns its approved operators list
func GetCmdQueryApprovedNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "approvals [denomId] [tokenId]",
		Short:   "Query the approved addresses for a NFT.",
		Long:    "Get the approved addresses for the NFT.",
		Example: fmt.Sprintf("$ %s query nft approvals <denomId> <tokenId>", version.AppName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			denomId := args[0]
			if err := types.ValidateDenomID(denomId); err != nil {
				return err
			}

			tokenId := args[1]
			if err := types.ValidateTokenID(tokenId); err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.GetApprovalsNFT(context.Background(), &types.QueryApprovalsNFTRequest{
				DenomId: denomId,
				TokenId: tokenId,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	cmd.Flags().AddFlagSet(FsQueryOwner)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryIsApprovedForAll queries if the operator address is authorized for owner address
func GetCmdQueryIsApprovedForAll() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "is-approved-for-all [owner] [operator]",
		Short:   "Query if an address is approved operator",
		Long:    "Query if an address is an authorized operator for another address",
		Example: fmt.Sprintf("$ %s query nft isApprovedForAll <owner> <operator>", version.AppName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			owner := args[0]
			if _, err := sdk.AccAddressFromBech32(owner); err != nil {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
			}

			operator := args[1]
			if _, err := sdk.AccAddressFromBech32(operator); err != nil {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid operator address (%s)", err)
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.QueryApprovalsIsApprovedForAll(context.Background(), &types.QueryApprovalsIsApprovedForAllRequest{
				Owner:    owner,
				Operator: operator,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	cmd.Flags().AddFlagSet(FsQueryOwner)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
