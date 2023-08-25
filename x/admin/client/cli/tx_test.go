package cli_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/CudoVentures/cudos-node/simapp"
	"github.com/CudoVentures/cudos-node/testutil"
	"github.com/CudoVentures/cudos-node/testutil/sample"
	"github.com/CudoVentures/cudos-node/x/admin/client/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"

	distrcli "github.com/cosmos/cosmos-sdk/x/distribution/client/cli"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	bankcli "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type TxTestSuite struct {
	suite.Suite
	config network.Config
}

func TestTxTestSuite(t *testing.T) {
	suite.Run(t, new(TxTestSuite))
}

func (s *TxTestSuite) SetupSuite() {
	s.T().Log("setting up tx test suite")

	s.config = simapp.NewConfig(s.T().TempDir())
	s.config.NumValidators = 1
}

func (s *TxTestSuite) TearDownSuite() {
	s.T().Log("tearing down tx test suite")
}

func (s *TxTestSuite) TestAdminSpendCommunityPool() {
	network, err := testutil.RunNetwork(s.T(), s.config)
	require.NoError(s.T(), err)

	communityPoolReceiver := sample.AccAddress()

	clientCtx := network.Validators[0].ClientCtx
	valAddr := network.Validators[0].Address.String()

	s.T().Run("community sped", func(t *testing.T) {
		txRes, txErr := clitestutil.ExecTestCLICmd(clientCtx, distrcli.GetCmdQueryCommunityPool(), make([]string, 0))
		require.NoError(s.T(), txErr)

		var queryCommunityPoolResponse distrtypes.QueryCommunityPoolResponse
		require.NoError(t, network.Config.Codec.UnmarshalJSON(txRes.Bytes(), &queryCommunityPoolResponse))

		args := []string{
			communityPoolReceiver,
			fmt.Sprintf("%s%s", queryCommunityPoolResponse.Pool[0].Amount.RoundInt().String(), queryCommunityPoolResponse.Pool[0].Denom),
			fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
			fmt.Sprintf("--%s=%s", flags.FlagFrom, valAddr),
			fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(network.Config.BondDenom, sdk.NewInt(10))).String()),
		}
		txRes, txErr = clitestutil.ExecTestCLICmd(clientCtx, cli.CmdAdminSpendCommunityPool(), args)
		require.NoError(s.T(), txErr)
		testutil.WaitForBlock()

		txResp, txErr := testutil.QueryJustBroadcastedTx(clientCtx, txRes)
		require.NoError(t, txErr)
		require.Equal(t, uint32(0), txResp.Code)

		txRes, txErr = clitestutil.ExecTestCLICmd(clientCtx, bankcli.GetBalancesCmd(), []string{
			communityPoolReceiver,
		})
		require.NoError(s.T(), txErr)

		var queryAllBalancesResponse banktypes.QueryAllBalancesResponse
		require.NoError(t, network.Config.Codec.UnmarshalJSON(txRes.Bytes(), &queryAllBalancesResponse))
		require.Equal(t, queryCommunityPoolResponse.Pool[0].Amount.RoundInt(), queryAllBalancesResponse.Balances[0].Amount)
		require.Equal(t, queryCommunityPoolResponse.Pool[0].Denom, queryAllBalancesResponse.Balances[0].Denom)
	})

	network.Cleanup()
}
