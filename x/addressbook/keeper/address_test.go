package keeper_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	keepertest "github.com/CudoVentures/cudos-node/testutil/keeper"
	"github.com/CudoVentures/cudos-node/testutil/nullify"
	"github.com/CudoVentures/cudos-node/testutil/sample"
	"github.com/CudoVentures/cudos-node/x/addressbook/keeper"
	"github.com/CudoVentures/cudos-node/x/addressbook/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNAddress(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Address {
	items := make([]types.Address, n)
	for i := range items {
		items[i].Creator = sample.AccAddress()
		items[i].Network = "BTC"
		items[i].Label = "1@testdenom"

		keeper.CreateNewAddress(ctx, items[i])
	}
	return items
}

func TestAddressGet(t *testing.T) {
	keeper, ctx := keepertest.AddressbookKeeper(t)
	items := createNAddress(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetAddress(ctx,
			item.Creator,
			item.Network,
			item.Label,
		)
		tc := tc
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}

func TestAddressRemove(t *testing.T) {
	keeper, ctx := keepertest.AddressbookKeeper(t)
	items := createNAddress(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveAddress(ctx,
			item.Creator,
			item.Network,
			item.Label,
		)
		_, found := keeper.GetAddress(ctx,
			item.Creator,
			item.Network,
			item.Label,
		)
		require.False(t, found)
	}
}

func TestAddressGetAll(t *testing.T) {
	keeper, ctx := keepertest.AddressbookKeeper(t)
	items := createNAddress(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllAddress(ctx)),
	)
}
