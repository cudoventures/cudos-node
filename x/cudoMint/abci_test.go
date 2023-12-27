package cudoMint_test

import (
	"testing"
)

func TestCalculateMintedCoins(t *testing.T) {
	// app := simapp.Setup(false)

	// ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	// app.CudoMintKeeper.SetParams(ctx, types.NewParams(sdk.NewInt(10)))
	// totalBlocks := int64(100000)
	// for height := int64(1); height <= totalBlocks; height++ {
	// 	ctx = ctx.WithBlockHeight(height)
	// 	cudoMint.BeginBlocker(ctx, app.CudoMintKeeper)
	// }

	// expectedNormTimePassed, _ := sdk.NewDecFromStr("10.0003")
	// require.True(t, app.CudoMintKeeper.GetMinter(ctx).NormTimePassed.LT(expectedNormTimePassed))

	// expectedSupply, _ := sdk.NewIntFromString("1530000000000000000000000000")
	// require.Equal(t, expectedSupply.String(), app.BankKeeper.GetSupply(ctx, "acudos").Amount.String())
}
