package keeper_test

import (
	"strconv"
	"testing"

	keepertest "github.com/CudoVentures/cudos-node/testutil/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/CudoVentures/cudos-node/testutil/nullify"
	"github.com/CudoVentures/cudos-node/x/token/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestQueryTokenByDenom(t *testing.T) {
	k, ctx := keepertest.TestTokenKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNToken(k, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryTokenByDenomRequest
		response *types.QueryTokenByDenomResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryTokenByDenomRequest{
				Denom: msgs[0].Denom,
			},
			response: &types.QueryTokenByDenomResponse{Token: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryTokenByDenomRequest{
				Denom: msgs[1].Denom,
			},
			response: &types.QueryTokenByDenomResponse{Token: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryTokenByDenomRequest{
				Denom: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := k.TokenByDenom(wctx, tc.request)
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

func TestQueryAllTokens(t *testing.T) {
	k, ctx := keepertest.TestTokenKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNToken(k, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllTokensRequest {
		return &types.QueryAllTokensRequest{
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
			resp, err := k.AllTokens(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Tokens), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Tokens),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := k.AllTokens(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Tokens), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Tokens),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := k.AllTokens(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Tokens),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := k.AllTokens(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
