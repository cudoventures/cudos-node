package cudoMint

import (
	"cudos.org/cudos-node/x/cudoMint/keeper"
	"cudos.org/cudos-node/x/cudoMint/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

// Minting based on pre defined yearly values
var (
	// based on the assumption that we have 1 block per 5 seconds
	blocksPerMonth sdk.Int = sdk.NewInt(525657)
	denom          string  = "acudos"
	tokensPerYear          = map[int]sdk.Dec{
		2021: sdk.NewDec(306_000_000).Mul(sdk.NewDec(10).Power(18)),
		2022: sdk.NewDec(272_000_000).Mul(sdk.NewDec(10).Power(18)),
		2023: sdk.NewDec(238_000_000).Mul(sdk.NewDec(10).Power(18)),
		2024: sdk.NewDec(204_000_000).Mul(sdk.NewDec(10).Power(18)),
		2025: sdk.NewDec(170_000_000).Mul(sdk.NewDec(10).Power(18)),
		2026: sdk.NewDec(68_000_000).Mul(sdk.NewDec(10).Power(18)),
		2027: sdk.NewDec(68_000_000).Mul(sdk.NewDec(10).Power(18)),
		2028: sdk.NewDec(68_000_000).Mul(sdk.NewDec(10).Power(18)),
		2029: sdk.NewDec(68_000_000).Mul(sdk.NewDec(10).Power(18)),
		2030: sdk.NewDec(68_000_000).Mul(sdk.NewDec(10).Power(18)),
	}
)

func calculateMintedCoins(year int, minter types.Minter) sdk.Dec {
	if yearlyTokens, ok := tokensPerYear[year]; ok {
		return yearlyTokens.QuoInt(blocksPerMonth).Add(minter.MintRemainder)
	} else {
		return sdk.ZeroDec()
	}
}

// BeginBlocker mints new tokens for the previous block.
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	minter := k.GetMinter(ctx)
	mintAmountDec := calculateMintedCoins(ctx.BlockTime().Year(), minter)
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
			sdk.NewAttribute(types.AttributeMintedDenom, denom),
			sdk.NewAttribute(types.AttributeMintedTokens, mintAmountInt.String()),
		),
	)
}
