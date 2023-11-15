package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
)

func (k Keeper) CollectionByDenomId(goCtx context.Context, req *types.QueryCollectionByDenomIdRequest) (*types.QueryCollectionByDenomIdResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	collection, found := k.GetCollectionByDenomID(ctx, req.DenomId)
	if !found {
		return &types.QueryCollectionByDenomIdResponse{}, nil
	}

	return &types.QueryCollectionByDenomIdResponse{Collection: collection}, nil
}
