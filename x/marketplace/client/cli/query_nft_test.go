package cli_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/CudoVentures/cudos-node/simapp"
	"github.com/CudoVentures/cudos-node/testutil/nullify"
	"github.com/CudoVentures/cudos-node/x/marketplace/client/cli"
	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
)

func networkWithNftObjects(t *testing.T, n int) (*network.Network, []types.Nft) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	for i := 0; i < n; i++ {
		nft := types.Nft{
			Id: uint64(i),
		}
		nullify.Fill(&nft)
		state.NftList = append(state.NftList, nft)
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), state.NftList
}

type QueryNftIntegrationTestSuite struct {
	suite.Suite
	network *network.Network
	nftList []types.Nft
}

func TestQueryNftIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(QueryNftIntegrationTestSuite))
}

func (s *QueryNftIntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up query nft integration test suite")

	cfg := simapp.NewConfig()
	cfg.NumValidators = 1

	state := types.GenesisState{}
	require.NoError(s.T(), cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	for i := 0; i < 5; i++ {
		nft := types.Nft{
			Id: uint64(i),
		}
		nullify.Fill(&nft)
		state.NftList = append(state.NftList, nft)
	}
	s.nftList = state.NftList

	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(s.T(), err)
	cfg.GenesisState[types.ModuleName] = buf

	s.network = network.New(s.T(), cfg)

	_, err = s.network.WaitForHeight(3) // The network is fully initialized after 3 blocks
	s.Require().NoError(err)
}

func (s *QueryNftIntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down query nft integration test suite")
	s.network.Cleanup()
}

func (s *QueryNftIntegrationTestSuite) TestShowNft() {
	ctx := s.network.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	for _, tc := range []struct {
		desc string
		id   string
		args []string
		err  error
		obj  types.Nft
	}{
		{
			desc: "found",
			id:   fmt.Sprintf("%d", s.nftList[0].Id),
			args: common,
			obj:  s.nftList[0],
		},
		{
			desc: "not found",
			id:   "not_found",
			args: common,
			err:  status.Error(codes.NotFound, "not found"),
		},
	} {
		tc := tc
		s.T().Run(tc.desc, func(t *testing.T) {
			args := []string{tc.id}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowNft(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryGetNftResponse
				require.NoError(t, s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.Nft)
				require.Equal(t,
					nullify.Fill(&tc.obj),
					nullify.Fill(&resp.Nft),
				)
			}
		})
	}
}

func (s *QueryNftIntegrationTestSuite) TestListNft() {
	ctx := s.network.Validators[0].ClientCtx
	request := func(next []byte, offset, limit uint64, total bool) []string {
		args := []string{
			fmt.Sprintf("--%s=json", tmcli.OutputFlag),
		}
		if next == nil {
			args = append(args, fmt.Sprintf("--%s=%d", flags.FlagOffset, offset))
		} else {
			args = append(args, fmt.Sprintf("--%s=%s", flags.FlagPageKey, next))
		}
		args = append(args, fmt.Sprintf("--%s=%d", flags.FlagLimit, limit))
		if total {
			args = append(args, fmt.Sprintf("--%s", flags.FlagCountTotal))
		}
		return args
	}
	s.T().Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(s.nftList); i += step {
			args := request(nil, uint64(i), uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListNft(), args)
			require.NoError(t, err)
			var resp types.QueryAllNftResponse
			require.NoError(t, s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.Nft), step)
			require.Subset(t,
				nullify.Fill(s.nftList),
				nullify.Fill(resp.Nft),
			)
		}
	})
	s.T().Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(s.nftList); i += step {
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListNft(), args)
			require.NoError(t, err)
			var resp types.QueryAllNftResponse
			require.NoError(t, s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.Nft), step)
			require.Subset(t,
				nullify.Fill(s.nftList),
				nullify.Fill(resp.Nft),
			)
			next = resp.Pagination.NextKey
		}
	})
	s.T().Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(s.nftList)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListNft(), args)
		require.NoError(t, err)
		var resp types.QueryAllNftResponse
		require.NoError(t, s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NoError(t, err)
		require.Equal(t, len(s.nftList), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(s.nftList),
			nullify.Fill(resp.Nft),
		)
	})
}
