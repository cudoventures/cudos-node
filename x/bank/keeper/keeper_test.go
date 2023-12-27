package keeper_test

import (
	"testing"

	"github.com/CudoVentures/cudos-node/app/apptesting"
	cudoMinttypes "github.com/CudoVentures/cudos-node/x/cudoMint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/suite"
)

const initialPower = int64(100)

var burnerAcc = authtypes.NewEmptyModuleAccount(authtypes.Burner, authtypes.Burner)

// The default power validators are initialized to have within tests
var initTokens = sdk.TokensFromConsensusPower(initialPower, sdk.DefaultPowerReduction)
var initCoins = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initTokens))

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
	s.SetupTest()

	s.App.AccountKeeper.SetModuleAccount(s.Ctx, burnerAcc)
	s.Require().NoError(s.App.BankKeeper.MintCoins(s.Ctx, cudoMinttypes.ModuleName, initCoins))

	s.Require().NoError(
		s.App.BankKeeper.SendCoinsFromModuleToModule(s.Ctx, cudoMinttypes.ModuleName, govtypes.ModuleName, initCoins),
	)

	dbBefore := s.App.BankKeeper.GetBalance(s.Ctx, s.App.DistrKeeper.GetDistributionAccount(s.Ctx).GetAddress(), "acudos")
	supplyAfterInflation, _, err := s.App.BankKeeper.GetPaginatedTotalSupply(s.Ctx, &query.PageRequest{})
	s.Require().NoError(err)
	err = s.App.BankKeeper.BurnCoins(s.Ctx, govtypes.ModuleName, initCoins)
	supplyAfterBurn, _, err := s.App.BankKeeper.GetPaginatedTotalSupply(s.Ctx, &query.PageRequest{})
	s.Require().Equal(supplyAfterInflation.Sub(initCoins), supplyAfterBurn)
	dbAfter := s.App.BankKeeper.GetBalance(s.Ctx, s.App.DistrKeeper.GetDistributionAccount(s.Ctx).GetAddress(), "acudos")
	s.Require().Equal(dbBefore, dbAfter, "destribution module balance shouldn't change when burning non acudos")

	acudosAmount := sdk.NewCoins(sdk.NewCoin("acudos", sdk.NewInt(100000000000000)))
	s.Require().NoError(s.App.BankKeeper.MintCoins(s.Ctx, cudoMinttypes.ModuleName, acudosAmount))
	s.Require().NoError(
		s.App.BankKeeper.SendCoinsFromModuleToModule(s.Ctx, cudoMinttypes.ModuleName, govtypes.ModuleName, acudosAmount),
	)
	dbBefore1 := s.App.BankKeeper.GetBalance(s.Ctx, s.App.DistrKeeper.GetDistributionAccount(s.Ctx).GetAddress(), "acudos")
	supplyAfterInflation1, _, err := s.App.BankKeeper.GetPaginatedTotalSupply(s.Ctx, &query.PageRequest{})
	s.Require().NoError(err)
	err = s.App.BankKeeper.BurnCoins(s.Ctx, govtypes.ModuleName, acudosAmount)
	supplyAfterBurn1, _, err := s.App.BankKeeper.GetPaginatedTotalSupply(s.Ctx, &query.PageRequest{})
	s.Require().Equal(supplyAfterInflation1, supplyAfterBurn1)
	dbAfter1 := s.App.BankKeeper.GetBalance(s.Ctx, s.App.DistrKeeper.GetDistributionAccount(s.Ctx).GetAddress(), "acudos")
	s.Require().Equal(dbBefore1.Add(acudosAmount[0]), dbAfter1, "destribution module balance shouldn't change when burning non acudos")
}
