package keeper_test

import (
	"fmt"
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	keepertest "github.com/CudoVentures/cudos-node/testutil/keeper"
	"github.com/CudoVentures/cudos-node/x/addressbook/keeper"
	"github.com/CudoVentures/cudos-node/x/addressbook/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestAddressMsgServerCreate(t *testing.T) {
	k, ctx := keepertest.AddressbookKeeper(t)
	srv := keeper.NewMsgServerImpl(*k)
	wctx := sdk.WrapSDKContext(ctx)
	creator := "A"
	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateAddress{
			Creator: creator,
			Network: "ETH",
			Label:   fmt.Sprintf("%d@testdenom", i),
		}
		_, err := srv.CreateAddress(wctx, expected)
		require.NoError(t, err)
		rst, found := k.GetAddress(ctx, expected.Creator, expected.Network, expected.Label)
		require.True(t, found)
		require.Equal(t, expected.Creator, rst.Creator)
	}
}

func TestAddressMsgServerUpdate(t *testing.T) {
	creator := "A"

	for _, tc := range []struct {
		desc    string
		request *types.MsgUpdateAddress
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgUpdateAddress{
				Creator: creator,
				Network: "ETH",
				Label:   "0@testdenom",
			},
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgUpdateAddress{
				Creator: "B",
				Network: "UNK",
				Label:   "0@testdenom",
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.AddressbookKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)
			expected := &types.MsgCreateAddress{
				Creator: creator,
				Network: tc.request.Network,
				Label:   tc.request.Label,
			}
			_, err := srv.CreateAddress(wctx, expected)
			require.NoError(t, err)

			_, err = srv.UpdateAddress(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				rst, found := k.GetAddress(ctx, tc.request.Creator, tc.request.Network, tc.request.Label)
				require.True(t, found)
				require.Equal(t, expected.Creator, rst.Creator)
			}
		})
	}
}

func TestAddressMsgServerDelete(t *testing.T) {
	creator := "A"

	for _, tc := range []struct {
		desc    string
		request *types.MsgDeleteAddress
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgDeleteAddress{
				Creator: creator,
				Network: "ETH",
				Label:   "1@newdenom",
			},
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgDeleteAddress{
				Creator: creator,
				Network: "ETH",
				Label:   "3@newdenom",
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.AddressbookKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)

			_, err := srv.CreateAddress(wctx, &types.MsgCreateAddress{
				Creator: creator,
				Network: "ETH",
				Label:   "1@newdenom",
			})
			require.NoError(t, err)
			_, err = srv.DeleteAddress(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				_, found := k.GetAddress(ctx, tc.request.Creator, tc.request.Network, tc.request.Label)
				require.False(t, found)
			}
		})
	}
}
