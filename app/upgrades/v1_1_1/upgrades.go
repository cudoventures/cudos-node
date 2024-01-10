package v1_1_1

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// CreateUpgradeHandler creates an SDK upgrade handler for v1
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	gk wasm.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", UpgradeName)
		UpdateWasmParams(ctx, gk)
		logger.Info("running module migrations ...")
		return mm.RunMigrations(ctx, configurator, vm)
	}
}

func UpdateWasmParams(ctx sdk.Context, gk wasm.Keeper) {
	params := gk.GetParams(ctx)
	params.CodeUploadAccess.Permission = wasmtypes.AccessTypeNobody
	gk.SetParams(ctx, params)
}
