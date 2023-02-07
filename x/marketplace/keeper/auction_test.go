package keeper_test

import (
	"testing"

	keepertest "github.com/CudoVentures/cudos-node/testutil/keeper"
	"github.com/CudoVentures/cudos-node/x/marketplace/keeper"
	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func createNAuction(k *keeper.Keeper, ctx sdk.Context, n int) []types.Auction {
	auctions := make([]types.Auction, n)
	for i := 0; i < n; i++ {
		var a types.Auction
		if i%2 == 0 {
			a = &types.EnglishAuction{BaseAuction: &types.BaseAuction{}}
		} else {
			a = &types.DutchAuction{BaseAuction: &types.BaseAuction{}}
		}

		id, err := k.AppendAuction(ctx, a)
		if err != nil {
			continue
		}

		a.SetId(id)
		auctions[i] = a
	}
	return auctions
}

func TestAuctionGet(t *testing.T) {
	k, _, _, ctx := keepertest.MarketplaceKeeper(t)
	auctions := createNAuction(k, ctx, 2)

	for _, a := range auctions {
		haveAuction, err := k.GetAuction(ctx, a.GetId())
		require.NoError(t, err)
		require.Equal(t, a, haveAuction)
	}

	_, err := k.GetAuction(ctx, 9999)
	require.Error(t, err)
}

func TestAuctionRemove(t *testing.T) {
	k, _, _, ctx := keepertest.MarketplaceKeeper(t)
	auctions := createNAuction(k, ctx, 2)

	for _, a := range auctions {
		k.RemoveAuction(ctx, a.GetId())
		_, err := k.GetAuction(ctx, a.GetId())
		require.Error(t, err)
	}
}

func TestAuctionGetAll(t *testing.T) {
	k, _, _, ctx := keepertest.MarketplaceKeeper(t)
	auctions := createNAuction(k, ctx, 2)

	haveAuctions, err := k.GetAllAuction(ctx)
	require.NoError(t, err)
	require.Equal(t, auctions, haveAuctions)
}

func TestAuctionCount(t *testing.T) {
	k, _, _, ctx := keepertest.MarketplaceKeeper(t)
	auctions := createNAuction(k, ctx, 2)

	count := uint64(len(auctions))
	require.Equal(t, count, k.GetAuctionCount(ctx))
}
