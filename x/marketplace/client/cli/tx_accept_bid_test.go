package cli_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"

	"github.com/CudoVentures/cudos-node/simapp"
	"github.com/CudoVentures/cudos-node/testutil/network"
	"github.com/CudoVentures/cudos-node/x/marketplace/client/cli"
	"github.com/CudoVentures/cudos-node/x/marketplace/types"
)

func TestAcceptBid(t *testing.T) {
	cfg := simapp.NewConfig()
	cfg.NumValidators = 1
	net := network.New(t, cfg)
	_, err := net.WaitForHeight(3)
	require.NoError(t, err)
	val := net.Validators[0]

	flags := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	}

	for _, tc := range []struct {
		desc    string
		args    []string
		wantErr error
	}{
		{
			desc: "valid",
			args: []string{"1"},
		},
		{
			desc:    "invalid auction id",
			args:    []string{"invalid"},
			wantErr: types.ErrInvalidAuctionId,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			args := append(tc.args, flags...)
			_, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.CmdAcceptBid(), args)
			require.ErrorIs(t, err, tc.wantErr)
		})
	}
}
