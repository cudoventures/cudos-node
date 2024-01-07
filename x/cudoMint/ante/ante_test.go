package ante_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/CudoVentures/cudos-node/app"
	"github.com/CudoVentures/cudos-node/app/apptesting"
	appparams "github.com/CudoVentures/cudos-node/app/params"
	simulation "github.com/CudoVentures/cudos-node/x/simulation"
	"github.com/cosmos/cosmos-sdk/client"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// AnteTestSuite is a test suite to be used with ante handler tests.
type AnteTestSuite struct {
	apptesting.KeeperTestHelper
	clientCtx client.Context
	txBuilder client.TxBuilder
}

// SetupTest setups a new test, with new app, context, and anteHandler.
func (suite *AnteTestSuite) SetupTest() {
	suite.Setup(suite.T(), apptesting.SimAppChainID)

	// Set up TxConfig.
	encodingConfig := suite.SetupEncoding()
	// suite.SetupCudoMint()

	suite.clientCtx = client.Context{}.
		WithTxConfig(encodingConfig.TxConfig)

}

func (suite *AnteTestSuite) SetupEncoding() appparams.EncodingConfig {
	encodingConfig := app.MakeEncodingConfig()
	// We're using TestMsg encoding in some tests, so register it here.
	encodingConfig.Amino.RegisterConcrete(&testdata.TestMsg{}, "testdata.TestMsg", nil)

	return encodingConfig
}

// CreateTestTx is a helper function to create a tx given multiple inputs.
func (suite *AnteTestSuite) CreateTestTx(privs []cryptotypes.PrivKey, accNums []uint64, accSeqs []uint64, chainID string) (xauthsigning.Tx, error) {
	// First round: we gather all the signer infos. We use the "set empty
	// signature" hack to do that.
	var sigsV2 []signing.SignatureV2
	for i, priv := range privs {
		sigV2 := signing.SignatureV2{
			PubKey: priv.PubKey(),
			Data: &signing.SingleSignatureData{
				SignMode:  suite.clientCtx.TxConfig.SignModeHandler().DefaultMode(),
				Signature: nil,
			},
			Sequence: accSeqs[i],
		}

		sigsV2 = append(sigsV2, sigV2)
	}
	err := suite.txBuilder.SetSignatures(sigsV2...)
	if err != nil {
		return nil, err
	}

	// Second round: all signer infos are set, so each signer can sign.
	sigsV2 = []signing.SignatureV2{}
	for i, priv := range privs {
		signerData := xauthsigning.SignerData{
			ChainID:       chainID,
			AccountNumber: accNums[i],
			Sequence:      accSeqs[i],
		}
		sigV2, err := tx.SignWithPrivKey(
			suite.clientCtx.TxConfig.SignModeHandler().DefaultMode(), signerData,
			suite.txBuilder, priv, suite.clientCtx.TxConfig, accSeqs[i])
		if err != nil {
			return nil, err
		}

		sigsV2 = append(sigsV2, sigV2)
	}
	err = suite.txBuilder.SetSignatures(sigsV2...)
	if err != nil {
		return nil, err
	}

	return suite.txBuilder.GetTx(), nil
}

func (suite *AnteTestSuite) CreateValidator(tokens *big.Int) (cryptotypes.PrivKey, cryptotypes.PubKey, stakingtypes.Validator, authtypes.AccountI) {
	suite.Ctx = suite.Ctx.WithBlockHeight(suite.Ctx.BlockHeight() + 1)
	suite.App.BeginBlock(abci.RequestBeginBlock{Header: suite.Ctx.BlockHeader()})
	// at the end of commit, deliverTx is set to nil, which is why we need to get newest instance of deliverTx ctx here after committing
	// update ctx to new deliver tx context
	suite.Ctx = suite.App.NewContext(false, suite.Ctx.BlockHeader())

	priv, pub, addr := testdata.KeyTestPubAddr()
	consKey, valPub, _ := suite.Ed25519PubAddr()
	valAddr := sdk.ValAddress(addr)

	sendCoins := sdk.NewCoins(sdk.NewCoin(appparams.BondDenom, sdk.NewIntFromBigInt(tokens.Mul(tokens, big.NewInt(2)))))
	suite.FundAcc(addr, sendCoins)

	// set account in account keeper
	account := authtypes.NewBaseAccountWithAddress(addr)
	account.SetPubKey(pub)
	account.SetAccountNumber(0)
	account.SetSequence(0)
	// create sim account to pass into deliver tx context
	// TODO: refactor, why need simAccount?
	simAccount := simtypes.Account{
		PrivKey: priv,
		Address: addr,
		PubKey:  pub,
		ConsKey: consKey,
	}
	// set account in deliver tx context
	suite.App.AccountKeeper.SetAccount(suite.Ctx, account)

	commissionRates := stakingtypes.NewCommissionRates(
		sdk.NewDecWithPrec(1, 2), sdk.NewDecWithPrec(1, 0),
		sdk.NewDecWithPrec(1, 0),
	)

	delegationCoin := sdk.NewCoin(appparams.BondDenom, sdk.NewIntFromBigInt(tokens))
	desc := stakingtypes.NewDescription("moniker", "", "", "", "")

	msgCreateValidator, err := stakingtypes.NewMsgCreateValidator(
		valAddr,
		valPub,
		delegationCoin,
		desc,
		commissionRates,
		sdk.NewIntFromBigInt(tokens),
	)
	suite.Require().NoError(err)

	// deliver Tx
	_, _, err = simulation.GenAndDeliverTx(suite.App.BaseApp, suite.clientCtx.TxConfig, msgCreateValidator, sdk.NewCoins(), suite.Ctx, simAccount, suite.App.AccountKeeper, stakingtypes.ModuleName)
	suite.Require().NoError(err)
	fmt.Println("deliver tx success")
	// turn block for validator updates
	suite.App.EndBlock(abci.RequestEndBlock{Height: suite.Ctx.BlockHeight()})
	suite.App.Commit()

	retval, found := suite.App.StakingKeeper.GetValidator(suite.Ctx, valAddr)
	suite.Require().Equal(true, found)

	updatedAccount := suite.App.AccountKeeper.GetAccount(suite.Ctx, addr)

	return priv, pub, retval, updatedAccount
}

