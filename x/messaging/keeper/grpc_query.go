package keeper

import (
	"github.com/CudoVentures/cudos-node/x/messaging/types"
)

var _ types.QueryServer = Keeper{}
