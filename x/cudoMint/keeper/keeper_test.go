package keeper_test

import (
	"testing"

	"github.com/CudoVentures/cudos-node/app/apptesting"
	"github.com/stretchr/testify/suite"

	"github.com/CudoVentures/cudos-node/x/cudoMint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	s.T().Log("setting up keeper test suite")
	s.Setup()
}

func (s *KeeperTestSuite) TestMintCoins() {
	s.SetupTest()
	mintCoins := sdk.NewCoins(
		sdk.NewCoin("acudos", sdk.NewInt(100000)),
		sdk.NewCoin("cudosAdmin", sdk.NewInt(1000)),
	)

	err := s.App.CudoMintKeeper.MintCoins(s.Ctx, mintCoins)
	s.NoError(err)

	s.Require().Equal(sdk.NewInt(100000), s.App.BankKeeper.GetBalance(s.Ctx, authtypes.NewModuleAddress(types.ModuleName), "acudos").Amount)
	s.Require().Equal(sdk.NewInt(1000), s.App.BankKeeper.GetBalance(s.Ctx, authtypes.NewModuleAddress(types.ModuleName), "cudosAdmin").Amount)
}

func (s *KeeperTestSuite) TestAddCollectedFees() {
	s.SetupTest()
	mintCoins := sdk.NewCoins(
		sdk.NewCoin("acudos", sdk.NewInt(100000)),
		sdk.NewCoin("cudosAdmin", sdk.NewInt(1000)),
	)

	err := s.App.CudoMintKeeper.MintCoins(s.Ctx, mintCoins)
	s.NoError(err)
	err = s.App.CudoMintKeeper.AddCollectedFees(s.Ctx, mintCoins)
	s.NoError(err)
	feeAcc := authtypes.NewModuleAddressOrBech32Address(authtypes.FeeCollectorName)
	nativeBalance := s.App.BankKeeper.GetBalance(s.Ctx, feeAcc, "acudos").Amount
	adminBalance := s.App.BankKeeper.GetBalance(s.Ctx, feeAcc, "cudosAdmin").Amount
	s.Require().Equal(sdk.NewInt(100000), nativeBalance)
	s.Require().Equal(sdk.NewInt(1000), adminBalance)
}
