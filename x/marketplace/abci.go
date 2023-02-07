package marketplace

import (
	"github.com/CudoVentures/cudos-node/x/marketplace/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abcitypes "github.com/tendermint/tendermint/abci/types"
)

func EndBlock(ctx sdk.Context, k keeper.Keeper) []abcitypes.ValidatorUpdate {
	if err := k.AuctionEndBlocker(ctx); err != nil {
		panic(err)
	}

	return []abcitypes.ValidatorUpdate{}
}
