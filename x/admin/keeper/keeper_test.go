package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/CudoVentures/cudos-node/simapp"
	"github.com/CudoVentures/cudos-node/testutil/sample"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	cudominttypes "github.com/CudoVentures/cudos-node/x/cudoMint/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app simapp.CudosSimApp
	ctx sdk.Context
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	s.T().Log("setting up keeper test suite")
	app := simapp.Setup(s.T(), false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	s.app, s.ctx = *app, ctx
}

func (s *KeeperTestSuite) TestAdminSpendCommunityPool() {
	communityPoolReceiver := sample.AccAddress()
	bondDenom := s.app.StakingKeeper.BondDenom(s.ctx)

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
			s.app.BankKeeper.MintCoins(s.ctx, cudominttypes.ModuleName, sdk.NewCoins(tc.feeAmount))
			s.app.BankKeeper.SendCoinsFromModuleToModule(s.ctx, cudominttypes.ModuleName, distrtypes.ModuleName, sdk.NewCoins(tc.feeAmount))
			s.app.DistrKeeper.SetFeePool(s.ctx, newFeePool)

			err := s.app.AdminKeeper.AdminDistributeFromFeePool(s.ctx, sdk.NewCoins(tc.withdrawAmount), sdk.AccAddress(communityPoolReceiver))
			if tc.expError {
				s.Error(err)
			} else {
				s.NoError(err)
				s.Require().Equal(tc.withdrawAmount.Amount, s.app.BankKeeper.GetBalance(s.ctx, sdk.AccAddress(communityPoolReceiver), bondDenom).Amount)
			}
		})
	}
}
