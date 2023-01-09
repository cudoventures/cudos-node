package keeper

import (
	"context"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AuctionAll(
	c context.Context, req *types.QueryAllAuctionRequest,
) (*types.QueryAllAuctionResponse, error) {
	if req == nil {
		return nil, sdkerrors.ErrInvalidRequest
	}

	var auctions []*codectypes.Any
	ctx := sdk.UnwrapSDKContext(c)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AuctionKey))

	pageRes, err := query.Paginate(store, req.Pagination, func(key, value []byte) error {
		var a types.Auction
		if err := k.cdc.UnmarshalInterface(value, &a); err != nil {
			return err
		}

		auctionAny, err := types.PackAuction(a)
		if err != nil {
			return err
		}

		auctions = append(auctions, auctionAny)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &types.QueryAllAuctionResponse{Auctions: auctions, Pagination: pageRes}, nil
}

func (k Keeper) Auction(
	ctx context.Context, req *types.QueryGetAuctionRequest,
) (*types.QueryGetAuctionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	a, err := k.GetAuction(sdk.UnwrapSDKContext(ctx), req.Id)
	if err != nil {
		return nil, types.ErrAuctionNotFound
	}

	auctionAny, err := types.PackAuction(a)
	if err != nil {
		return nil, err
	}

	return &types.QueryGetAuctionResponse{Auction: auctionAny}, nil
}