func (suite *AnteTestSuite) TestAnte_HappyPathCreateValidator() {
	const MinSelfDelegation = "2000000000000000000000000"
	suite.SetupTest() // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()
	suite.txBuilder.SetGasLimit(400_000)

	msd, _ := new(big.Int).SetString(MinSelfDelegation, 10)
	_, _, _, _ = suite.CreateValidator(msd)

	suite.Ctx = suite.Ctx.WithBlockHeight(suite.Ctx.BlockHeight() + 1)
	suite.App.BeginBlock(abci.RequestBeginBlock{Header: suite.Ctx.BlockHeader()})
	// update ctx to new deliver tx context
	suite.Ctx = suite.App.NewContext(false, suite.Ctx.BlockHeader())
	suite.App.EndBlock(abci.RequestEndBlock{Height: suite.Ctx.BlockHeight()})
	suite.App.Commit()
	// tx, err := suite.CreateTestTx([]cryptotypes.PrivKey{priv1}, []uint64{acc.GetAccountNumber()}, []uint64{acc.GetSequence()}, suite.Ctx.ChainID())
	// suite.Require().NoError(err)
	// dyncomm := suite.App.DyncommKeeper.CalculateDynCommission(suite.Ctx, val1)
	// invalidtarget := dyncomm.Mul(sdk.NewDecWithPrec(9, 1))

	// // invalid tx fails, not updating account sequence in account keeper
	// editmsg := stakingtypes.NewMsgEditValidator(
	// 	val1.GetOperator(),
	// 	val1.Description, &invalidtarget, &val1.MinSelfDelegation,
	// )

	// err := suite.txBuilder.SetMsgs(editmsg)
	// suite.Require().NoError(err)

	// // due to submitting a create validator tx before, thus account sequence is now 1
	// for i := 0; i < 5; i++ {
	// 	suite.Ctx = suite.Ctx.WithBlockHeight(suite.Ctx.BlockHeight() + 1)
	// 	suite.App.BeginBlock(abci.RequestBeginBlock{Header: suite.Ctx.BlockHeader()})
	// 	// update ctx to new deliver tx context
	// 	suite.Ctx = suite.App.NewContext(false, suite.Ctx.BlockHeader())

	// tx, err := suite.CreateTestTx([]cryptotypes.PrivKey{priv1}, []uint64{acc.GetAccountNumber()}, []uint64{acc.GetSequence()}, suite.Ctx.ChainID())
	// suite.Require().NoError(err)

	// 	_, checkRes, err := suite.App.SimCheck(suite.clientCtx.TxConfig.TxEncoder(), tx)
	// 	fmt.Printf("check response: %+v, error = %v \n", checkRes, err)
	// 	suite.Ctx = suite.Ctx.WithIsCheckTx(false)
	// 	_, deliverRes, err := suite.App.SimDeliver(suite.clientCtx.TxConfig.TxEncoder(), tx)
	// 	fmt.Printf("deliver response: %+v, error = %v \n", deliverRes, err)
	// 	suite.App.EndBlock(abci.RequestEndBlock{Height: suite.Ctx.BlockHeight()})
	// 	suite.App.Commit()

	// 	// check and update account keeper
	// 	acc = suite.App.AccountKeeper.GetAccount(suite.CheckCtx, acc.GetAddress())
	// 	checkSeq := acc.GetSequence()
	// 	// checkSeq not updated when checkTx fails
	// 	suite.Require().Equal(uint64(1), checkSeq)
	// 	acc = suite.App.AccountKeeper.GetAccount(suite.Ctx, acc.GetAddress())
	// 	deliverSeq := acc.GetSequence()
	// 	// deliverSeq not updated when deliverTx fails
	// 	suite.Require().Equal(uint64(1), deliverSeq)
	// }
}

func TestAnteTestSuite(t *testing.T) {
	suite.Run(t, new(AnteTestSuite))
}
