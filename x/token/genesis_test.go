package token

import (
	"testing"

	keepertest "github.com/CudoVentures/cudos-node/testutil/keeper"
	"github.com/CudoVentures/cudos-node/testutil/nullify"
	"github.com/CudoVentures/cudos-node/x/token/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState *GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &GenesisState{

				Tokens: []types.Token{
					{
						Denom: "0",
					},
					{
						Denom: "1",
					},
				},
			},
			valid: true,
		},
		{
			desc: "duplicated token",
			genState: &GenesisState{
				Tokens: []types.Token{
					{
						Denom: "0",
					},
					{
						Denom: "0",
					},
				},
			},
			valid: false,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestGenesis(t *testing.T) {
	genesisState := GenesisState{
		Tokens: []types.Token{
			{
				Denom: "0",
			},
			{
				Denom: "1",
			},
		},
	}

	k, ctx := keepertest.TestTokenKeeper(t)
	InitGenesis(ctx, *k, genesisState)
	got := ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.Tokens, got.Tokens)
}
