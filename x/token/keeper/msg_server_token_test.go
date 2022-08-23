package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	keepertest "github.com/CudoVentures/cudos-node/testutil/keeper"
	"github.com/CudoVentures/cudos-node/x/token/keeper"
	"github.com/CudoVentures/cudos-node/x/token/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestTokenMsgServerCreate(t *testing.T) {
	k, ctx := keepertest.TokenKeeper(t)
	srv := keeper.NewMsgServerImpl(*k)
	wctx := sdk.WrapSDKContext(ctx)
	owner := "A"
	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateToken{Owner: owner,
			Denom: strconv.Itoa(i),
		}
		_, err := srv.CreateToken(wctx, expected)
		require.NoError(t, err)
		rst, found := k.GetToken(ctx,
			expected.Denom,
		)
		require.True(t, found)
		require.Equal(t, expected.Owner, rst.Owner)
	}
}

func TestTokenMsgServerUpdate(t *testing.T) {
	owner := "A"

	for _, tc := range []struct {
		desc    string
		request *types.MsgUpdateToken
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgUpdateToken{Owner: owner,
				Denom: strconv.Itoa(0),
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgUpdateToken{Owner: "B",
				Denom: strconv.Itoa(0),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgUpdateToken{Owner: owner,
				Denom: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.TokenKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)
			expected := &types.MsgCreateToken{Owner: owner,
				Denom: strconv.Itoa(0),
			}
			_, err := srv.CreateToken(wctx, expected)
			require.NoError(t, err)

			_, err = srv.UpdateToken(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				rst, found := k.GetToken(ctx,
					expected.Denom,
				)
				require.True(t, found)
				require.Equal(t, expected.Owner, rst.Owner)
			}
		})
	}
}

func TestTokenMsgServerDelete(t *testing.T) {
	owner := "A"

	for _, tc := range []struct {
		desc    string
		request *types.MsgDeleteToken
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgDeleteToken{Owner: owner,
				Denom: strconv.Itoa(0),
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgDeleteToken{Owner: "B",
				Denom: strconv.Itoa(0),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgDeleteToken{Owner: owner,
				Denom: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.TokenKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)

			_, err := srv.CreateToken(wctx, &types.MsgCreateToken{Owner: owner,
				Denom: strconv.Itoa(0),
			})
			require.NoError(t, err)
			_, err = srv.DeleteToken(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				_, found := k.GetToken(ctx,
					tc.request.Denom,
				)
				require.False(t, found)
			}
		})
	}
}
