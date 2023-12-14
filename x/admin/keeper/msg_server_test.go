package keeper_test

import (
	"github.com/CudoVentures/cudos-node/testutil/sample"

	adminkeeper "github.com/CudoVentures/cudos-node/x/admin/keeper"
	admintypes "github.com/CudoVentures/cudos-node/x/admin/types"
	cudominttypes "github.com/CudoVentures/cudos-node/x/cudoMint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

func (s *KeeperTestSuite) TestMsgAdminSpendCommunityPool() {
	communityPoolReceiver := sample.AccAddress()
	bondDenom := s.App.StakingKeeper.BondDenom(s.Ctx)
	adminAddr := "cosmos1qae52r3vdtl92am2klvqe9rtn3534mllsf3sfj"
	notAdminAddr := sample.AccAddress()

	testCases := []struct {
		name           string
		sender         string
		withdrawAmount sdk.Coin
		feeAmount      sdk.Coin
		expError       bool
	}{
		{
			name:           "fee 45acudos, withdraw 30acudos, success",
			sender:         adminAddr,
			withdrawAmount: sdk.NewCoin(bondDenom, sdk.NewInt(30)),
			feeAmount:      sdk.NewCoin(bondDenom, sdk.NewInt(45)),
			expError:       false,
		},
		{
			name:           "fee 30acudos, withdraw 32acudos, failed",
			sender:         adminAddr,
			withdrawAmount: sdk.NewCoin(bondDenom, sdk.NewInt(32)),
			feeAmount:      sdk.NewCoin(bondDenom, sdk.NewInt(30)),
			expError:       true,
		},
		{
			name:           "fee 45acudos, withdraw 30acudos, not admin, failed",
			sender:         notAdminAddr,
			withdrawAmount: sdk.NewCoin(bondDenom, sdk.NewInt(45)),
			feeAmount:      sdk.NewCoin(bondDenom, sdk.NewInt(30)),
			expError:       true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			s.SetupTest()
			addrAcc, _ := sdk.AccAddressFromBech32(adminAddr)
			adminCoins := sdk.NewCoins(sdk.NewCoin(bondDenom, sdk.NewInt(45)), sdk.NewCoin(admintypes.AdminDenom, sdk.OneInt()))
			s.App.BankKeeper.MintCoins(s.Ctx, cudominttypes.ModuleName, adminCoins)
			s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx, cudominttypes.ModuleName, addrAcc, adminCoins)

			msgServer := adminkeeper.NewMsgServerImpl(s.App.AdminKeeper)
			newFeePool := distrtypes.FeePool{
				CommunityPool: sdk.NewDecCoinsFromCoins(tc.feeAmount),
			}
			s.App.BankKeeper.MintCoins(s.Ctx, cudominttypes.ModuleName, sdk.NewCoins(tc.feeAmount))
			s.App.BankKeeper.SendCoinsFromModuleToModule(s.Ctx, cudominttypes.ModuleName, distrtypes.ModuleName, sdk.NewCoins(tc.feeAmount))
			s.App.DistrKeeper.SetFeePool(s.Ctx, newFeePool)

			msgAdminSpendCommunityPool := admintypes.MsgAdminSpendCommunityPool{
				Initiator: tc.sender,
				ToAddress: communityPoolReceiver,
				Coins:     sdk.NewCoins(tc.withdrawAmount),
			}
			_, err := msgServer.AdminSpendCommunityPool(s.Ctx, &msgAdminSpendCommunityPool)
			comAcc, _ := sdk.AccAddressFromBech32(communityPoolReceiver)
			if tc.expError {
				s.Error(err)
			} else {
				s.NoError(err)
				s.Equal(tc.withdrawAmount, s.App.BankKeeper.GetBalance(s.Ctx, comAcc, bondDenom))
			}
		})
	}
}
