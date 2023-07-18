package cli_test

import (
	"fmt"
	"testing"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/CudoVentures/cudos-node/simapp"
	"github.com/CudoVentures/cudos-node/testutil/nullify"
	"github.com/CudoVentures/cudos-node/x/marketplace/client/cli"
	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
)

type QueryCollectionIntegrationTestSuite struct {
	suite.Suite
	network        *network.Network
	collectionList []types.Collection
}

func TestQueryCollectionIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(QueryCollectionIntegrationTestSuite))
}

func (s *QueryCollectionIntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up query collection integration test suite")

	cfg := simapp.NewConfig(s.T().TempDir())
	cfg.NumValidators = 1

	state := types.GenesisState{}
	require.NoError(s.T(), cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	for i := 0; i < 5; i++ {
		collection := types.Collection{
			Id: uint64(i),
		}
		nullify.Fill(&collection)
		state.CollectionList = append(state.CollectionList, collection)
	}
	s.collectionList = state.CollectionList

	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(s.T(), err)
	cfg.GenesisState[types.ModuleName] = buf

	s.network, err = network.New(s.T(), s.T().TempDir(), cfg)
	s.Require().NoError(err)

	_, err = s.network.WaitForHeight(3) // The network is fully initialized after 3 blocks
	s.Require().NoError(err)
}

func (s *QueryCollectionIntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down query collection integration test suite")
	s.network.Cleanup()
}

func (s *QueryCollectionIntegrationTestSuite) TestShowCollection() {
	ctx := s.network.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	for _, tc := range []struct {
		desc string
		id   string
		args []string
		err  error
		obj  types.Collection
	}{
		{
			desc: "found",
			id:   fmt.Sprintf("%d", s.collectionList[0].Id),
			args: common,
			obj:  s.collectionList[0],
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
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowCollection(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryGetCollectionResponse
				require.NoError(t, s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.Collection)
				require.Equal(t,
					nullify.Fill(&tc.obj),
					nullify.Fill(&resp.Collection),
				)
			}
		})
	}
}

func (s *QueryCollectionIntegrationTestSuite) TestListCollection() {
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
		for i := 0; i < len(s.collectionList); i += step {
			args := request(nil, uint64(i), uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListCollection(), args)
			require.NoError(t, err)
			var resp types.QueryAllCollectionResponse
			require.NoError(t, s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.Collection), step)
			require.Subset(t,
				nullify.Fill(s.collectionList),
				nullify.Fill(resp.Collection),
			)
		}
	})
	s.T().Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(s.collectionList); i += step {
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListCollection(), args)
			require.NoError(t, err)
			var resp types.QueryAllCollectionResponse
			require.NoError(t, s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.Collection), step)
			require.Subset(t,
				nullify.Fill(s.collectionList),
				nullify.Fill(resp.Collection),
			)
			next = resp.Pagination.NextKey
		}
	})
	s.T().Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(s.collectionList)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListCollection(), args)
		require.NoError(t, err)
		var resp types.QueryAllCollectionResponse
		require.NoError(t, s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NoError(t, err)
		require.Equal(t, len(s.collectionList), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(s.collectionList),
			nullify.Fill(resp.Collection),
		)
	})
}
