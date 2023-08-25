package cli_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/CudoVentures/cudos-node/simapp"
	"github.com/CudoVentures/cudos-node/testutil"
	"github.com/cosmos/cosmos-sdk/testutil/network"
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

func (s *TxTestSuite) TestIssueDenom() {
	network, err := testutil.RunNetwork(s.T(), s.config)
	require.NoError(s.T(), err)

	// clientCtx := network.Validators[0].ClientCtx
	// valAddr := network.Validators[0].Address.String()

	network.Cleanup()
}
