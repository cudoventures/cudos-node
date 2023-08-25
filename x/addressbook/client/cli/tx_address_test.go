package cli_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/CudoVentures/cudos-node/simapp"
	"github.com/CudoVentures/cudos-node/testutil"
	"github.com/CudoVentures/cudos-node/x/addressbook/client/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"
)

type TxAddressIntegrationTestSuite struct {
	suite.Suite
	config network.Config
}

func TestTxAddressIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(TxAddressIntegrationTestSuite))
}

func (s *TxAddressIntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up tx address integration test suite")

	s.config = simapp.NewConfig(s.T().TempDir())
	s.config.NumValidators = 1
}

func (s *TxAddressIntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down tx address integration test suite")
}

func (s *TxAddressIntegrationTestSuite) TestCreateAddress() {
	network, err := testutil.RunNetwork(s.T(), s.config)
	require.NoError(s.T(), err)

	ctx := network.Validators[0].ClientCtx
	valAddr := network.Validators[0].Address.String()

	fields := []string{"network", "label", "value"}
	for _, tc := range []struct {
		desc string

		args []string
		err  error
		code uint32
	}{
		{
			desc: "valid",
			args: []string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, valAddr),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(network.Config.BondDenom, sdk.NewInt(10))).String()),
			},
		},
	} {
		s.T().Run(tc.desc, func(t *testing.T) {
			args := []string{}
			args = append(args, fields...)
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdCreateAddress(), args)
			require.NoError(t, err)
			testutil.WaitForBlock()
			txResp, err := testutil.QueryJustBroadcastedTx(ctx, out)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.code, txResp.Code)
			}
		})
	}

	network.Cleanup()
}

func (s *TxAddressIntegrationTestSuite) TestUpdateAddress() {
	network, err := testutil.RunNetwork(s.T(), s.config)
	require.NoError(s.T(), err)

	ctx := network.Validators[0].ClientCtx
	valAddr := network.Validators[0].Address.String()

	existingKey := []string{"network", "label", "newvalue"}
	notFoundKey := []string{"network", "label1", "newvalue"}
	common := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, valAddr),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(network.Config.BondDenom, sdk.NewInt(10))).String()),
	}

	_, err = clitestutil.ExecTestCLICmd(ctx, cli.CmdCreateAddress(), append(existingKey, common...))
	require.NoError(s.T(), err)
	testutil.WaitForBlock()

	for _, tc := range []struct {
		desc string

		args []string
		code uint32
		err  error
	}{
		{
			desc: "valid",

			args: append(existingKey, common...),
		},
		{
			desc: "key not found",

			args: append(notFoundKey, common...),
			code: sdkerrors.ErrKeyNotFound.ABCICode(),
		},
	} {
		s.T().Run(tc.desc, func(t *testing.T) {
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdUpdateAddress(), tc.args)
			require.NoError(t, err)
			testutil.WaitForBlock()
			txResp, err := testutil.QueryJustBroadcastedTx(ctx, out)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.code, txResp.Code)
			}
		})
	}

	network.Cleanup()
}

func (s *TxAddressIntegrationTestSuite) TestDeleteAddress() {
	network, err := testutil.RunNetwork(s.T(), s.config)
	require.NoError(s.T(), err)

	ctx := network.Validators[0].ClientCtx
	valAddr := network.Validators[0].Address.String()

	existingKey := []string{"network", "label"}
	notFoundKey := []string{"network", "label1"}
	common := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, valAddr),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(network.Config.BondDenom, sdk.NewInt(10))).String()),
	}

	_, err = clitestutil.ExecTestCLICmd(ctx, cli.CmdCreateAddress(), append(append(existingKey, "value"), common...))
	require.NoError(s.T(), err)
	testutil.WaitForBlock()

	for _, tc := range []struct {
		desc string

		args []string
		code uint32
		err  error
	}{
		{
			desc: "valid",

			args: append(existingKey, common...),
		},
		{
			desc: "key not found",

			args: append(notFoundKey, common...),
			code: sdkerrors.ErrKeyNotFound.ABCICode(),
		},
	} {
		s.T().Run(tc.desc, func(t *testing.T) {
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdDeleteAddress(), tc.args)
			require.NoError(t, err)
			testutil.WaitForBlock()
			txResp, err := testutil.QueryJustBroadcastedTx(ctx, out)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.code, txResp.Code)
			}
		})
	}

	network.Cleanup()
}
