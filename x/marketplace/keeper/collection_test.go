package keeper_test

import (
	"testing"

	keepertest "github.com/CudoVentures/cudos-node/testutil/keeper"
	"github.com/CudoVentures/cudos-node/testutil/nullify"
	"github.com/CudoVentures/cudos-node/x/marketplace/keeper"
	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func createNCollection(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Collection {
	items := make([]types.Collection, n)
	for i := range items {
		items[i].Id = keeper.AppendCollection(ctx, items[i])
	}
	return items
}

func TestCollectionGet(t *testing.T) {
	keeper, ctx := keepertest.MarketplaceKeeper(t)
	items := createNCollection(keeper, ctx, 10)
	for _, item := range items {
		got, found := keeper.GetCollection(ctx, item.Id)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&got),
		)
	}
}

func TestCollectionRemove(t *testing.T) {
	keeper, ctx := keepertest.MarketplaceKeeper(t)
	items := createNCollection(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveCollection(ctx, item.Id)
		_, found := keeper.GetCollection(ctx, item.Id)
		require.False(t, found)
	}
}

func TestCollectionGetAll(t *testing.T) {
	keeper, ctx := keepertest.MarketplaceKeeper(t)
	items := createNCollection(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllCollection(ctx)),
	)
}

func TestCollectionCount(t *testing.T) {
	keeper, ctx := keepertest.MarketplaceKeeper(t)
	items := createNCollection(keeper, ctx, 10)
	count := uint64(len(items))
	require.Equal(t, count, keeper.GetCollectionCount(ctx))
}
