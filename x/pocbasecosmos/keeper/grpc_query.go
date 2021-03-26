package keeper

import (
	"github.com/rdpnd/poc-base-cosmos/x/pocbasecosmos/types"
)

var _ types.QueryServer = Keeper{}
