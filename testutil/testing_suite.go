package testutil

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/testutil"

	"github.com/cosmos/gogoproto/proto"

	"github.com/cosmos/cosmos-sdk/client"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
)

func RunNetwork(t *testing.T, cfg network.Config) (*network.Network, error) {
	network, err := network.New(t, t.TempDir(), cfg)
	if err != nil {
		return nil, err
	}

	if _, err := network.WaitForHeight(3); err != nil {
		return nil, err
	}

	return network, nil
}

func WaitForBlock() {
	time.Sleep(time.Duration(3) * time.Second)
}

func QueryJustBroadcastedTx(clientCtx client.Context, bz testutil.BufferWriter) (*sdk.TxResponse, error) {
	respType := proto.Message(&sdk.TxResponse{})
	if err := clientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType); err != nil {
		return nil, err
	}

	txResp := respType.(*sdk.TxResponse)
	bz, err := clitestutil.ExecTestCLICmd(clientCtx, authcmd.QueryTxCmd(), []string{
		txResp.TxHash,
	})
	if err != nil {
		return nil, err
	}

	if err = clientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType); err != nil {
		return nil, err
	}

	return respType.(*sdk.TxResponse), nil

}
