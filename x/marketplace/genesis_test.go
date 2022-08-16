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

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.MarketplaceKeeper(t)
	marketplace.InitGenesis(ctx, *k, genesisState)
	got := marketplace.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
