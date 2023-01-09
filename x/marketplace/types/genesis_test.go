package types_test

import (
	"testing"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	auctions := []types.Auction{
		&types.EnglishAuction{BaseAuction: &types.BaseAuction{Id: 0}},
		&types.DutchAuction{BaseAuction: &types.BaseAuction{Id: 1}},
	}
	auctionsAny, err := types.PackAuctions(auctions)
	require.NoError(t, err)

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
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated collection",
			genState: &types.GenesisState{
				CollectionList: []types.Collection{
					{
						Id: 0,
					},
					{
						Id: 0,
					},
				},
			},
			valid: false,
		},
		{
			desc: "invalid collection count",
			genState: &types.GenesisState{
				CollectionList: []types.Collection{
					{
						Id: 1,
					},
				},
				CollectionCount: 0,
			},
			valid: false,
		},
		{
			desc: "duplicated nft",
			genState: &types.GenesisState{
				NftList: []types.Nft{
					{
						Id: 0,
					},
					{
						Id: 0,
					},
				},
			},
			valid: false,
		},
		{
			desc: "invalid nft count",
			genState: &types.GenesisState{
				NftList: []types.Nft{
					{
						Id: 1,
					},
				},
				NftCount: 0,
			},
			valid: false,
		},
		{
			desc: "duplicated auction",
			genState: &types.GenesisState{
				AuctionList: []*codectypes.Any{
					{
						Value: []byte("1"),
					},
					{
						Value: []byte("1"),
					},
				},
			},
			valid: false,
		},
		{
			desc: "invalid auction count",
			genState: &types.GenesisState{
				AuctionList: []*codectypes.Any{
					{
						Value: []byte("1"),
					},
				},
				AuctionCount: 0,
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
