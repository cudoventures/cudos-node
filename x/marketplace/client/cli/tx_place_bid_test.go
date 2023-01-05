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
)

// todo check refactor
func TestPlaceBid(t *testing.T) {
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
		desc   string
		args   []string
		errMsg string
	}{
		{
			desc: "valid",
			args: []string{"1", "1acudos"},
		},
		{
			desc:   "invalid msg.ValidateBasic",
			args:   []string{"1", "0acudos"},
			errMsg: "amount must be positive: invalid price",
		},
		{
			desc:   "invalid coin arg",
			args:   []string{"1", "invalid"},
			errMsg: "invalid decimal coin expression: invalid",
		},
		{
			desc:   "invalid auctionID arg",
			args:   []string{"invalid", "1acudos"},
			errMsg: "strconv.ParseUint: parsing \"invalid\": invalid syntax",
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			args := append(tc.args, flags...)
			_, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.CmdPlaceBid(), args)

			if tc.errMsg != "" {
				require.EqualError(t, err, tc.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
