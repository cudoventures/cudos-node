package app

import (
	"github.com/CudoVentures/cudos-node/app/upgrades/v1_1_1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// BeginBlockForks executes any necessary fork logic based upon the current block height.
func BeginBlockForks(ctx sdk.Context, app *App) {
	upgradePlan := upgradetypes.Plan{
		Height: ctx.BlockHeight(),
	}

	switch ctx.BlockHeight() {
	case v1_1_1.UpgradeHeight:
		upgradePlan.Name = v1_1_1.UpgradeName
		upgradePlan.Info = v1_1_1.UpgradeInfo
	default:
		// do nothing
		return
	}

	err := app.UpgradeKeeper.ScheduleUpgrade(ctx, upgradePlan)
	if err != nil {
		panic(err)
	}
}
