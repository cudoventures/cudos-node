package keeper_test

import (
	"testing"

	"github.com/CudoVentures/cudos-node/app/apptesting"
	"github.com/stretchr/testify/suite"

	"github.com/CudoVentures/cudos-node/testutil/sample"

	cudominttypes "github.com/CudoVentures/cudos-node/x/cudoMint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	s.T().Log("setting up admin keeper test suite")
	s.Setup()
}

func (s *KeeperTestSuite) TestAdminSpendCommunityPool() {
	communityPoolReceiver := sample.AccAddress()
	bondDenom := s.App.StakingKeeper.BondDenom(s.Ctx)

	testCases := []struct {
		name           string
		withdrawAmount sdk.Coin
		feeAmount      sdk.Coin
		expError       bool
	}{
		{
			name:           "fee 45acudos, withdraw 30acudos, success",
			withdrawAmount: sdk.NewCoin(bondDenom, sdk.NewInt(30)),
			feeAmount:      sdk.NewCoin(bondDenom, sdk.NewInt(45)),
			expError:       false,
		},
		{
			name:           "fee 30acudos, withdraw 32acudos, failed",
			withdrawAmount: sdk.NewCoin(bondDenom, sdk.NewInt(32)),
			feeAmount:      sdk.NewCoin(bondDenom, sdk.NewInt(30)),
			expError:       true,
		},
		{
			name:           "fee 0, withdraw 0, success",
			withdrawAmount: sdk.NewCoin(bondDenom, sdk.NewInt(0)),
			feeAmount:      sdk.NewCoin(bondDenom, sdk.NewInt(0)),
			expError:       false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			s.SetupTest()
			newFeePool := distrtypes.FeePool{
				CommunityPool: sdk.NewDecCoinsFromCoins(tc.feeAmount),
			}
			err := s.App.BankKeeper.MintCoins(s.Ctx, cudominttypes.ModuleName, sdk.NewCoins(tc.feeAmount))
			s.Require().NoError(err)
			err = s.App.BankKeeper.SendCoinsFromModuleToModule(s.Ctx, cudominttypes.ModuleName, distrtypes.ModuleName, sdk.NewCoins(tc.feeAmount))
			s.Require().NoError(err)
			s.App.DistrKeeper.SetFeePool(s.Ctx, newFeePool)
			err = s.App.AdminKeeper.AdminDistributeFromFeePool(s.Ctx, sdk.NewCoins(tc.withdrawAmount), sdk.AccAddress(communityPoolReceiver))
			if tc.expError {
				s.Error(err)
			} else {
				s.NoError(err)
				s.Require().Equal(tc.withdrawAmount.Amount, s.App.BankKeeper.GetBalance(s.Ctx, sdk.AccAddress(communityPoolReceiver), bondDenom).Amount)
			}
		})
	}
}
