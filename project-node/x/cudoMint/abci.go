package cudoMint

import (
	"cudos.org/cudos-node/x/cudoMint/keeper"
	"cudos.org/cudos-node/x/cudoMint/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

// Minting based on the polynomial "=(28653300+5000*n)/(1.5+0.004*n^2)"
// where n indicates number of months since Ethereum minting (n = 1 means Jan 2021, n = 2 means Feb 2021, etc.)
var (
	// based on the assumption that we have 1 block per 5 seconds
	blocksPerMonth sdk.Int = sdk.NewInt(525657)
	// regulate offset of n
	monthsOffset int64  = 0
	monthsActive int64  = 9
	denom        string = "cudos"
)

func calculateMintedCoins(monthsPassed sdk.Int, minter types.Minter) sdk.Dec {
	monthsPassed = monthsPassed.AddRaw(1) // the algorithm is 1-based
	monthsDenominator := sdk.MustNewDecFromStr("1.5").Add(sdk.MustNewDecFromStr("0.004")).Mul(monthsPassed.ToDec().Power(2))
	coinsForMonth := sdk.NewDec(28653300 + monthsPassed.MulRaw(5000).Int64()).Quo(monthsDenominator)
	return (coinsForMonth.QuoInt(blocksPerMonth)).Add(minter.MintRemainder)
}

// BeginBlocker mints new tokens for the previous block.
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	blockHeight := ctx.BlockHeight()
	monthsPassed := sdk.NewInt(blockHeight).Quo(blocksPerMonth).AddRaw(monthsOffset)
	if monthsPassed.GT(sdk.NewInt(monthsActive)) {
		return
	}

	minter := k.GetMinter(ctx)
	mintAmountDec := calculateMintedCoins(monthsPassed, minter)
	mintAmountInt := mintAmountDec.TruncateInt()
	mintedCoin := sdk.NewCoin(denom, mintAmountInt)
	mintedCoins := sdk.NewCoins(mintedCoin)
	err := k.MintCoins(ctx, mintedCoins)
	if err != nil {
		panic(err)
	}
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
			sdk.NewAttribute(types.AttributeMintedDenom, "stake"),
			sdk.NewAttribute(types.AttributeMintedTokens, mintAmountInt.String()),
		),
	)
}
