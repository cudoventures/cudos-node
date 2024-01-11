package keeper_test

import (
	appparams "github.com/CudoVentures/cudos-node/app/params"
	adminkeeper "github.com/CudoVentures/cudos-node/x/admin/keeper"
	admintypes "github.com/CudoVentures/cudos-node/x/admin/types"
	cudominttypes "github.com/CudoVentures/cudos-node/x/cudoMint/types"
	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

func (s *KeeperTestSuite) TestMsgAdminSpendCommunityPool() {
	communityPoolReceiver := s.KeeperTestHelper.RandomAccountAddresses(1)[0]
	bondDenom := appparams.BondDenom
	// random admin address
	adminAddr := s.KeeperTestHelper.RandomAccountAddresses(1)[0]
	// random address
	notAdminAddr := s.KeeperTestHelper.RandomAccountAddresses(1)[0]

	testCases := []struct {
		name           string
		sender         string
		withdrawAmount sdk.Coin
		cpAmount       sdk.Coin // Community pool amount
		expectedError  bool
	}{
		{
			name:           "available community pool: 40, withdraw: 30",
			sender:         adminAddr.String(),
			withdrawAmount: sdk.NewCoin(bondDenom, sdk.NewInt(30)),
			cpAmount:       sdk.NewCoin(bondDenom, sdk.NewInt(40)),
			expectedError:  false,
		},
		{
			name:           "available community pool: 40, withdraw: 50",
			sender:         adminAddr.String(),
			withdrawAmount: sdk.NewCoin(bondDenom, sdk.NewInt(50)),
			cpAmount:       sdk.NewCoin(bondDenom, sdk.NewInt(40)),
			expectedError:  true,
		},
		{
			name:           "available community pool: 40, withdraw: 40",
			sender:         adminAddr.String(),
			withdrawAmount: sdk.NewCoin(bondDenom, sdk.NewInt(40)),
			cpAmount:       sdk.NewCoin(bondDenom, sdk.NewInt(40)),
			expectedError:  false,
		},
		{
			name:           "available community pool: 40, withdraw: 0",
			sender:         adminAddr.String(),
			withdrawAmount: sdk.NewCoin(bondDenom, sdk.NewInt(0)),
			cpAmount:       sdk.NewCoin(bondDenom, sdk.NewInt(40)),
			expectedError:  false,
		},
		{
			name:           "available community pool: 0, withdraw: 0",
			sender:         adminAddr.String(),
			withdrawAmount: sdk.NewCoin(bondDenom, sdk.NewInt(0)),
			cpAmount:       sdk.NewCoin(bondDenom, sdk.NewInt(0)),
			expectedError:  false,
		},
		{
			name:           "available community pool: 30, withdraw 20,  not admin",
			sender:         notAdminAddr.String(),
			withdrawAmount: sdk.NewCoin(bondDenom, sdk.NewInt(20)),
			cpAmount:       sdk.NewCoin(bondDenom, sdk.NewInt(30)),
			expectedError:  true,
		},
		{
			name:           "available community pool: 30, withdraw 40,  not admin",
			sender:         notAdminAddr.String(),
			withdrawAmount: sdk.NewCoin(bondDenom, sdk.NewInt(40)),
			cpAmount:       sdk.NewCoin(bondDenom, sdk.NewInt(30)),
			expectedError:  true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			s.Setup(s.T(), chainID)
			adminCoins := sdk.NewCoins(sdk.NewCoin(bondDenom, sdk.NewInt(100)), sdk.NewCoin(appparams.AdminTokenDenom, sdk.NewInt(100)))

			// mint and send to admin
			s.App.BankKeeper.MintCoins(s.KeeperTestHelper.Ctx, cudominttypes.ModuleName, adminCoins)
			s.App.BankKeeper.SendCoinsFromModuleToAccount(s.KeeperTestHelper.Ctx, cudominttypes.ModuleName, adminAddr, adminCoins)

			msgServer := adminkeeper.NewMsgServerImpl(s.AdminKeeper)
			newFeePool := distrtypes.FeePool{
				CommunityPool: sdk.NewDecCoinsFromCoins(tc.cpAmount),
			}
			s.App.DistrKeeper.SetFeePool(s.KeeperTestHelper.Ctx, newFeePool)

			// Transfer token to the community pool
			s.App.BankKeeper.MintCoins(s.KeeperTestHelper.Ctx, cudominttypes.ModuleName, sdk.NewCoins(tc.cpAmount))
			s.App.BankKeeper.SendCoinsFromModuleToModule(s.KeeperTestHelper.Ctx, cudominttypes.ModuleName, distrtypes.ModuleName, sdk.NewCoins(tc.cpAmount))

			// Construct and send message

			msg := admintypes.MsgAdminSpendCommunityPool{
				Initiator: tc.sender,
				ToAddress: communityPoolReceiver.String(),
				Coins:     sdk.NewCoins(tc.withdrawAmount),
			}

			_, err := msgServer.AdminSpendCommunityPool(types.WrapSDKContext(s.KeeperTestHelper.Ctx), &msg)

			if tc.expectedError {
				s.Error(err)
			} else {
				s.NoError(err)
				afterCPBalance := s.App.DistrKeeper.GetFeePool(s.KeeperTestHelper.Ctx).CommunityPool.AmountOf(bondDenom)

				s.Suite.Equal(tc.withdrawAmount, s.App.BankKeeper.GetBalance(s.KeeperTestHelper.Ctx, communityPoolReceiver, bondDenom))
				s.Suite.Equal(tc.cpAmount.Sub(tc.withdrawAmount).Amount.ToDec(), afterCPBalance)
			}
		})
	}
}
