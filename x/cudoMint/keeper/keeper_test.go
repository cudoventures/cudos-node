package keeper_test

import (
	"testing"

	"github.com/CudoVentures/cudos-node/app/apptesting"
	appparams "github.com/CudoVentures/cudos-node/app/params"
	"github.com/CudoVentures/cudos-node/x/cudoMint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/suite"
)

var chainID = "cudos-app"

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// TestMintCoins tests if the coins are minted and added to the cudosMint module account
func (s *KeeperTestSuite) TestMintCoins() {
	s.Setup(s.T(), chainID)
	mintCoins := sdk.NewCoins(
		sdk.NewCoin("acudos", sdk.NewInt(100000)),
		sdk.NewCoin("cudosAdmin", sdk.NewInt(1000)),
	)

	err := s.CudoMintKeeper.MintCoins(s.Ctx, mintCoins)
	s.NoError(err)

	s.Require().Equal(sdk.NewInt(100000), s.App.BankKeeper.GetBalance(s.Ctx, authtypes.NewModuleAddress(types.ModuleName), "acudos").Amount)
	s.Require().Equal(sdk.NewInt(1000), s.App.BankKeeper.GetBalance(s.Ctx, authtypes.NewModuleAddress(types.ModuleName), "cudosAdmin").Amount)
}

func (s *KeeperTestSuite) TestMintNothing() {
	s.Setup(s.T(), chainID)
	mintCoins := sdk.NewCoins()

	err := s.CudoMintKeeper.MintCoins(s.Ctx, mintCoins)
	s.NoError(err)
}

// TestAddCollectedFees tests if the collected fees are transferred from the cudosMint module account to the fee collector account
func (s *KeeperTestSuite) TestAddCollectedFees() {
	s.Setup(s.T(), chainID)
	mintCoins := sdk.NewCoins(
		sdk.NewCoin(appparams.BondDenom, sdk.NewInt(100000)),
		sdk.NewCoin(appparams.AdminTokenDenom, sdk.NewInt(1000)),
	)

	err := s.CudoMintKeeper.MintCoins(s.Ctx, mintCoins)
	s.NoError(err)

	err = s.CudoMintKeeper.AddCollectedFees(s.Ctx, mintCoins)
	s.NoError(err)

	feeAcc := authtypes.NewModuleAddress(authtypes.FeeCollectorName)

	// Check if the coins are added to the fee collector account
	s.Require().Equal(sdk.NewInt(100000), s.App.BankKeeper.GetBalance(s.Ctx, feeAcc, appparams.BondDenom).Amount)
	s.Require().Equal(sdk.NewInt(1000), s.App.BankKeeper.GetBalance(s.Ctx, feeAcc, appparams.AdminTokenDenom).Amount)
}
