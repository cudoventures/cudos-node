package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/CudoVentures/cudos-node/simapp"
	"github.com/CudoVentures/cudos-node/x/cudoMint/types"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
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

func (s *KeeperTestSuite) TestMintCoins() {
	s.SetupTest()
	mintCoins := sdk.NewCoins(
		sdk.NewCoin("acudos", sdk.NewInt(100000)),
		sdk.NewCoin("cudosAdmin", sdk.NewInt(1000)),
	)

	err := s.app.CudoMintKeeper.MintCoins(s.ctx, mintCoins)
	s.NoError(err)

	s.Require().Equal(sdk.NewInt(100000), s.app.BankKeeper.GetBalance(s.ctx, authtypes.NewModuleAddress(types.ModuleName), "acudos").Amount)
	s.Require().Equal(sdk.NewInt(1000), s.app.BankKeeper.GetBalance(s.ctx, authtypes.NewModuleAddress(types.ModuleName), "cudosAdmin").Amount)
}

func (s *KeeperTestSuite) TestAddCollectedFees() {
	s.SetupTest()
	mintCoins := sdk.NewCoins(
		sdk.NewCoin("acudos", sdk.NewInt(100000)),
		sdk.NewCoin("cudosAdmin", sdk.NewInt(1000)),
	)

	err := s.app.CudoMintKeeper.MintCoins(s.ctx, mintCoins)
	s.NoError(err)
	err = s.app.CudoMintKeeper.AddCollectedFees(s.ctx, mintCoins)
	s.NoError(err)
	feeAcc := authtypes.NewModuleAddressOrBech32Address(authtypes.FeeCollectorName)
	nativeBalance := s.app.BankKeeper.GetBalance(s.ctx, feeAcc, "acudos").Amount
	adminBalance := s.app.BankKeeper.GetBalance(s.ctx, feeAcc, "cudosAdmin").Amount
	s.Require().Equal(sdk.NewInt(100000), nativeBalance)
	s.Require().Equal(sdk.NewInt(1000), adminBalance)
}