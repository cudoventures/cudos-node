package decorators_test

import (
	"testing"

	"github.com/CudoVentures/cudos-node/app"
	"github.com/CudoVentures/cudos-node/app/apptesting"
	appparams "github.com/CudoVentures/cudos-node/app/params"
	cometbftdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/client"
	sims "github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/suite"
)

type AnteTestSuite struct {
	apptesting.KeeperTestHelper
	clientCtx client.Context
	txBuilder client.TxBuilder
}

// returns context and app with params set on account keeper
func createTestApp(isCheckTx bool, tempDir string) (*app.App, sdk.Context) {
	db := cometbftdb.NewMemDB()
	encCdc := app.MakeEncodingConfig()

	cudosApp := app.New(
		log.NewNopLogger(),
		db,
		nil,
		true,
		map[int64]bool{},
		tempDir,
		5,
		encCdc,
		sims.EmptyAppOptions{},
	)
	ctx := cudosApp.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	cudosApp.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())

	return cudosApp, ctx
}

// SetupTest setups a new test, with new app, context, and anteHandler.
func (suite *AnteTestSuite) SetupTest(isCheckTx bool) {
	tempDir := suite.T().TempDir()
	suite.KeeperTestHelper.App, suite.KeeperTestHelper.Ctx = createTestApp(isCheckTx, tempDir)
	suite.KeeperTestHelper.Ctx = suite.KeeperTestHelper.Ctx.WithBlockHeight(1)

	// Set up TxConfig.
	encodingConfig := suite.SetupEncoding()

	suite.clientCtx = client.Context{}.
		WithTxConfig(encodingConfig.TxConfig)
}

func (suite *AnteTestSuite) SetupEncoding() appparams.EncodingConfig {
	encodingConfig := app.MakeEncodingConfig()
	// We're using TestMsg encoding in some tests, so register it here.
	encodingConfig.Amino.RegisterConcrete(&testdata.TestMsg{}, "testdata.TestMsg", nil)
	testdata.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	return encodingConfig
}

func TestAnteTestSuite(t *testing.T) {
	suite.Run(t, new(AnteTestSuite))
}
