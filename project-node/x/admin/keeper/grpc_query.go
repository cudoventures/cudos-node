package keeper

import (
	"cudos.org/cudos-node/x/admin/types"
)

var _ types.QueryServer = Keeper{}
