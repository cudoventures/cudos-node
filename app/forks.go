package app

import (
	"github.com/CudoVentures/cudos-node/app/upgrades/v1_1_1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlockForks executes any necessary fork logic based upon the current block height.
func BeginBlockForks(ctx sdk.Context, app *App) {
	v1_1_1.UpdateWasmParams(ctx, app.wasmKeeper)
}
