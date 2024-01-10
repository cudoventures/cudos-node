package decorators_test

import (
	"errors"

	"github.com/CudoVentures/cudos-node/app/decorators"
	"github.com/CudoVentures/cudos-node/app/params"
	cudoMinttypes "github.com/CudoVentures/cudos-node/x/cudoMint/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
)

// TestCrisisOnlyAdmin tests that the decorator properly checks for admin tokens
func (suite AnteTestSuite) TestCrisisOnlyAdmin() {
	testCases := []struct {
		name            string    // Name of the test case
		denom           string    // Denom of the coin to burn
		mintCoin        sdk.Coins // Initial coins to mint
		sendCoin        sdk.Coins // Initial coins in the account
		withAccountAddr func() (priv cryptotypes.PrivKey, addr sdk.AccAddress, sender string)
		expectedErr     error
	}{
		{
			name:     "success when enough coins in admin account",
			denom:    params.AdminTokenDenom,
			mintCoin: sdk.NewCoins(sdk.NewCoin(params.AdminTokenDenom, sdk.NewInt(100))),
			sendCoin: sdk.NewCoins(sdk.NewCoin(params.AdminTokenDenom, sdk.NewInt(100))),
			withAccountAddr: func() (priv cryptotypes.PrivKey, addr sdk.AccAddress, sender string) {
				priv, _, addr = testdata.KeyTestPubAddr()
				return priv, addr, addr.String()
			},
		},
		{
			name:     "failed when sender in msg is invalid",
			denom:    params.AdminTokenDenom,
			mintCoin: sdk.NewCoins(sdk.NewCoin(params.AdminTokenDenom, sdk.NewInt(100))),
			sendCoin: sdk.NewCoins(sdk.NewCoin(params.AdminTokenDenom, sdk.NewInt(100))),
			withAccountAddr: func() (priv cryptotypes.PrivKey, addr sdk.AccAddress, sender string) {
				priv, _, addr = testdata.KeyTestPubAddr()
				return priv, addr, "invalid sender"
			},
			expectedErr: errors.New("decoding bech32 failed: invalid character in string: ' '"),
		},
		{
			name:     "failed when no coins in admin account",
			denom:    params.AdminTokenDenom,
			mintCoin: sdk.NewCoins(sdk.NewCoin(params.AdminTokenDenom, sdk.NewInt(100))),
			sendCoin: sdk.NewCoins(sdk.NewCoin(params.AdminTokenDenom, sdk.NewInt(0))),
			withAccountAddr: func() (priv cryptotypes.PrivKey, addr sdk.AccAddress, sender string) {
				priv, _, addr = testdata.KeyTestPubAddr()
				return priv, addr, addr.String()
			},
			expectedErr: errors.New("sender has no admin tokens"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// setup
			suite.SetupTest(true)
			suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

			// Get address of the account
			priv, addr, sender := tc.withAccountAddr()

			// Mint coins
			suite.Require().NoError(suite.app.BankKeeper.MintCoins(suite.ctx, cudoMinttypes.ModuleName, tc.mintCoin))

			// Send coins to the account
			suite.Require().NoError(
				suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, cudoMinttypes.ModuleName, addr, tc.sendCoin),
			)

			// build and sign the transaction
			decorator := decorators.NewOnlyAdminVerifyInvariantDecorator(suite.app.BankKeeper)
			antehandler := sdk.ChainAnteDecorators(decorator)
			suite.Require().NoError(suite.txBuilder.SetMsgs(
				&crisistypes.MsgVerifyInvariant{
					Sender:              sender,
					InvariantModuleName: cudoMinttypes.ModuleName,
					InvariantRoute:      cudoMinttypes.RouterKey,
				},
			))

			privs, accNums, accSeqs := []cryptotypes.PrivKey{priv}, []uint64{0, 1}, []uint64{0, 0}
			tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
			suite.Require().NoError(err)

			// When
			_, err = antehandler(suite.ctx, tx, false)

			// Then
			if tc.expectedErr != nil {
				suite.Require().Equal(tc.expectedErr.Error(), err.Error())
			} else {
				suite.Require().NoError(err)
			}

		})
	}
}
