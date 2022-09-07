package keeper

import (
	"github.com/CudoVentures/cudos-node/x/marketplace/types"
)

var _ types.QueryServer = Keeper{}
