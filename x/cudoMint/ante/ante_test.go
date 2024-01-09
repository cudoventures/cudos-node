package ante_test

import (
	"fmt"
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
	"github.com/cosmos/cosmos-sdk/client"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
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
	fmt.Printf("sigsV2: %+v \n", sigsV2)

	err = suite.txBuilder.SetSignatures(sigsV2...)
	if err != nil {
		return nil, err
	}

	return suite.txBuilder.GetTx(), nil
}

func (suite *AnteTestSuite) CreateValidator(tokens sdk.Int) (cryptotypes.PrivKey, cryptotypes.PubKey, stakingtypes.Validator, authtypes.AccountI, error) {
	suite.Ctx = suite.Ctx.WithBlockHeight(suite.Ctx.BlockHeight() + 1)
	suite.App.BeginBlock(abci.RequestBeginBlock{Header: suite.Ctx.BlockHeader()})
	// at the end of commit, deliverTx is set to nil, which is why we need to get newest instance of deliverTx ctx here after committing
	// update ctx to new deliver tx context
	suite.Ctx = suite.App.NewContext(false, suite.Ctx.BlockHeader())

	priv, pub, addr := testdata.KeyTestPubAddr()
	_, valPub, _ := suite.Ed25519PubAddr()
	valAddr := sdk.ValAddress(addr)
	fmt.Println("valAddr: ", valAddr.String())
	sendCoins := sdk.NewCoins(sdk.NewCoin(appparams.BondDenom, tokens.Mul(sdk.NewInt(2))))
	suite.FundAcc(addr, sendCoins)

	// set account in account keeper
	account := authtypes.NewBaseAccountWithAddress(addr)
	account.SetPubKey(pub)
	account.SetAccountNumber(0)
	account.SetSequence(0)
	// set account in check tx and deliver tx context
	suite.App.AccountKeeper.SetAccount(suite.CheckCtx, account)
	suite.App.AccountKeeper.SetAccount(suite.Ctx, account)

	commissionRates := stakingtypes.NewCommissionRates(
		sdk.NewDecWithPrec(1, 2), sdk.NewDecWithPrec(1, 0),
		sdk.NewDecWithPrec(1, 0),
	)

	delegationCoin := sdk.NewCoin(appparams.BondDenom, tokens.Mul(sdk.NewInt(2)))
	desc := stakingtypes.NewDescription("moniker", "", "", "", "")

	msgCreateValidator, err := stakingtypes.NewMsgCreateValidator(
		valAddr,
		valPub,
		delegationCoin,
		desc,
		commissionRates,
		tokens,
	)
	if err != nil {
		return nil, nil, stakingtypes.Validator{}, nil, err
	}

	err = suite.txBuilder.SetMsgs(msgCreateValidator)
	if err != nil {
		return nil, nil, stakingtypes.Validator{}, nil, err
	}
	tx, err := suite.CreateTestTx([]cryptotypes.PrivKey{priv}, []uint64{account.GetAccountNumber()}, []uint64{account.GetSequence()}, suite.Ctx.ChainID())
	if err != nil {
		return nil, nil, stakingtypes.Validator{}, nil, err
	}

	_, _, err = suite.App.Check(suite.clientCtx.TxConfig.TxEncoder(), tx)
	if err != nil {
		return nil, nil, stakingtypes.Validator{}, nil, err
	}
	_, _, err = suite.App.Deliver(suite.clientCtx.TxConfig.TxEncoder(), tx)

	if err != nil {
		return nil, nil, stakingtypes.Validator{}, nil, err
	}

	// turn block for validator updates
	suite.App.EndBlock(abci.RequestEndBlock{Height: suite.Ctx.BlockHeight()})
	suite.App.Commit()

	retval, found := suite.App.StakingKeeper.GetValidator(suite.Ctx, valAddr)

	if !found {
		return nil, nil, stakingtypes.Validator{}, nil, fmt.Errorf("validator not found")
	}

	updatedAccount := suite.App.AccountKeeper.GetAccount(suite.Ctx, addr)

	return priv, pub, retval, updatedAccount, nil
}

