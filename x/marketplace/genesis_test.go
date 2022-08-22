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
		NftCount: 2,
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.MarketplaceKeeper(t)
	marketplace.InitGenesis(ctx, *k, genesisState)
	got := marketplace.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.CollectionList, got.CollectionList)
	require.Equal(t, genesisState.CollectionCount, got.CollectionCount)
	require.ElementsMatch(t, genesisState.NftList, got.NftList)
	require.Equal(t, genesisState.NftCount, got.NftCount)
	// this line is used by starport scaffolding # genesis/test/assert
}
