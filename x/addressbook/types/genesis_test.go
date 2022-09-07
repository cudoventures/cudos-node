package types_test

import (
	"testing"

	"github.com/CudoVentures/cudos-node/testutil/sample"
	"github.com/CudoVentures/cudos-node/x/addressbook/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	creator := sample.AccAddress()
	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{

				AddressList: []types.Address{
					{
						Creator: creator,
						Network: "BTC",
						Label:   "100@testdenom",
					},
					{
						Creator: creator,
						Network: "BTC",
						Label:   "101@testdenom",
					},
				},
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated address",
			genState: &types.GenesisState{
				AddressList: []types.Address{
					{
						Creator: creator,
						Network: "BTC",
						Label:   "100@testdenom",
					},
					{
						Creator: creator,
						Network: "BTC",
						Label:   "100@testdenom",
					},
				},
			},
			valid: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
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
