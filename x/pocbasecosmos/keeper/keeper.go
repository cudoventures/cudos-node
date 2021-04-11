package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/rdpnd/poc-base-cosmos/x/pocbasecosmos/types"
)

type (
	Keeper struct {
		cdc                codec.Marshaler
		storeKey           sdk.StoreKey
		memKey             sdk.StoreKey
		distributionKeeper types.DistributionKeeper
		bankKeeper         types.BankKeeper
	}
)

func NewKeeper(cdc codec.Marshaler, storeKey, memKey sdk.StoreKey,
	dk types.DistributionKeeper, bk types.BankKeeper) *Keeper {
	return &Keeper{
		cdc:                cdc,
		storeKey:           storeKey,
		memKey:             memKey,
		distributionKeeper: dk,
		bankKeeper:         bk,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) AdminDistributeFromFeePool(ctx sdk.Context, amount sdk.Coins, receiveAddr sdk.AccAddress) error {
	return k.distributionKeeper.DistributeFromFeePool(ctx, amount, receiveAddr)
}
