package custom_bindings

import (
	"encoding/json"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/CudoVentures/cudos-node/x/marketplace/keeper"
	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
)

func PerformCustomMarketplaceQuery(keeper keeper.Keeper) wasmkeeper.CustomQuerier {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		var custom marketplaceCustomQuery
		err := json.Unmarshal(request, &custom)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
		}

		switch {
		case custom.QueryCollection != nil:
			collection, err := keeper.Collection(ctx, &types.QueryGetCollectionRequest{Id: custom.QueryCollection.Id})
			if err != nil {
				return nil, err
			}
			return json.Marshal(collection)
		case custom.QueryAllCollections != nil:
			collection, err := keeper.CollectionAll(ctx, &types.QueryAllCollectionRequest{Pagination: custom.QueryAllCollections.Pagination})
			if err != nil {
				return nil, err
			}
			return json.Marshal(collection)
		case custom.QueryCollectionByDenomId != nil:
			collection, err := keeper.CollectionByDenomId(ctx, &types.QueryCollectionByDenomIdRequest{DenomId: custom.QueryCollectionByDenomId.DenomId})
			if err != nil {
				return nil, err
			}
			return json.Marshal(collection)
		case custom.QueryNft != nil:
			nft, err := keeper.Nft(ctx, &types.QueryGetNftRequest{Id: custom.QueryNft.Id})
			if err != nil {
				return nil, err
			}
			return json.Marshal(nft)
		case custom.QueryNftAll != nil:
			nft, err := keeper.NftAll(ctx, &types.QueryAllNftRequest{Pagination: custom.QueryNftAll.Pagination})
			if err != nil {
				return nil, err
			}
			return json.Marshal(nft)
		case custom.QueryListAdmins != nil:
			admins, err := keeper.ListAdmins(ctx, &types.QueryListAdminsRequest{})
			if err != nil {
				return nil, err
			}
			return json.Marshal(admins)

		}
		return nil, sdkerrors.Wrap(wasmtypes.ErrInvalidMsg, "Unknown Custom query variant")
	}
}

type marketplaceCustomQuery struct {
	QueryCollection          *QueryCollection          `json:"query_collection_marketplace,omitempty"`
	QueryCollectionByDenomId *QueryCollectionByDenomId `json:"query_collection_by_denom_id,omitempty"`
	QueryAllCollections      *QueryAllCollections      `json:"query_all_collections,omitempty"`
	QueryNft                 *QueryNft                 `json:"query_nft,omitempty"`
	QueryNftAll              *QueryNftAll              `json:"query_all_nfts,omitempty"`
	QueryListAdmins          *QueryListAdmins          `json:"query_list_admins,omitempty"`
}

type QueryCollection struct {
	Id uint64 `json:"id"`
}

type QueryAllCollections struct {
	Pagination *query.PageRequest `json:"pagination,omitempty"`
}

type QueryCollectionByDenomId struct {
	DenomId string `json:"denom_id"`
}

type QueryNft struct {
	Id uint64 `json:"id"`
}

type QueryNftAll struct {
	Pagination *query.PageRequest `json:"pagination,omitempty"`
}

type QueryListAdmins struct {
}
