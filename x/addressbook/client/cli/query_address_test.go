package cli_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/CudoVentures/cudos-node/testutil/network"
	"github.com/CudoVentures/cudos-node/testutil/nullify"
	"github.com/CudoVentures/cudos-node/testutil/sample"
	"github.com/CudoVentures/cudos-node/x/addressbook/client/cli"
	"github.com/CudoVentures/cudos-node/x/addressbook/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func networkWithAddressObjects(t *testing.T, n int) (*network.Network, []types.Address) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	for i := 0; i < n; i++ {
		address := types.Address{
			Creator: sample.AccAddress(),
			Network: "BTC",
			Label:   fmt.Sprintf("%d@testdenom", i),
		}
		nullify.Fill(&address)
		state.AddressList = append(state.AddressList, address)
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), state.AddressList
}

func TestShowAddress(t *testing.T) {
	net, objs := networkWithAddressObjects(t, 2)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	for _, tc := range []struct {
		desc    string
		creator string
		network string
		label   string

		args []string
		err  error
		obj  types.Address
	}{
		{
			desc:    "found",
			creator: objs[0].Creator,
			network: objs[0].Network,
			label:   objs[0].Label,

			args: common,
			obj:  objs[0],
		},
		{
			desc:    "not found",
			creator: sample.AccAddress(),
			network: "BTC",
			label:   "0@testdenom",

			args: common,
			err:  status.Error(codes.NotFound, "not found"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				tc.creator,
				tc.network,
				tc.label,
			}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowAddress(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryGetAddressResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.Address)
				require.Equal(t,
					nullify.Fill(&tc.obj),
					nullify.Fill(&resp.Address),
				)
			}
		})
	}
}

func TestListAddress(t *testing.T) {
	net, objs := networkWithAddressObjects(t, 5)

	ctx := net.Validators[0].ClientCtx
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
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(objs); i += step {
			args := request(nil, uint64(i), uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListAddress(), args)
			require.NoError(t, err)
			var resp types.QueryAllAddressResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.Address), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.Address),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListAddress(), args)
			require.NoError(t, err)
			var resp types.QueryAllAddressResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.Address), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.Address),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListAddress(), args)
		require.NoError(t, err)
		var resp types.QueryAllAddressResponse
		require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(objs),
			nullify.Fill(resp.Address),
		)
	})
}
