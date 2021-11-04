package cudoMint

import (
	"cudos.org/cudos-node/x/cudoMint/types"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"
)

func TestCalculateMintedCoins(t *testing.T) {
	blocksPerDay := sdk.NewInt(100)
	minter := types.NewMinter(sdk.NewDec(0), sdk.NewDec(0))
	totalBlocks := blocksPerDay.Mul(totalDays).Int64()
	incr := normalizeBlockHeightInc(blocksPerDay)
	mintedCoins := sdk.NewDec(0)
	for i := int64(0); i < totalBlocks; i++ {
		coins := calculateMintedCoins(minter, incr)
		mintedCoins = mintedCoins.Add(coins)
		minter.NormTimePassed = minter.NormTimePassed.Add(incr)
		if i%10000 == 0 {
			blocksPerDay = blocksPerDay.Add(sdk.NewInt(24))
			fmt.Println(fmt.Sprintf("%v: Printed %v, Got total %v mil, norm time passed %v, blocks per day %v", i, coins, mintedCoins, minter.NormTimePassed, blocksPerDay))
		}
	}
	fmt.Println(fmt.Sprintf("Got total %v mil, norm time passed %v", mintedCoins, minter.NormTimePassed))
	expectedCoins := sdk.MustNewDecFromStr("1530000000000000000000000000.000000000000000000")
	oneCudo := sdk.NewDec(10).Power(18)
	if mintedCoins.Add(oneCudo).LT(expectedCoins) || mintedCoins.GT(expectedCoins) {
		t.Errorf("Got unexpected amount of coins %v; wanted %v +/- epsilon", mintedCoins, expectedCoins)
	}
	if minter.NormTimePassed.Add(incr).LT(FinalNormTimePassed) || minter.NormTimePassed.GT(FinalNormTimePassed) {
		t.Errorf("Got unexpected normalized time passed %v; wanted %v +/- epsilon", minter.NormTimePassed, FinalNormTimePassed)
	}

}
