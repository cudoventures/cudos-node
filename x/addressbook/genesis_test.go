package addressbook_test

import (
	"testing"

	keepertest "github.com/CudoVentures/cudos-node/testutil/keeper"
	"github.com/CudoVentures/cudos-node/testutil/nullify"
	"github.com/CudoVentures/cudos-node/testutil/sample"
	"github.com/CudoVentures/cudos-node/x/addressbook"
	"github.com/CudoVentures/cudos-node/x/addressbook/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		AddressList: []types.Address{
			{
				Creator: sample.AccAddress(),
				Network: "BTC",
				Label:   "1@testdenom",
			},
			{
				Creator: sample.AccAddress(),
				Network: "ETH",
				Label:   "2@newdenom",
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.AddressbookKeeper(t)
	addressbook.InitGenesis(ctx, *k, genesisState)
	got := addressbook.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.AddressList, got.AddressList)
	// this line is used by starport scaffolding # genesis/test/assert
}
