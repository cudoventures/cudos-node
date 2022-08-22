package keeper

import (
	"github.com/CudoVentures/cudos-node/x/token/types"
)

var _ types.QueryServer = Keeper{}
