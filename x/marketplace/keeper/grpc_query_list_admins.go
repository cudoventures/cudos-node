package keeper

import (
	"context"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ListAdmins(goCtx context.Context, req *types.QueryListAdminsRequest) (*types.QueryListAdminsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	admins, err := k.GetAdmins(sdk.UnwrapSDKContext(goCtx))
	if err != nil {
		return &types.QueryListAdminsResponse{}, err
	}

	return &types.QueryListAdminsResponse{Admins: admins}, nil
}
