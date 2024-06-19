package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/CudoVentures/cudos-node/x/messaging/types"
    "github.com/CudoVentures/cudos-node/x/messaging/keeper"
    keepertest "github.com/CudoVentures/cudos-node/testutil/keeper"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.MessagingKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
