package custom_bindings

import (
	"encoding/json"

	wasmKeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmd/x/wasm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"

	nftKeeper "github.com/CudoVentures/cudos-node/x/nft/keeper"
	nftTypes "github.com/CudoVentures/cudos-node/x/nft/types"
)

func PerformCustomNftQuery(keeper nftKeeper.Keeper) wasmKeeper.CustomQuerier {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		var custom nftCustomQuery
		err := json.Unmarshal(request, &custom)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
		}

		switch {
		case custom.QueryDenomById != nil:
			denom, err := keeper.GetDenom(ctx, custom.QueryDenomById.DenomID)
			if err != nil {
				return nil, err
			}
			return json.Marshal(nftTypes.QueryDenomResponse{Denom: &denom})
		case custom.QueryDenomByName != nil:
			denom, err := keeper.GetDenomByName(ctx, custom.QueryDenomByName.DenomName)
			if err != nil {
				return nil, err
			}
			return json.Marshal(nftTypes.QueryDenomResponse{Denom: &denom})
		case custom.QueryDenomBySymbol != nil:
			denom, err := keeper.GetDenomBySymbol(ctx, custom.QueryDenomBySymbol.Symbol)
			if err != nil {
				return nil, err
			}
			return json.Marshal(nftTypes.QueryDenomResponse{Denom: &denom})
		case custom.QueryDenoms != nil:
			denoms, err := keeper.Denoms(sdk.WrapSDKContext(ctx), &nftTypes.QueryDenomsRequest{Pagination: custom.QueryDenoms.Pagination})
			if err != nil {
				return nil, err
			}
			return json.Marshal(denoms)
		case custom.QueryCollection != nil:
			collection, err := keeper.Collection(sdk.WrapSDKContext(ctx), &nftTypes.QueryCollectionRequest{
				DenomID:    custom.QueryCollection.DenomID,
				Pagination: custom.QueryCollection.Pagination,
			})
			if err != nil {
				return nil, err
			}
			return json.Marshal(collection)
		case custom.QuerySupply != nil:
			denom, err := keeper.GetDenom(ctx, custom.QuerySupply.DenomID) // Otherwise queries for non-existing denom ID's will return 0, instead of erro.
			if err != nil {
				return nil, err
			}
			totalSupply := keeper.GetTotalSupply(ctx, denom.Id)
			return json.Marshal(nftTypes.QuerySupplyResponse{Amount: totalSupply})
		case custom.QueryOwner != nil:
			if len(custom.QueryOwner.Address) > 0 {
				owner, err := keeper.Owner(sdk.WrapSDKContext(ctx), &nftTypes.QueryOwnerRequest{
					DenomID:    custom.QueryOwner.DenomID,
					Owner:      custom.QueryOwner.Address,
					Pagination: custom.QueryOwner.Pagination,
				})
				if err != nil {
					return nil, err
				}
				return json.Marshal(owner)
			}
			return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Owner address is empty!")
		case custom.QueryToken != nil:
			nft, err := keeper.GetBaseNFT(ctx, custom.QueryToken.DenomID, custom.QueryToken.TokenId)
			if err != nil {
				return nil, err
			}
			return json.Marshal(nftTypes.QueryNFTResponse{NFT: &nft})
		case custom.QueryApprovals != nil:
			approvedAddressesForNft, err := keeper.GetNFTApprovedAddresses(ctx, custom.QueryApprovals.DenomID, custom.QueryApprovals.TokenId)
			if err != nil {
				return nil, err
			}
			return json.Marshal(nftTypes.QueryApprovalsNFTResponse{ApprovedAddresses: approvedAddressesForNft})
		case custom.QueryApprovedForAll != nil:
			if len(custom.QueryApprovedForAll.OwnerAddress) > 0 && len(custom.QueryApprovedForAll.OperatorAddress) > 0 {
				ownerAddress, err := sdk.AccAddressFromBech32(custom.QueryApprovedForAll.OwnerAddress)
				if err != nil {
					return nil, err
				}

				operatorAddress, err := sdk.AccAddressFromBech32(custom.QueryApprovedForAll.OperatorAddress)
				if err != nil {
					return nil, err
				}

				isApproved := keeper.IsApprovedOperator(ctx, ownerAddress, operatorAddress)
				return json.Marshal(nftTypes.QueryApprovalsIsApprovedForAllResponse{IsApproved: isApproved})
			}

		}
		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown Custom query variant")
	}
}

type nftCustomQuery struct {
	QueryDenomById      *QueryDenomById      `json:"query_denom_by_id,omitempty"`
	QueryDenomByName    *QueryDenomByName    `json:"query_denom_by_name,omitempty"`
	QueryDenomBySymbol  *QueryDenomBySymbol  `json:"query_denom_by_symbol,omitempty"`
	QueryDenoms         *QueryAllDenoms      `json:"query_denoms,omitempty"`
	QueryCollection     *QueryCollection     `json:"query_collection,omitempty"`
	QuerySupply         *QuerySupply         `json:"query_supply,omitempty"`
	QueryOwner          *QueryOwner          `json:"query_owner,omitempty"`
	QueryToken          *QueryToken          `json:"query_token,omitempty"`
	QueryApprovals      *QueryApprovals      `json:"query_approvals,omitempty"`
	QueryApprovedForAll *QueryApprovedForAll `json:"query_approved_for_all,omitempty"`
}

type QueryDenomById struct {
	DenomID string `json:"denom_id"`
}

type QueryDenomByName struct {
	DenomName string `json:"denom_name"`
}

type QueryDenomBySymbol struct {
	Symbol string `json:"denom_symbol"`
}

type QueryAllDenoms struct {
	Pagination *query.PageRequest `json:"pagination,omitempty"`
}

type QueryCollection struct {
	DenomID    string             `json:"denom_id"`
	Pagination *query.PageRequest `json:"pagination,omitempty"`
}

type QuerySupply struct {
	DenomID string `json:"denom_id"`
}

type QueryOwner struct {
	Address    string             `json:"address"`
	DenomID    string             `json:"denom_id,omitempty"`
	Pagination *query.PageRequest `json:"pagination,omitempty"`
}

type QueryToken struct {
	DenomID string `json:"denom_id"`
	TokenId string `json:"token_id"`
}

type QueryApprovals struct {
	DenomID string `json:"denom_id"`
	TokenId string `json:"token_id"`
}

type QueryApprovedForAll struct {
	OwnerAddress    string `json:"owner_address"`
	OperatorAddress string `json:"operator_address"`
}
