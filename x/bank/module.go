package bank

import (
	"fmt"

	custombankkeeper "github.com/CudoVentures/cudos-node/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankmodule "github.com/cosmos/cosmos-sdk/x/bank"
	bankexported "github.com/cosmos/cosmos-sdk/x/bank/exported"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

type AppModule struct {
	bankmodule.AppModule

	keeper custombankkeeper.Keeper

	legacySubspace bankexported.Subspace
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper custombankkeeper.Keeper, accountKeeper types.AccountKeeper, ss bankexported.Subspace) AppModule {
	bankModule := bankmodule.NewAppModule(cdc, keeper, accountKeeper, ss)
	return AppModule{
		AppModule:      bankModule,
		keeper:         keeper,
		legacySubspace: ss,
	}
}

// RegisterServices registers module services.
// NOTE: Overriding this method as not doing so will cause a panic
// when trying to force this custom keeper into a bankkeeper.BaseKeeper
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), bankkeeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), am.keeper)
	m := bankkeeper.NewMigrator(am.keeper.BaseKeeper, am.legacySubspace)
	if err := cfg.RegisterMigration(types.ModuleName, 1, m.Migrate1to2); err != nil {
		panic(fmt.Sprintf("failed to migrate x/bank from version 1 to 2: %v", err))
	}
}
