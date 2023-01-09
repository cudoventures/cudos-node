package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/CudoVentures/cudos-node/testutil/keeper"
	"github.com/CudoVentures/cudos-node/x/marketplace/types"
)

func TestAuctionQuerySingle(t *testing.T) {
	k, _, _, ctx := keepertest.MarketplaceKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	auctions := createNAuction(k, ctx, 2)
	auctionsAny, err := types.PackAuctions(auctions)
	require.NoError(t, err)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetAuctionRequest
		response *types.QueryGetAuctionResponse
		wantErr  error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetAuctionRequest{Id: auctions[0].GetId()},
			response: &types.QueryGetAuctionResponse{Auction: auctionsAny[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetAuctionRequest{Id: auctions[1].GetId()},
			response: &types.QueryGetAuctionResponse{Auction: auctionsAny[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetAuctionRequest{Id: uint64(len(auctions))},
			wantErr: types.ErrAuctionNotFound,
		},
		{
			desc:    "InvalidRequest",
			wantErr: status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			haveResp, err := k.Auction(wctx, tc.request)
			require.ErrorIs(t, err, tc.wantErr)
			if tc.wantErr == nil {
				require.Equal(t, tc.response, haveResp)
			}
		})
	}
}

func TestAuctionQueryPaginated(t *testing.T) {
	keeper, _, _, ctx := keepertest.MarketplaceKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	auctions := createNAuction(keeper, ctx, 5)

	q := func(next []byte, offset, limit uint64, total bool) *types.QueryAllAuctionRequest {
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
		step := uint64(2)
		for i := 0; i < len(auctions); i += int(step) {
			resp, err := keeper.AuctionAll(wctx, q(nil, uint64(i), step, false))
			require.NoError(t, err)
			haveAuctions, err := types.UnpackAuctions(resp.Auctions)
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Auctions), int(step))
			require.Subset(t, auctions, haveAuctions)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(auctions); i += step {
			resp, err := keeper.AuctionAll(wctx, q(next, 0, uint64(step), false))
			require.NoError(t, err)
			haveAuctions, err := types.UnpackAuctions(resp.Auctions)
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Auctions), step)
			require.Subset(t, auctions, haveAuctions)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.AuctionAll(wctx, q(nil, 0, 0, true))
		require.NoError(t, err)
		haveAuctions, err := types.UnpackAuctions(resp.Auctions)
		require.NoError(t, err)
		require.Equal(t, len(auctions), int(resp.Pagination.Total))
		require.Equal(t, auctions, haveAuctions)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.AuctionAll(wctx, nil)
		require.Error(t, err)
	})
	t.Run("InvalidPageRequest", func(t *testing.T) {
		_, err := keeper.AuctionAll(wctx, q([]byte("invalid"), 1, 0, true))
		require.Error(t, err)
	})
}
