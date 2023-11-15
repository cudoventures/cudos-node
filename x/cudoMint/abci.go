package cudoMint

import (
	"time"

	"github.com/CudoVentures/cudos-node/x/cudoMint/keeper"
	"github.com/CudoVentures/cudos-node/x/cudoMint/types"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Minting based on the formula f(t)=358 - 53 * t + 1.8 * t^2, where t is number of years passed since the release = 150mil, 7 sec - 150mils/(7)

/*
Minting is done on steps. The step is calculated in normalizeBlockHeightInc function. All blocks have the same steps.

The step, with 17280 blocker per day, is 0.000000158462131350. This value has 18 decimal digits precision.
It is calculated by dividing 10 (~years) at total number of blocks. This could lead to an infinity number of decimals which results in precision loss by rounding up to the 18th decimal digit.

The calculation of how many tokens should be minted is done in the calculateMintedCoins function.
It returns a decimal multiplied by the 10^24.
Having in mind that the decimal that is multiplied has 18 decimal digits precision, the multiplication olaways results in a number without any decimal digits.
That's why minter.MintRemainder is always zero => we do not need to add it to the mintAmountDec

In the calculateMintedCoins function, the actual calculation is done by solving an integration in range [A; B].
We ensure that the max argument pass to the integral is no larger than FinalNormTimePassed.
This solves the problem with the precision loss that is aggregated in the accumulator (minter.NormTimePassed)


minter.NormTimePassed holds the current step as accumulator. Each block it is incremented by the step. Thus resulting in no loss in precision because minter.NormTimePassed = blockNumber * step, which is number that has no more than 18 decimal digits.
*/

var (
	// based on the assumption that we have 1 block per 5 seconds
	// if actual blocks are generated at slower rate then the network will mint tokens more than 3652 days (~10 years)
	denom                 = "acudos"         // Hardcoded to the acudos currency. Its not changeable, because some of the math depends on the size of this denomination
	totalDays             = sdk.NewInt(3652) // Hardcoded to 10 years
	InitialNormTimePassed = sdk.NewDecWithPrec(53172694105988, 14)
	FinalNormTimePassed   = sdk.NewDec(10)
	zeroPointSix          = sdk.MustNewDecFromStr("0.6")
	twentySixPointFive    = sdk.MustNewDecFromStr("26.5")
)

// Normalize block height incrementation
func normalizeBlockHeightInc(incrementModifier sdk.Int) sdk.Dec {
	totalBlocks := incrementModifier.Mul(totalDays)
	return (sdk.NewDec(1).QuoInt(totalBlocks)).Mul(FinalNormTimePassed)
}

// Integral of f(t) is 0,6 * t^3  - 26.5 * t^2 + 358 * t
// The function extrema is ~10.48 so after that the function is decreasing
func calculateIntegral(t sdk.Dec) sdk.Dec {
	return (zeroPointSix.Mul(t.Power(3))).Sub(twentySixPointFive.Mul(t.Power(2))).Add(sdk.NewDec(358).Mul(t))
}

// func calculateIntegralInNorm(t sdk.Dec) sdk.Dec {
// 	if t.LT(InitialNormTimePassed) {
// 		return sdk.NewDec(0)
// 	}

// 	if t.GT(FinalNormTimePassed) {
// 		return calculateIntegral(FinalNormTimePassed)
// 	}

// 	integralUpperbound := calculateIntegral(t)
// 	integralLowerbound := calculateIntegral(InitialNormTimePassed)
// 	return integralUpperbound.Sub(integralLowerbound)
// }

func calculateMintedCoins(minter types.Minter, increment sdk.Dec) sdk.Dec {
	prevStep := calculateIntegral(sdk.MinDec(minter.NormTimePassed, FinalNormTimePassed))
	nextStep := calculateIntegral(sdk.MinDec(minter.NormTimePassed.Add(increment), FinalNormTimePassed))
	return (nextStep.Sub(prevStep)).Mul(sdk.NewDec(10).Power(24)) // formula calculates in mil of cudos + converting to acudos
}

func logMintingInfo(ctx sdk.Context, k keeper.Keeper, minter types.Minter) {
	initiallySkipped := calculateIntegral(InitialNormTimePassed)
	mintedSoFar := calculateIntegral(sdk.MinDec(minter.NormTimePassed, FinalNormTimePassed))
	mintedSoFar = mintedSoFar.Sub(initiallySkipped).Mul(sdk.NewDec(10).Power(24))
	total := calculateIntegral(FinalNormTimePassed)
	total = total.Sub(initiallySkipped).Mul(sdk.NewDec(10).Power(24))
	k.Logger(ctx).Info("CudosMint module", "minted_so_far", mintedSoFar.TruncateInt().String()+denom, "left", total.Sub(mintedSoFar).TruncateInt().String()+denom, "total", total.TruncateInt().String()+denom)
}

// BeginBlocker mints new tokens for the previous block.
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	minter := k.GetMinter(ctx)
	params := k.GetParams(ctx)

	if minter.NormTimePassed.GT(FinalNormTimePassed) {
		return
	}

	incr := normalizeBlockHeightInc(params.IncrementModifier)
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

	logMintingInfo(ctx, k, minter)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeMintedDenom, denom),
			sdk.NewAttribute(types.AttributeMintedTokens, mintAmountInt.String()),
		),
	)
}
