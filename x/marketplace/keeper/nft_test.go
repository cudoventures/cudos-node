package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	keepertest "github.com/CudoVentures/cudos-node/testutil/keeper"
	"github.com/CudoVentures/cudos-node/testutil/nullify"
	"github.com/CudoVentures/cudos-node/x/marketplace/keeper"
	"github.com/CudoVentures/cudos-node/x/marketplace/types"
)

func createNNft(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Nft {
	items := make([]types.Nft, n)
	for i := range items {
		items[i].Id = keeper.AppendNft(ctx, items[i])
	}
	return items
}

func TestNftGet(t *testing.T) {
	keeper, _, _, ctx := keepertest.MarketplaceKeeper(t)
	items := createNNft(keeper, ctx, 10)
	for _, item := range items {
		item := item
		got, found := keeper.GetNft(ctx, item.Id)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&got),
		)
	}
}

func TestNftRemove(t *testing.T) {
	keeper, _, _, ctx := keepertest.MarketplaceKeeper(t)
	items := createNNft(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveNft(ctx, item.Id)
		_, found := keeper.GetNft(ctx, item.Id)
		require.False(t, found)
	}
}

func TestNftGetAll(t *testing.T) {
	keeper, _, _, ctx := keepertest.MarketplaceKeeper(t)
	items := createNNft(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllNft(ctx)),
	)
}

func TestNftCount(t *testing.T) {
	keeper, _, _, ctx := keepertest.MarketplaceKeeper(t)
	items := createNNft(keeper, ctx, 10)
	count := uint64(len(items))
	require.Equal(t, count, keeper.GetNftCount(ctx))
}
