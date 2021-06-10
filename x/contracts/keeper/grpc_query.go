package keeper

import (
	"cudos.org/cudos-node/x/contracts/types"
)

var _ types.QueryServer = Keeper{}
