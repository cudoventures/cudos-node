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

func (s *KeeperTestSuite) TestBurnSingleCoin() {
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
			name:        "Burn non-acudos coins, should be burned",
			denom:       sdk.DefaultBondDenom,
			initCoin:    sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(200))),
			coinsToBurn: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(200))),
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
			name:        "Burn acudos coins, should be sent to community pool",
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
		{
			name:        "Burn cudosAdmin coins, should be ignored",
			denom:       params.BondDenom,
			initCoin:    sdk.NewCoins(sdk.NewCoin(params.AdminTokenDenom, sdk.NewInt(1000))),
			coinsToBurn: sdk.NewCoins(sdk.NewCoin(params.AdminTokenDenom, sdk.NewInt(1000))),
			expectedTotalSupply: func(supplyBeforeBurn sdk.Coins, burnCoin sdk.Coins) sdk.Coins {
				// expected total supply is the same as before when burning acudos coins
				return supplyBeforeBurn
			},
			expectedTotalDistributed: func(distributedBeforeBurn sdk.Coin, burnCoin sdk.Coin) sdk.Coin {
				// expected total distributed is distributedBeforeBur
				return distributedBeforeBurn
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			// Mint coin
			s.Require().NoError(s.App.BankKeeper.MintCoins(s.Ctx, cudoMinttypes.ModuleName, tc.initCoin))
			supplyAfterInflation, _, _ := s.App.BankKeeper.GetPaginatedTotalSupply(s.Ctx, &query.PageRequest{})
			distributedAfterInflation := s.App.BankKeeper.GetBalance(s.Ctx, s.App.DistrKeeper.GetDistributionAccount(s.Ctx).GetAddress(), tc.denom)

			// Given coin to burn
			s.Require().NoError(s.App.BankKeeper.SendCoinsFromModuleToModule(s.Ctx, cudoMinttypes.ModuleName, govtypes.ModuleName, tc.coinsToBurn))

			// When burn happens
			err := s.App.BankKeeper.BurnCoins(s.Ctx, govtypes.ModuleName, tc.coinsToBurn)
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

func (s *KeeperTestSuite) TestBurnMultipleCoins() {
	randomDenom := "randomDenom"
	testCases := []struct {
		name                     string
		denoms                   []string
		initCoins                []sdk.Coins
		coinsToBurn              []sdk.Coins
		expectedBurnCoins        sdk.Coins
		expectedToCommunityPool  sdk.Coins
		expectedTotalSupply      func(suppliesBeforeBurn sdk.Coins, burnCoins sdk.Coins) sdk.Coins
		expectedTotalDistributed func(distributedBeforeBurn sdk.Coins, burnCoins sdk.Coins) sdk.Coins
	}{
		// Test cases for burning non-acudos coins
		{
			name:   "Burn non-acudos coins, should be burned",
			denoms: []string{sdk.DefaultBondDenom, randomDenom},
			initCoins: []sdk.Coins{
				sdk.NewCoins(sdk.NewCoin(randomDenom, sdk.NewInt(100))),
				sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(200))),
			},
			coinsToBurn: []sdk.Coins{
				sdk.NewCoins(sdk.NewCoin(randomDenom, sdk.NewInt(100))),
				sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(200))),
			},
			expectedBurnCoins: sdk.Coins{
				sdk.NewCoin(randomDenom, sdk.NewInt(100)),
				sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(200)),
			},
			expectedTotalSupply: func(suppliesBeforeBurn sdk.Coins, burnCoins sdk.Coins) sdk.Coins {
				// expected total supply is distributedBeforeBurn - burnCoin
				return suppliesBeforeBurn.Sub(burnCoins)
			},
			expectedTotalDistributed: func(distributedBeforeBurn sdk.Coins, distributedCoins sdk.Coins) sdk.Coins {
				// expected total distributed is the same as before when burning non-acudos coins
				return distributedBeforeBurn
			},
		},
		// Test cases for burning acudos and cudosAdmin coins
		{
			name:   "Burn acudos and cudosAdmin coins, acudos should be sent to community pool and cudosAdmin should be ignored",
			denoms: []string{params.BondDenom, params.AdminTokenDenom},
			initCoins: []sdk.Coins{
				sdk.NewCoins(sdk.NewCoin(params.BondDenom, sdk.NewInt(100))),
				sdk.NewCoins(sdk.NewCoin(params.AdminTokenDenom, sdk.NewInt(1000))),
			},
			coinsToBurn: []sdk.Coins{
				sdk.NewCoins(sdk.NewCoin(params.BondDenom, sdk.NewInt(100))),
				sdk.NewCoins(sdk.NewCoin(params.AdminTokenDenom, sdk.NewInt(1000))),
			},
			expectedToCommunityPool: sdk.Coins{
				sdk.NewCoin(params.BondDenom, sdk.NewInt(100)),
			},
			expectedTotalSupply: func(suppliesBeforeBurn sdk.Coins, burnCoins sdk.Coins) sdk.Coins {
				return suppliesBeforeBurn
			},
			expectedTotalDistributed: func(distributedBeforeBurn sdk.Coins, distributedCoins sdk.Coins) sdk.Coins {
				return distributedBeforeBurn.Add(distributedCoins...)
			},
		},
		// Testcase for burning random and cudosAdmin coins
		{
			name:   "Burn random and cudosAdmin coins, cudosAdmin should be ignored",
			denoms: []string{randomDenom, params.AdminTokenDenom},
			initCoins: []sdk.Coins{
				sdk.NewCoins(sdk.NewCoin(params.AdminTokenDenom, sdk.NewInt(1000))),
				sdk.NewCoins(sdk.NewCoin(randomDenom, sdk.NewInt(100))),
			},
			coinsToBurn: []sdk.Coins{
				sdk.NewCoins(sdk.NewCoin(params.AdminTokenDenom, sdk.NewInt(1000))),
				sdk.NewCoins(sdk.NewCoin(randomDenom, sdk.NewInt(100))),
			},
			expectedBurnCoins: sdk.Coins{
				sdk.NewCoin(randomDenom, sdk.NewInt(100)),
			},
			expectedTotalSupply: func(suppliesBeforeBurn sdk.Coins, burnCoins sdk.Coins) sdk.Coins {
				return suppliesBeforeBurn.Sub(burnCoins)
			},
			expectedTotalDistributed: func(distributedBeforeBurn sdk.Coins, distributedCoins sdk.Coins) sdk.Coins {
				return distributedBeforeBurn
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			s.SetupTest()
			// Mint coins
			for i := range tc.initCoins {
				s.Require().NoError(s.App.BankKeeper.MintCoins(s.Ctx, cudoMinttypes.ModuleName, tc.initCoins[i]))
			}
			supplyAfterInflation, _, _ := s.App.BankKeeper.GetPaginatedTotalSupply(s.Ctx, &query.PageRequest{})
			distributedAfterInflation := s.App.BankKeeper.GetAllBalances(s.Ctx, s.App.DistrKeeper.GetDistributionAccount(s.Ctx).GetAddress())

			// Given coins to burn
			for i := range tc.coinsToBurn {
				s.Require().NoError(s.App.BankKeeper.SendCoinsFromModuleToModule(s.Ctx, cudoMinttypes.ModuleName, govtypes.ModuleName, tc.coinsToBurn[i]))
			}

			// When burn happens
			for i := range tc.coinsToBurn {
				err := s.App.BankKeeper.BurnCoins(s.Ctx, govtypes.ModuleName, tc.coinsToBurn[i])
				s.Require().NoError(err)
			}

			// Then
			supplyAfterBurn, _, err := s.App.BankKeeper.GetPaginatedTotalSupply(s.Ctx, &query.PageRequest{})
			s.Require().NoError(err)
			totalDistributedCoinsAfterBurn := s.App.BankKeeper.GetAllBalances(s.Ctx, s.App.DistrKeeper.GetDistributionAccount(s.Ctx).GetAddress())
			s.Require().Equal(tc.expectedTotalSupply(supplyAfterInflation, tc.expectedBurnCoins), supplyAfterBurn)
			s.Require().Equal(tc.expectedTotalDistributed(distributedAfterInflation, tc.expectedToCommunityPool), totalDistributedCoinsAfterBurn)
		})
	}
}
