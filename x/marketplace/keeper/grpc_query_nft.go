package keeper

import (
	"context"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) NftAll(c context.Context, req *types.QueryAllNftRequest) (*types.QueryAllNftResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var nfts []types.Nft
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	nftStore := prefix.NewStore(store, types.KeyPrefix(types.NftKey))

	pageRes, err := query.Paginate(nftStore, req.Pagination, func(key []byte, value []byte) error {
		var nft types.Nft
		if err := k.cdc.Unmarshal(value, &nft); err != nil {
			return err
		}

		nfts = append(nfts, nft)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllNftResponse{Nft: nfts, Pagination: pageRes}, nil
}

func (k Keeper) Nft(c context.Context, req *types.QueryGetNftRequest) (*types.QueryGetNftResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	nft, found := k.GetNft(ctx, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryGetNftResponse{Nft: nft}, nil
}
