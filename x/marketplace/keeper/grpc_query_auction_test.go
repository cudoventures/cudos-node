package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/CudoVentures/cudos-node/testutil/keeper"
	"github.com/CudoVentures/cudos-node/testutil/nullify"
	"github.com/CudoVentures/cudos-node/x/marketplace/types"
)

func TestAuctionQuerySingle(t *testing.T) {
	keeper, _, _, ctx := keepertest.MarketplaceKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	// todo these auctions are completely empty and nothing is tested add auctionType and etc
	msgs := createNAuction(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetAuctionRequest
		response *types.QueryGetAuctionResponse
		errMsg   string
	}{
		{
			desc:     "First",
			request:  &types.QueryGetAuctionRequest{Id: msgs[0].Id},
			response: &types.QueryGetAuctionResponse{Auction: msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetAuctionRequest{Id: msgs[1].Id},
			response: &types.QueryGetAuctionResponse{Auction: msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetAuctionRequest{Id: uint64(len(msgs))},
			errMsg:  "key not found",
		},
		{
			desc:   "InvalidRequest",
			errMsg: "rpc error: code = InvalidArgument desc = invalid request",
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Auction(wctx, tc.request)
			if tc.errMsg != "" {
				require.EqualError(t, err, tc.errMsg)
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

func TestAuctionQueryPaginated(t *testing.T) {
	keeper, _, _, ctx := keepertest.MarketplaceKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNAuction(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllAuctionRequest {
		return &types.QueryAllAuctionRequest{
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
			resp, err := keeper.AuctionAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Auction), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Auction),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.AuctionAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Auction), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Auction),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.AuctionAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Auction),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.AuctionAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
	t.Run("InvalidPageRequest", func(t *testing.T) {
		_, err := keeper.AuctionAll(wctx, request([]byte("invalid"), 1, 0, true))
		require.ErrorIs(t, err, status.Error(codes.Internal, "invalid request, either offset or key is expected, got both"))
	})
}
