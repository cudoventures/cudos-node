package keeper

import (
	"context"

	"github.com/CudoVentures/cudos-node/x/token/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) AllTokens(c context.Context, req *types.QueryAllTokensRequest) (*types.QueryAllTokensResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var tokens []types.Token
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	tokenStore := prefix.NewStore(store, types.KeyPrefix(types.TokenKeyPrefix))

	pageRes, err := query.Paginate(tokenStore, req.Pagination, func(key []byte, value []byte) error {
		var token types.Token
		if err := k.cdc.Unmarshal(value, &token); err != nil {
			return err
		}

		tokens = append(tokens, token)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllTokensResponse{Tokens: tokens, Pagination: pageRes}, nil
}

func (k Keeper) TokenByDenom(c context.Context, req *types.QueryTokenByDenomRequest) (*types.QueryTokenByDenomResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetTokenByDenom(
		ctx,
		req.Denom,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryTokenByDenomResponse{Token: val}, nil
}
