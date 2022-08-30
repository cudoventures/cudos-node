package keeper_test

import (
	"testing"

	testkeeper "github.com/CudoVentures/cudos-node/testutil/keeper"
	"github.com/CudoVentures/cudos-node/testutil/sample"
	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func genAddresses(count int) []string {
	addresses := []string{}
	for i := 0; i < count; i++ {
		addresses = append(addresses, sample.AccAddress())
	}
	return addresses
}

func TestDistributeRoyalties(t *testing.T) {
	addresses := genAddresses(4)

	firstRoyaltyPrcent, err := sdk.NewDecFromStr("22.9")
	require.NoError(t, err)

	secondRoyaltyPercent, err := sdk.NewDecFromStr("0.01")
	require.NoError(t, err)

	thirdRoyaltyPercent, err := sdk.NewDecFromStr("30")
	require.NoError(t, err)

	kp, bankKeeper, ctx := testkeeper.MarketplaceKeeper(t)

	price := int64(10000)
	require.NoError(t, bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin("acudos", sdk.NewInt(price)))))

	err = kp.DistributeRoyalties(ctx, sdk.NewCoin("acudos", sdk.NewInt(price)), addresses[0], []types.Royalty{
		{
			Address: addresses[1],
			Percent: firstRoyaltyPrcent,
		},
		{
			Address: addresses[2],
			Percent: secondRoyaltyPercent,
		},
		{
			Address: addresses[3],
			Percent: thirdRoyaltyPercent,
		},
	})
	require.NoError(t, err)

	accountsBalances := bankKeeper.GetAccountsBalances(ctx)
	accBalancesMap := make(map[string]int64)

	distributedRoyalties := int64(0)

	for _, accBalance := range accountsBalances {
		accBalancesMap[accBalance.Address] = accBalance.Coins[0].Amount.Int64()

		distributedRoyalties += accBalance.Coins[0].Amount.Int64()
	}

	require.Equal(t, accBalancesMap[addresses[0]], int64(4709))
	require.Equal(t, accBalancesMap[addresses[1]], int64(2290))
	require.Equal(t, accBalancesMap[addresses[2]], int64(1))
	require.Equal(t, accBalancesMap[addresses[3]], int64(3000))

	require.Equal(t, distributedRoyalties, price)
}
