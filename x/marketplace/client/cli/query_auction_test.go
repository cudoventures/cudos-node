package cli_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	tmcli "github.com/tendermint/tendermint/libs/cli"

	"github.com/CudoVentures/cudos-node/simapp"
	"github.com/CudoVentures/cudos-node/testutil/network"
	"github.com/CudoVentures/cudos-node/testutil/sample"
	"github.com/CudoVentures/cudos-node/x/marketplace/client/cli"
	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	date   = time.Date(2023, 1, 15, 15, 00, 00, 000, time.UTC)
	amount = sdk.NewCoin("acudos", sdk.NewInt(100))
	addr1  = sample.AccAddress()
)

func networkWithAuctionObjects(t *testing.T, n int) (*network.Network, []types.Auction) {
	cfg := simapp.NewConfig()
	cfg.NumValidators = 1

	state := types.GenesisState{}
	err := cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state)
	require.NoError(t, err)

	auctions := make([]types.Auction, n)
	for i := 0; i < n; i++ {
		var a types.Auction
		if i%2 == 0 {
			a = types.NewEnglishAuction(
				addr1,
				"asd",
				"1",
				amount,
				date,
				date.Add(time.Hour*24),
			)
		} else {
			a = types.NewDutchAuction(
				addr1,
				"asd",
				"1",
				amount,
				amount.SubAmount(sdk.NewInt(50)),
				date,
				date.Add(time.Hour*24),
			)
		}
		a.SetId(uint64(i))
		auctions[i] = a
	}
	auctionsAny, err := types.PackAuctions(auctions)
	require.NoError(t, err)
	state.AuctionList = auctionsAny
	state.AuctionCount = uint64(len(auctionsAny))

	bz, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = bz

	net := network.New(t, cfg)
	_, err = net.WaitForHeight(3)
	require.NoError(t, err)

	return net, auctions
}

func TestShowAuction(t *testing.T) {
	net, auctions := networkWithAuctionObjects(t, 2)
	ctx := net.Validators[0].ClientCtx

	for _, tc := range []struct {
		desc    string
		id      string
		wantErr error
	}{
		{
			desc: "valid",
			id:   strconv.FormatUint(auctions[0].GetId(), 10),
		},
		{
			desc:    "not existing auction",
			id:      "111",
			wantErr: types.ErrAuctionNotFound,
		},
		{
			desc:    "invalid auction id",
			id:      "invalid",
			wantErr: types.ErrInvalidAuctionId,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{tc.id}
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowAuction(), args)
			require.ErrorIs(t, err, tc.wantErr)

			if tc.wantErr != nil {
				return
			}

			require.NoError(t, err)
			var resp types.QueryGetAuctionResponse
			err = net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp)
			require.NoError(t, err)
			require.NotNil(t, resp.Auction)

			haveAuction, err := types.UnpackAuction(resp.Auction)
			require.NoError(t, err)
			require.Equal(t, auctions[0], haveAuction)

			args = []string{strconv.FormatInt(int64(auctions[1].GetId()), 10)}
			out, err = clitestutil.ExecTestCLICmd(ctx, cli.CmdShowAuction(), args)
			require.NoError(t, err)

			var resp2 types.QueryGetAuctionResponse
			err = net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp2)
			require.NoError(t, err)
			require.NotNil(t, resp2.Auction)

			haveAuction, err = types.UnpackAuction(resp2.Auction)
			require.NoError(t, err)
			require.Equal(t, auctions[1], haveAuction)
		})
	}
}

func TestListAuction(t *testing.T) {
	net, auctions := networkWithAuctionObjects(t, 5)
	ctx := net.Validators[0].ClientCtx
	request := func(next []byte, offset, limit uint64, total bool) []string {
		args := []string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)}
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
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(auctions); i += step {
			args := request(nil, uint64(i), uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListAuction(), args)
			require.NoError(t, err)

			var resp types.QueryAllAuctionResponse
			err = net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp)
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Auctions), step)

			haveAuctions, err := types.UnpackAuctions(resp.Auctions)
			require.NoError(t, err)
			require.Subset(t, auctions, haveAuctions)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(auctions); i += step {
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListAuction(), args)
			require.NoError(t, err)

			var resp types.QueryAllAuctionResponse
			err = net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp)
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Auctions), step)

			haveAuctions, err := types.UnpackAuctions(resp.Auctions)
			require.NoError(t, err)
			require.Subset(t, auctions, haveAuctions)

			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(auctions)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListAuction(), args)
		require.NoError(t, err)

		var resp types.QueryAllAuctionResponse
		err = net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp)
		require.NoError(t, err)
		require.NoError(t, err)
		require.Equal(t, len(auctions), int(resp.Pagination.Total))

		haveAuctions, err := types.UnpackAuctions(resp.Auctions)
		require.NoError(t, err)
		require.Subset(t, auctions, haveAuctions)
	})
	t.Run("InvalidPageRequest", func(t *testing.T) {
		args := request(nil, 1, uint64(len(auctions)), true)
		args = append(args, fmt.Sprintf("--%s=%d", flags.FlagPage, 2))
		_, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListAuction(), args)
		require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
	})
}
