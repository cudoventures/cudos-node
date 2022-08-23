package keeper

import (
	"github.com/CudoVentures/cudos-node/x/addressbook/types"
)

var _ types.QueryServer = Keeper{}
