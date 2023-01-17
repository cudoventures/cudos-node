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
			desc: "valid english auction",
			args: []string{
				"xyz",
				"1",
				"25h",
				`{"@type":"/cudoventures.cudosnode.marketplace.EnglishAuction","minPrice":{"denom":"acudos","amount":"1"}}`,
			},
		},
		{
			desc: "valid dutch auction",
			args: []string{
				"xyz",
				"1",
				"25h",
				`{"@type":"/cudoventures.cudosnode.marketplace.DutchAuction","startPrice":{"denom":"acudos","amount":"10"},"minPrice":{"denom":"acudos","amount":"1"}}`,
			},
		},
		{
			desc: "invalid denom id",
			args: []string{
				"123",
				"1",
				"25h",
				`{"@type":"/cudoventures.cudosnode.marketplace.EnglishAuction","minPrice":{"denom":"acudos","amount":"1"}}`,
			},
			wantErr: nfttypes.ErrInvalidDenom,
		},
		{
			desc: "invalid token id",
			args: []string{
				"xyz",
				"invalid",
				"25h",
				`{"@type":"/cudoventures.cudosnode.marketplace.EnglishAuction","minPrice":{"denom":"acudos","amount":"1"}}`,
			},
			wantErr: nfttypes.ErrInvalidTokenID,
		},
		{
			desc: "invalid auction type",
			args: []string{
				"xyz",
				"1",
				"25h",
				`{"@type":"/Invalid","minPrice":{"denom":"acudos","amount":"0"}}`,
			},
			wantErr: sdkerrors.ErrInvalidType,
		},
		{
			desc: "invalid duration",
			args: []string{
				"xyz",
				"1",
				"24",
				"",
			},
			wantErr: types.ErrInvalidAuctionDuration,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			ctx := val.ClientCtx
			args := append(tc.args, flags...)
			_, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdPublishAuction(), args)
			require.ErrorIs(t, err, tc.wantErr)
		})
	}
}
