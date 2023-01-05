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

func TestPublishAuction(t *testing.T) {
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
			// todo make a valid test case with DutchAuction
			desc: "valid",
			args: []string{
				"1",
				"xyz",
				"25h",
				`{"@type":"/cudoventures.cudosnode.marketplace.EnglishAuction","minPrice":{"denom":"stake","amount":"1"}}`,
			},
		},
		{
			desc: "invalid tokenID arg",
			args: []string{
				"invalid",
				"xyz",
				"25h",
				`{"@type":"/cudoventures.cudosnode.marketplace.EnglishAuction","minPrice":{"denom":"stake","amount":"1"}}`,
			},
			errMsg: "invalid nft: invalid address",
		},
		{
			desc: "invalid auctionType arg",
			args: []string{
				"1",
				"xyz",
				"25h",
				`{"@type":"/Invalid","minPrice":{"denom":"stake","amount":"0"}}`,
			},
			errMsg: "unable to resolve type URL /Invalid",
		},
		{
			desc: "invalid duration arg",
			args: []string{
				"1",
				"xyz",
				"24",
				"",
			},
			errMsg: "time: missing unit in duration \"24\"",
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			args := append(tc.args, flags...)
			_, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.CmdPublishAuction(), args)

			if tc.errMsg != "" {
				require.EqualError(t, err, tc.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
