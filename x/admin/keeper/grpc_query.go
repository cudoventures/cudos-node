package keeper

import (
	"github.com/CudoVentures/cudos-node/x/admin/types"
)

var _ types.QueryServer = Keeper{}
