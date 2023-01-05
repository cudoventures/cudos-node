package cli_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/CudoVentures/cudos-node/simapp"
	"github.com/CudoVentures/cudos-node/testutil/network"
	"github.com/CudoVentures/cudos-node/x/marketplace/client/cli"
	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	nfttypes "github.com/CudoVentures/cudos-node/x/nft/types"
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
		desc    string
		args    []string
		wantErr error
	}{
		{
			desc: "valid",
			args: []string{
				"1",
				"xyz",
				"25h",
				`{"@type":"/cudoventures.cudosnode.marketplace.EnglishAuction","minPrice":{"denom":"stake","amount":"1"}}`,
			},
		},
		// todo dutch auction
		{
			desc: "invalid denom id",
			args: []string{
				"1",
				"123",
				"25h",
				`{"@type":"/cudoventures.cudosnode.marketplace.EnglishAuction","minPrice":{"denom":"stake","amount":"1"}}`,
			},
			wantErr: nfttypes.ErrInvalidNFT,
		},
		{
			desc: "invalid token id",
			args: []string{
				"invalid",
				"xyz",
				"25h",
				`{"@type":"/cudoventures.cudosnode.marketplace.EnglishAuction","minPrice":{"denom":"stake","amount":"1"}}`,
			},
			wantErr: nfttypes.ErrInvalidNFT,
		},
		{
			desc: "invalid auction type",
			args: []string{
				"1",
				"xyz",
				"25h",
				`{"@type":"/Invalid","minPrice":{"denom":"stake","amount":"0"}}`,
			},
			wantErr: sdkerrors.ErrInvalidType,
		},
		{
			desc: "invalid duration",
			args: []string{
				"1",
				"xyz",
				"24",
				"",
			},
			wantErr: types.ErrInvalidAuctionDuration,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			args := append(tc.args, flags...)
			_, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.CmdPublishAuction(), args)
			require.ErrorIs(t, err, tc.wantErr)
		})
	}
}
