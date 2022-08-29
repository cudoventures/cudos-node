package keeper_test

import (
	"testing"

	testkeeper "github.com/CudoVentures/cudos-node/testutil/keeper"
	"github.com/CudoVentures/cudos-node/testutil/sample"
	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestDistributeRoyalties(t *testing.T) {
	kp, ctx := testkeeper.MarketplaceKeeper(t)

	seller := sample.AccAddress()
	err := kp.DistributeRoyalties(ctx, sdk.NewCoin("acudos", sdk.NewInt(100)), seller, []types.Royalty{})
	require.NoError(t, err)
}
