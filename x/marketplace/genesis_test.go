package marketplace_test

import (
	"testing"

	keepertest "github.com/CudoVentures/cudos-node/testutil/keeper"
	"github.com/CudoVentures/cudos-node/testutil/nullify"
	"github.com/CudoVentures/cudos-node/x/marketplace"
	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	auctions := []types.Auction{
		&types.EnglishAuction{BaseAuction: &types.BaseAuction{Id: 0}},
		&types.DutchAuction{BaseAuction: &types.BaseAuction{Id: 1}},
	}
	auctionsAny, err := types.PackAuctions(auctions)
	require.NoError(t, err)

	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		CollectionList: []types.Collection{
			{
				Id: 0,
			},
			{
				Id: 1,
			},
		},
		CollectionCount: 2,
		NftList: []types.Nft{
			{
				Id: 0,
			},
			{
				Id: 1,
			},
		},
		NftCount:     2,
		AuctionList:  auctionsAny,
		AuctionCount: 2,
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, _, _, ctx := keepertest.MarketplaceKeeper(t)
	marketplace.InitGenesis(ctx, *k, genesisState)
	got := marketplace.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.CollectionList, got.CollectionList)
	require.Equal(t, genesisState.CollectionCount, got.CollectionCount)
	require.ElementsMatch(t, genesisState.NftList, got.NftList)
	require.Equal(t, genesisState.NftCount, got.NftCount)
	require.ElementsMatch(t, genesisState.AuctionList, got.AuctionList)
	require.Equal(t, genesisState.AuctionCount, got.AuctionCount)
	// this line is used by starport scaffolding # genesis/test/assert
}
