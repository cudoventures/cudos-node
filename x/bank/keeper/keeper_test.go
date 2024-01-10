package keeper_test

import (
	"testing"

	"github.com/CudoVentures/cudos-node/app/apptesting"
	"github.com/CudoVentures/cudos-node/app/params"
	cudoMinttypes "github.com/CudoVentures/cudos-node/x/cudoMint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/suite"
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	s.T().Log("setting up bank keeper test suite")
	s.Setup(s.T(), apptesting.SimAppChainID)
}

func (s *KeeperTestSuite) TestBurnCoins() {
	testCases := []struct {
		name                     string                                                  // Name of the test case
		denom                    string                                                  // Denom of the coin to burn
		initCoin                 sdk.Coins                                               // Initial coins in the account
		coinsToBurn              sdk.Coins                                               // Coins to burn
		expectedTotalSupply      func(supplyBeforeBurn, burnCoin sdk.Coins) sdk.Coins    // Expected total supply after burning
		expectedTotalDistributed func(distributedBeforeBurn, burnCoin sdk.Coin) sdk.Coin // Expected total distributed coins after burning
	}{
		// Test cases for burning non-acudos coins
		{
			name:        "Burn non-acudos coins",
			denom:       sdk.DefaultBondDenom,
			initCoin:    sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
			coinsToBurn: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
			expectedTotalSupply: func(supplyBeforeBurn sdk.Coins, burnCoin sdk.Coins) sdk.Coins {
				// expected total supply is distributedBeforeBurn - burnCoin
				return supplyBeforeBurn.Sub(burnCoin)
			},
			expectedTotalDistributed: func(distributedBeforeBurn, burnCoin sdk.Coin) sdk.Coin {
				// expected total distributed is the same as before when burning non-acudos coins
				return distributedBeforeBurn
			},
		},
		// Test cases for burning acudos coins
		{
			name:        "Burn acudos coins",
			denom:       params.BondDenom,
			initCoin:    sdk.NewCoins(sdk.NewCoin(params.BondDenom, sdk.NewInt(100))),
			coinsToBurn: sdk.NewCoins(sdk.NewCoin(params.BondDenom, sdk.NewInt(100))),
			expectedTotalSupply: func(supplyBeforeBurn sdk.Coins, burnCoin sdk.Coins) sdk.Coins {
				// expected total supply is the same as before when burning acudos coins
				return supplyBeforeBurn
			},
			expectedTotalDistributed: func(distributedBeforeBurn sdk.Coin, burnCoin sdk.Coin) sdk.Coin {
				// expected total distributed is distributedBeforeBurn + burnCoin
				return distributedBeforeBurn.Add(burnCoin)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			// Mint coin
			s.Require().NoError(s.App.BankKeeper.MintCoins(s.Ctx, cudoMinttypes.ModuleName, tc.initCoin))
			supplyAfterInflation, _, err := s.App.BankKeeper.GetPaginatedTotalSupply(s.Ctx, &query.PageRequest{})
			distributedAfterInflation := s.App.BankKeeper.GetBalance(s.Ctx, s.App.DistrKeeper.GetDistributionAccount(s.Ctx).GetAddress(), tc.denom)

			// Given coin to burn
			s.Require().NoError(s.App.BankKeeper.SendCoinsFromModuleToModule(s.Ctx, cudoMinttypes.ModuleName, govtypes.ModuleName, tc.coinsToBurn))

			// When burn happens
			err = s.App.BankKeeper.BurnCoins(s.Ctx, govtypes.ModuleName, tc.coinsToBurn)
			s.Require().NoError(err)

			// Then
			supplyAfterBurn, _, err := s.App.BankKeeper.GetPaginatedTotalSupply(s.Ctx, &query.PageRequest{})
			s.Require().NoError(err)
			totalDistributedCoinsAfterBurn := s.App.BankKeeper.GetBalance(s.Ctx, s.App.DistrKeeper.GetDistributionAccount(s.Ctx).GetAddress(), tc.denom)
			s.Require().Equal(tc.expectedTotalSupply(supplyAfterInflation, tc.coinsToBurn), supplyAfterBurn)
			s.Require().Equal(tc.expectedTotalDistributed(distributedAfterInflation, tc.coinsToBurn[0]), totalDistributedCoinsAfterBurn)
		})
	}
}
