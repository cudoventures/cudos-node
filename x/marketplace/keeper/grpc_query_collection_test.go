package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/CudoVentures/cudos-node/testutil/keeper"
	"github.com/CudoVentures/cudos-node/testutil/nullify"
	"github.com/CudoVentures/cudos-node/x/marketplace/types"
)

func TestCollectionQuerySingle(t *testing.T) {
	keeper, _, ctx := keepertest.MarketplaceKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNCollection(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetCollectionRequest
		response *types.QueryGetCollectionResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetCollectionRequest{Id: msgs[0].Id},
			response: &types.QueryGetCollectionResponse{Collection: msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetCollectionRequest{Id: msgs[1].Id},
			response: &types.QueryGetCollectionResponse{Collection: msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetCollectionRequest{Id: uint64(len(msgs))},
			err:     sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Collection(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}

func TestCollectionQueryPaginated(t *testing.T) {
	keeper, _, ctx := keepertest.MarketplaceKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNCollection(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllCollectionRequest {
		return &types.QueryAllCollectionRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.CollectionAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Collection), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Collection),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.CollectionAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Collection), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Collection),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.CollectionAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Collection),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.CollectionAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
