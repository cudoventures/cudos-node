package keeper

import (
	"github.com/CudoVentures/cudos-node/x/cudomint/types"
)

var _ types.QueryServer = Keeper{}
