package keeper_test

import (
	"testing"

	"github.com/CudoVentures/cudos-node/app/apptesting"
	appparams "github.com/CudoVentures/cudos-node/app/params"
	cudominttypes "github.com/CudoVentures/cudos-node/x/cudoMint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/stretchr/testify/suite"
)

const (
	chainID = "cudos-app"
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) TestAdminDistributeFromFeePool() {
	bondDenom := appparams.BondDenom
	poolReceiver := s.KeeperTestHelper.RandomAccountAddresses(1)[0]
	testCases := []struct {
		name           string
		withdrawAmount sdk.Coin
		feeAmount      sdk.Coin
		expectedError  bool
	}{
		{
			name:           "available fee pool: 40, withdraw: 30",
			withdrawAmount: sdk.NewCoin(bondDenom, sdk.NewInt(30)),
			feeAmount:      sdk.NewCoin(bondDenom, sdk.NewInt(40)),
			expectedError:  false,
		},
		{
			name:           "available fee pool: 40, withdraw: 50",
			withdrawAmount: sdk.NewCoin(bondDenom, sdk.NewInt(50)),
			feeAmount:      sdk.NewCoin(bondDenom, sdk.NewInt(40)),
			expectedError:  true,
		},
		{
			name:           "available fee pool: 40, withdraw: 40",
			withdrawAmount: sdk.NewCoin(bondDenom, sdk.NewInt(40)),
			feeAmount:      sdk.NewCoin(bondDenom, sdk.NewInt(40)),
			expectedError:  false,
		},
		{
			name:           "available fee pool: 40, withdraw: 0",
			withdrawAmount: sdk.NewCoin(bondDenom, sdk.NewInt(0)),
			feeAmount:      sdk.NewCoin(bondDenom, sdk.NewInt(40)),
			expectedError:  false,
		},
		{
			name:           "available fee pool: 0, withdraw: 0",
			withdrawAmount: sdk.NewCoin(bondDenom, sdk.NewInt(0)),
			feeAmount:      sdk.NewCoin(bondDenom, sdk.NewInt(0)),
			expectedError:  false,
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.Setup(s.T(), chainID)
			// Set up fee pool and send coins to fee pool through bank keeper
			newFeePool := distrtypes.FeePool{
				CommunityPool: sdk.NewDecCoins(sdk.NewDecCoinFromCoin(tc.feeAmount)),
			}
			s.App.BankKeeper.MintCoins(s.Ctx, cudominttypes.ModuleName, sdk.NewCoins(tc.feeAmount))
			s.App.BankKeeper.SendCoinsFromModuleToModule(s.Ctx, cudominttypes.ModuleName, distrtypes.ModuleName, sdk.NewCoins(tc.feeAmount))
			s.App.DistrKeeper.SetFeePool(s.Ctx, newFeePool)

			// Distribute from fee pool
			err := s.AdminKeeper.AdminDistributeFromFeePool(s.Ctx, sdk.NewCoins(tc.withdrawAmount), poolReceiver)

			if tc.expectedError {
				s.Error(err)
			} else {
				s.NoError(err)
				// Assert that the received amount and the remaining amount in the fee pool are correc
				afterCPBalance := s.App.DistrKeeper.GetFeePool(s.Ctx).CommunityPool.AmountOf(bondDenom)
				s.Require().Equal(tc.withdrawAmount.Amount, s.App.BankKeeper.GetBalance(s.Ctx, poolReceiver, bondDenom).Amount)
				s.Require().Equal(sdk.NewDecFromInt(tc.feeAmount.Amount.Sub(tc.withdrawAmount.Amount)), afterCPBalance)
			}
		})

	}
}