func (suite *AnteTestSuite) EditValidator(tokens sdk.Int, priv cryptotypes.PrivKey) (stakingtypes.Validator, authtypes.AccountI, error) {
	suite.Ctx = suite.Ctx.WithBlockHeight(suite.Ctx.BlockHeight() + 1)
	suite.App.BeginBlock(abci.RequestBeginBlock{Header: suite.Ctx.BlockHeader()})
	// at the end of commit, deliverTx is set to nil, which is why we need to get newest instance of deliverTx ctx here after committing
	// update ctx to new deliver tx context
	suite.CheckCtx = suite.App.NewContext(true, suite.CheckCtx.BlockHeader())
	suite.Ctx = suite.App.NewContext(false, suite.Ctx.BlockHeader())
	pub := priv.PubKey()
	addr := sdk.AccAddress(pub.Address())
	valAddr := sdk.ValAddress(addr)
	account := suite.App.AccountKeeper.GetAccount(suite.Ctx, addr)
	_, found := suite.App.StakingKeeper.GetValidator(suite.Ctx, valAddr)
	if !found {
		return stakingtypes.Validator{}, nil, fmt.Errorf("validator not found")
	}

	desc := stakingtypes.NewDescription("moniker", "", "", "", "")
	msgEditValidator := stakingtypes.NewMsgEditValidator(
		valAddr,
		desc,
		nil,
		&tokens,
	)
	err := suite.txBuilder.SetMsgs(msgEditValidator)

	if err != nil {
		return stakingtypes.Validator{}, nil, err
	}

	tx, err := suite.CreateTestTx([]cryptotypes.PrivKey{priv}, []uint64{account.GetAccountNumber()}, []uint64{account.GetSequence()}, suite.Ctx.ChainID())
	if err != nil {
		return stakingtypes.Validator{}, nil, err
	}
	fmt.Println("before check tx: ", suite.App.AccountKeeper.GetAccount(suite.CheckCtx, addr).GetSequence())
	_, _, err = suite.App.Check(suite.clientCtx.TxConfig.TxEncoder(), tx)
	fmt.Println("after check tx: ", suite.App.AccountKeeper.GetAccount(suite.CheckCtx, addr).GetSequence())

	if err != nil {
		return stakingtypes.Validator{}, nil, err
	}
	fmt.Println("before deliver tx: ", suite.App.AccountKeeper.GetAccount(suite.Ctx, addr).GetSequence())
	_, _, err = suite.App.Deliver(suite.clientCtx.TxConfig.TxEncoder(), tx)
	fmt.Println("after deliver tx: ", suite.App.AccountKeeper.GetAccount(suite.Ctx, addr).GetSequence())

	if err != nil {
		return stakingtypes.Validator{}, nil, err
	}

	// turn block for validator updates
	suite.App.EndBlock(abci.RequestEndBlock{Height: suite.Ctx.BlockHeight()})
	suite.App.Commit()

	retval, found := suite.App.StakingKeeper.GetValidator(suite.Ctx, valAddr)

	if !found {
		return stakingtypes.Validator{}, nil, fmt.Errorf("error after edit validator: validator not found in store")
	}

	updatedAccount := suite.App.AccountKeeper.GetAccount(suite.Ctx, addr)

	return retval, updatedAccount, nil
}

func (suite *AnteTestSuite) TestAnte_CreateAndEditValidator() {
	MinSelfDelegation, _ := sdk.NewIntFromString("2000000000000000000000000")
	suite.SetupTest() // setup

	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()
	suite.txBuilder.SetGasLimit(400_000_000)

	priv1, _, _, _, err := suite.CreateValidator(MinSelfDelegation)
	suite.Require().NoError(err)
	// try create validator with insufficient min self delegation

	_, _, _, _, err = suite.CreateValidator(MinSelfDelegation.Sub(sdk.NewInt(1)))
	suite.Require().Error(err)
	// Revert the block height to the previous one, since the previous tx failed
	suite.Ctx = suite.Ctx.WithBlockHeight(suite.Ctx.BlockHeight() - 1)
	// try edit validator with sufficient min self delegation

	_, _, err = suite.EditValidator(MinSelfDelegation.Add(sdk.NewInt(1)), priv1)
	suite.Require().NoError(err)

	// try edit validator with insufficient min self delegation

	_, _, err = suite.EditValidator(MinSelfDelegation.Sub(sdk.NewInt(1)), priv1)
	suite.Require().Error(err)

}

func (suite *AnteTestSuite) TestAnte_CreateAndEditValidator_SeqNumber() {
	MinSelfDelegation, _ := sdk.NewIntFromString("2000000000000000000000000")

	suite.SetupTest() // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()
	suite.txBuilder.SetGasLimit(400_000_000)

	priv1, _, _, acc, err := suite.CreateValidator(MinSelfDelegation)
	suite.Require().NoError(err)

	// because successfully created validator, sequence number should be incremented
	suite.Require().Equal(uint64(1), acc.GetSequence())

	// try edit validator with sufficient min self delegation
	_, _, err = suite.EditValidator(MinSelfDelegation.Add(sdk.NewInt(1)), priv1)
	suite.Require().NoError(err)

	// because edit validator succeeded, sequence number should be incremented
	acc = suite.App.AccountKeeper.GetAccount(suite.CheckCtx, acc.GetAddress())
	checkSeq := acc.GetSequence()
	suite.Require().Equal(uint64(2), checkSeq)

	// because edit validator succeeded, sequence number should be incremented
	acc = suite.App.AccountKeeper.GetAccount(suite.Ctx, acc.GetAddress())
	delivSeq := acc.GetSequence()
	suite.Require().Equal(uint64(2), delivSeq)
}

func TestAnteTestSuite(t *testing.T) {
	suite.Run(t, new(AnteTestSuite))
}
