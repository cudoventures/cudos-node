package app_test

import (
	"testing"

	"github.com/CudoVentures/cudos-node/app/apptesting"
	cudoMinttypes "github.com/CudoVentures/cudos-node/x/cudoMint/types"
	gravitytypes "github.com/althea-net/cosmos-gravity-bridge/module/x/gravity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

type UpgradeTestSuite struct {
	apptesting.KeeperTestHelper
}

func TestUpgradeTestSuite(t *testing.T) {
	suite.Run(t, new(UpgradeTestSuite))
}

func (s *UpgradeTestSuite) TestUpgrade_V11_To_V12() {
	s.Setup(s.T(), "cudos-app")

	gravityAddress := s.App.AccountKeeper.GetModuleAddress(gravitytypes.ModuleName)
	// == SETUP ==
	// Set the initial balance for the gravity module account
	initialCoin := sdk.NewCoin("acudos", sdk.NewInt(1000000000000))
	s.App.BankKeeper.MintCoins(s.Ctx, cudoMinttypes.ModuleName, sdk.NewCoins(initialCoin))
	s.App.BankKeeper.SendCoinsFromModuleToModule(s.Ctx, cudoMinttypes.ModuleName, gravitytypes.ModuleName, sdk.NewCoins(initialCoin))

	balanceBefore := s.App.BankKeeper.GetAllBalances(s.Ctx, gravityAddress)
	s.Require().Equal(sdk.NewCoins(initialCoin), balanceBefore)
	// == UPGRADE ==
	upgradeHeight := int64(5)
	s.ConfirmUpgradeSucceeded("v1.2", upgradeHeight)

	// == CHECK ==

	//// Ensure the balance of the gravity module account is unchanged
	balanceAfter := s.App.BankKeeper.GetAllBalances(s.Ctx, gravityAddress)
	s.Require().Equal(balanceBefore, balanceAfter)

}

func (s *UpgradeTestSuite) TestUpgrade_V12_To_V13() {
	s.Setup(s.T(), "cudos-app")

	gravityAddress := s.App.AccountKeeper.GetModuleAddress(gravitytypes.ModuleName)
	// == SETUP ==
	// Set the initial balance for the gravity module account
	initialCoin := sdk.NewCoin("acudos", sdk.NewInt(1000000000000))
	s.App.BankKeeper.MintCoins(s.Ctx, cudoMinttypes.ModuleName, sdk.NewCoins(initialCoin))
	s.App.BankKeeper.SendCoinsFromModuleToModule(s.Ctx, cudoMinttypes.ModuleName, gravitytypes.ModuleName, sdk.NewCoins(initialCoin))

	balanceBefore := s.App.BankKeeper.GetAllBalances(s.Ctx, gravityAddress)
	s.Require().Equal(sdk.NewCoins(initialCoin), balanceBefore)
	// == UPGRADE ==
	upgradeHeight := int64(5)
	s.ConfirmUpgradeSucceeded("v1.3", upgradeHeight)

	// == CHECK ==

	//// Ensure the balance of the gravity module account is unchanged
	balanceAfter := s.App.BankKeeper.GetAllBalances(s.Ctx, gravityAddress)
	s.Require().Equal(balanceBefore, balanceAfter)
	s.Require().NotNil(s.App.GravityKeeper.GetBridgeContractAddress(s.Ctx))
	s.Require().NotEmpty(s.App.GravityKeeper.GetStaticValCosmosAddrs(s.Ctx))
}
