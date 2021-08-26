package keeper

import (
	"context"

	"cudos.org/cudos-node/x/nft/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) NFTAll(c context.Context, req *types.QueryAllNFTRequest) (*types.QueryAllNFTResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var Nfts []*types.NFT
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	NftStore := prefix.NewStore(store, types.KeyPrefix(types.NFTKey))

	pageRes, err := query.Paginate(NftStore, req.Pagination, func(key []byte, value []byte) error {
		var Nft types.NFT
		if err := k.cdc.UnmarshalInterface(value, &Nft); err != nil {
			return err
		}

		Nfts = append(Nfts, &Nft)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllNFTResponse{NFT: Nfts, Pagination: pageRes}, nil
}

func (k Keeper) Nft(c context.Context, req *types.QueryGetNFTRequest) (*types.QueryNftResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetNFT(ctx, req.Index)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryNftResponse{NFT: &val}, nil
}
