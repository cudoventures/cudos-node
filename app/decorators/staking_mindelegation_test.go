package decorators_test

import (
	"fmt"

	"github.com/CudoVentures/cudos-node/app/apptesting"
	"github.com/CudoVentures/cudos-node/app/decorators"
	"github.com/CudoVentures/cudos-node/app/params"
	cudoMinttypes "github.com/CudoVentures/cudos-node/x/cudoMint/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (suite *AnteTestSuite) TestStakingMin() {
	msd, _ := sdk.NewIntFromString(decorators.MinSelfDelegation)
	testCases := []struct {
		name              string    // Name of the test case
		denom             string    // Denom of the coin in this test
		mintCoin          sdk.Coins // Initial coins to mint
		sendCoin          sdk.Coins // Initial coins in the account
		minSelfDelegator  string
		withDelegatorAddr func() (priv cryptotypes.PrivKey, addr sdk.AccAddress)
		withValidatorAddr func() (priv cryptotypes.PrivKey, addr sdk.AccAddress)
		expectedErr       error
	}{
		{
			name:             "success when min self delegation is met",
			denom:            params.AdminTokenDenom,
			mintCoin:         sdk.NewCoins(sdk.NewCoin(params.AdminTokenDenom, sdk.NewInt(100))),
			sendCoin:         sdk.NewCoins(sdk.NewCoin(params.AdminTokenDenom, sdk.NewInt(100))),
			minSelfDelegator: decorators.MinSelfDelegation,
			withDelegatorAddr: func() (priv cryptotypes.PrivKey, addr sdk.AccAddress) {
				priv, _, addr = testdata.KeyTestPubAddr()
				return priv, addr
			},
			withValidatorAddr: func() (priv cryptotypes.PrivKey, addr sdk.AccAddress) {
				priv, _, addr = testdata.KeyTestPubAddr()
				return priv, addr
			},
		},
		{
			name:             "failed when min self delegation is not met",
			denom:            params.AdminTokenDenom,
			mintCoin:         sdk.NewCoins(sdk.NewCoin(params.AdminTokenDenom, sdk.NewInt(100))),
			sendCoin:         sdk.NewCoins(sdk.NewCoin(params.AdminTokenDenom, sdk.NewInt(100))),
			minSelfDelegator: "1",
			withDelegatorAddr: func() (priv cryptotypes.PrivKey, addr sdk.AccAddress) {
				priv, _, addr = testdata.KeyTestPubAddr()
				return priv, addr
			},
			withValidatorAddr: func() (priv cryptotypes.PrivKey, addr sdk.AccAddress) {
				priv, _, addr = testdata.KeyTestPubAddr()
				return priv, addr
			},
			expectedErr: sdkerrors.Wrap(
				sdkerrors.ErrInvalidRequest,
				fmt.Sprintf("minimum self delegation must be more than %v", msd),
			),
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// setup
			suite.SetupTest(true)
			suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

			// Get address of delegator and validator for each testcase
			privDelegator, addrDelegator := tc.withDelegatorAddr()
			privValidator, addrValidator := tc.withValidatorAddr()

			// Mint coins
			suite.Require().NoError(suite.KeeperTestHelper.App.BankKeeper.MintCoins(suite.KeeperTestHelper.Ctx, cudoMinttypes.ModuleName, tc.mintCoin))

			// Send coins to delegator
			suite.Require().NoError(
				suite.KeeperTestHelper.App.BankKeeper.SendCoinsFromModuleToAccount(suite.KeeperTestHelper.Ctx, cudoMinttypes.ModuleName, addrDelegator, tc.sendCoin),
			)

			// Build and sign a tx with a MsgCreateValidator
			decorator := decorators.NewMinSelfDelegationDecorator()
			antehandler := sdk.ChainAnteDecorators(decorator)
			minSelfDelegator, _ := sdk.NewIntFromString(tc.minSelfDelegator)
			suite.Require().NoError(suite.txBuilder.SetMsgs(
				&stakingtypes.MsgCreateValidator{
					Description:       stakingtypes.Description{},
					Commission:        stakingtypes.CommissionRates{},
					MinSelfDelegation: minSelfDelegator,
					DelegatorAddress:  string(addrDelegator),
					ValidatorAddress:  string(addrValidator),
					Pubkey:            nil,
					Value:             sdk.Coin{},
				},
			))

			privs, accNums, accSeqs := []cryptotypes.PrivKey{privDelegator, privValidator}, []uint64{0, 1}, []uint64{0, 0}
			tx, err := apptesting.CreateTestTx(privs, accNums, accSeqs, suite.KeeperTestHelper.Ctx.ChainID(), suite.clientCtx, suite.txBuilder)
			suite.Require().NoError(err)

			// When
			_, err = antehandler(suite.KeeperTestHelper.Ctx, tx, false)

			// Then
			if tc.expectedErr != nil {
				suite.Require().Equal(tc.expectedErr.Error(), err.Error())
			} else {
				suite.Require().NoError(err)
			}

		})
	}
}
