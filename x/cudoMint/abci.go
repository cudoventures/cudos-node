package cudoMint

import (
	"cudos.org/cudos-node/x/cudoMint/keeper"
	"cudos.org/cudos-node/x/cudoMint/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

// Minting based on the formula f(t)=358 - 53 * t + 1.8 * t^2, where t is number of years passed since the release = 150mil, 7 sec - 150mils/(7)
var (
	// based on the assumption that we have 1 block per 5 seconds
	denom               = "acudos" // Hardcoded to the acudos currency. Its not changeable, because some of the math depends on the size of this denomination
	totalDays           = sdk.NewInt(3652) // Hardcoded to 10 years
	FinalNormTimePassed = sdk.NewDec(10)
	zeroPointSix        = sdk.MustNewDecFromStr("0.6")
	twentySixPointFive  = sdk.MustNewDecFromStr("26.5")
)

// Normalize block height incrementation
func normalizeBlockHeightInc(blocksPerDay sdk.Int) sdk.Dec {
	totalBlocks := blocksPerDay.Mul(totalDays)
	return (sdk.NewDec(1).QuoInt(totalBlocks)).Mul(FinalNormTimePassed)
}

// Integral of f(t) is 0,6 * t^3  - 26.5 * t^2 + 358 * t
func calculateIntegral(t sdk.Dec) sdk.Dec {
	return (zeroPointSix.Mul(t.Power(3))).Sub(twentySixPointFive.Mul(t.Power(2))).Add(sdk.NewDec(358).Mul(t))
}

func calculateMintedCoins(minter types.Minter, increment sdk.Dec) sdk.Dec {
	prevStep := calculateIntegral(minter.NormTimePassed)
	nextStep := calculateIntegral(minter.NormTimePassed.Add(increment))
	return (nextStep.Sub(prevStep)).Mul(sdk.NewDec(10).Power(24)) // formula calculates in mil of cudos + converting to acudos
}

// BeginBlocker mints new tokens for the previous block.
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	minter := k.GetMinter(ctx)
	params := k.GetParams(ctx)

	if minter.NormTimePassed.GT(FinalNormTimePassed) {
		return
	}
	incr := normalizeBlockHeightInc(params.BlocksPerDay)
	mintAmountDec := calculateMintedCoins(minter, incr)
	mintAmountInt := mintAmountDec.TruncateInt()
	mintedCoin := sdk.NewCoin(denom, mintAmountInt)
	mintedCoins := sdk.NewCoins(mintedCoin)
	err := k.MintCoins(ctx, mintedCoins)
	if err != nil {
		panic(err)
	}
	minter.NormTimePassed = minter.NormTimePassed.Add(incr)
	minter.MintRemainder = mintAmountDec.Sub(mintAmountInt.ToDec())
	k.SetMinter(ctx, minter)

	// send the minted coins to the fee collector account
	err = k.AddCollectedFees(ctx, mintedCoins)
	if err != nil {
		panic(err)
	}

	if mintedCoin.Amount.IsInt64() {
		defer telemetry.ModuleSetGauge(types.ModuleName, float32(mintedCoin.Amount.Int64()), "minted_tokens")
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeMintedDenom, denom),
			sdk.NewAttribute(types.AttributeMintedTokens, mintAmountInt.String()),
		),
	)
}
