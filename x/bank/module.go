package bank

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankmodule "github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/bank/types"

	custombankkeeper "github.com/CudoVentures/cudos-node/x/bank/keeper"
)

type AppModule struct {
	bankmodule.AppModule

	keeper custombankkeeper.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper custombankkeeper.Keeper, accountKeeper types.AccountKeeper) AppModule {
	bankModule := bankmodule.NewAppModule(cdc, keeper, accountKeeper)
	return AppModule{
		AppModule: bankModule,
		keeper:    keeper,
	}
}

// RegisterServices registers module services.
// NOTE: Overriding this method as not doing so will cause a panic
// when trying to force this custom keeper into a bankkeeper.BaseKeeper
func (am AppModule) RegisterServices(cfg module.Configurator) {
	am.AppModule.RegisterServices(cfg)

}
