package keeper_test

import (
	"strconv"
	"testing"
	"time"

	keepertest "github.com/CudoVentures/cudos-node/testutil/keeper"
	"github.com/CudoVentures/cudos-node/testutil/nullify"
	"github.com/CudoVentures/cudos-node/x/marketplace/keeper"
	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func createNAuction(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Auction {
	auctions := make([]types.Auction, n)
	for i := 0; i < n; i++ {
		a := types.Auction{
			TokenId: strconv.FormatInt(int64(i), 10),
			DenomId: "test",
			EndTime: time.Date(2023, 1, 1, 15, 00, 00, 000, time.UTC),
			Creator: "test",
		}

		if i%2 == 0 {
			if err := a.SetAuctionType(
				&types.EnglishAuction{MinPrice: sdk.NewCoin("acudos", sdk.OneInt())},
			); err != nil {
				continue
			}
		} else {
			// todo dutch auction
			if err := a.SetAuctionType(&types.DutchAuction{}); err != nil {
				continue
			}
		}
		id, err := keeper.AppendAuction(ctx, a)
		if err != nil {
			continue
		}
		a.Id = id
		auctions[i] = a
	}
	return auctions
}

func TestAuctionGet(t *testing.T) {
	keeper, _, _, ctx := keepertest.MarketplaceKeeper(t)
	items := createNAuction(keeper, ctx, 10)
	for _, item := range items {
		got, found := keeper.GetAuction(ctx, item.Id)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&got),
		)
	}

	_, found := keeper.GetAuction(ctx, items[len(items)-1].Id+1)
	require.False(t, found)
}

func TestAuctionRemove(t *testing.T) {
	keeper, _, _, ctx := keepertest.MarketplaceKeeper(t)
	items := createNAuction(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveAuction(ctx, item.Id)
		_, found := keeper.GetAuction(ctx, item.Id)
		require.False(t, found)
	}
}

func TestAuctionGetAll(t *testing.T) {
	keeper, _, _, ctx := keepertest.MarketplaceKeeper(t)
	auctions := createNAuction(keeper, ctx, 10)
	require.Equal(t, auctions, keeper.GetAllAuction(ctx))
}

func TestAuctionCount(t *testing.T) {
	keeper, _, _, ctx := keepertest.MarketplaceKeeper(t)
	auctions := createNAuction(keeper, ctx, 10)
	count := uint64(len(auctions))
	require.Equal(t, count, keeper.GetAuctionCount(ctx))
}
