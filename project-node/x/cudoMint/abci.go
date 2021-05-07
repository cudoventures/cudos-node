package cudoMint

import (
	"cudos.org/cudos-node/x/cudoMint/keeper"
	"cudos.org/cudos-node/x/cudoMint/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

// Minting based on arithmetic progression
var (
	blockSpan int64 = 1_000_000
	mintTokens int64 = 10_000_000_000

	//Find d where: an = a1 + (n - 1)d, n is blockSpan and the sum of the series is mintTokens
	commonDifference sdk.Dec = sdk.NewDec(2 * mintTokens).Quo(sdk.NewDec(blockSpan * (blockSpan - 1)))
)


// BeginBlocker mints new tokens for the previous block.
func 	BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	logger := k.Logger(ctx)
	logger.Info("Test me?")
	blockHeight := ctx.BlockHeight()
	if blockHeight > blockSpan {
		return
	}
	logger.Info("Im here aint?")
	mintAmount := sdk.NewDec(blockSpan - blockHeight).Mul(commonDifference)
	mintedCoin := sdk.NewCoin("stake", mintAmount.TruncateInt())
	mintedCoins := sdk.NewCoins(mintedCoin)
	err := k.MintCoins(ctx, mintedCoins)
	if err != nil {
		panic(err)
	}

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
			sdk.NewAttribute(types.AttributeMintedTokens, mintAmount.String()),
		),
	)
}
