package keeper_test

import (
	"strconv"
	"testing"

	keepertest "github.com/CudoVentures/cudos-node/testutil/keeper"
	"github.com/CudoVentures/cudos-node/testutil/nullify"
	"github.com/CudoVentures/cudos-node/testutil/sample"
	"github.com/CudoVentures/cudos-node/x/addressbook/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestAddressQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.AddressbookKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNAddress(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetAddressRequest
		response *types.QueryGetAddressResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetAddressRequest{
				Creator: msgs[0].Creator,
				Network: msgs[0].Network,
				Label:   msgs[0].Label,
			},
			response: &types.QueryGetAddressResponse{Address: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetAddressRequest{
				Creator: msgs[1].Creator,
				Network: msgs[1].Network,
				Label:   msgs[1].Label,
			},
			response: &types.QueryGetAddressResponse{Address: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetAddressRequest{
				Creator: sample.AccAddress(),
				Network: "BTC",
				Label:   "1@testdenom",
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Address(wctx, tc.request)
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

func TestAddressQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.AddressbookKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNAddress(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllAddressRequest {
		return &types.QueryAllAddressRequest{
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
			resp, err := keeper.AddressAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Address), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Address),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.AddressAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Address), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Address),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.AddressAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Address),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.AddressAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
