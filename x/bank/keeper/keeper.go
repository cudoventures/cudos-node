package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	accountkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Copied from https://github.com/CudoVentures/cosmos-sdk/blob/3816012a2d4ea5c9bbb3d8e6174d3b96ff91a039/x/bank/keeper/keeper.go#L19C1-L20C1
var sendToCommunityDenom = "acudos"
var ignoredDenom = "cudosAdmin"

type Keeper struct {
	bankkeeper.BaseKeeper

	ak accountkeeper.AccountKeeper

	dk    distrkeeper.Keeper
	dkSet bool
}

var _ bankkeeper.Keeper = Keeper{}

func NewCustomKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ak accountkeeper.AccountKeeper,
	paramSpace paramtypes.Subspace,
	blockedAddrs map[string]bool,
) Keeper {
	keeper := Keeper{
		BaseKeeper: bankkeeper.NewBaseKeeper(cdc, storeKey, ak, paramSpace, blockedAddrs),
		ak:         ak,
	}
	return keeper
}

func (k *Keeper) SetDistrKeeper(dk distrkeeper.Keeper) {
	k.dk = dk
	k.dkSet = true
}

// Migrate from https://github.com/CudoVentures/cosmos-sdk/blob/3816012a2d4ea5c9bbb3d8e6174d3b96ff91a039/x/bank/keeper/keeper.go#L439
func (k Keeper) BurnCoins(ctx sdk.Context, moduleName string, amounts sdk.Coins) error {
	if !k.dkSet {
		panic("distr keeper not set for bank keeper")
	}
	burnAmts := sdk.Coins{}
	for _, amt := range amounts {
		// Send to community pool if denom is acudos, if cudosAdmin ignore, else burn
		if amt.Denom == sendToCommunityDenom {
			if err := k.dk.FundCommunityPool(ctx, sdk.NewCoins(amt), k.ak.GetModuleAddress(moduleName)); err != nil {
				return err
			}

		} else if amt.Denom == ignoredDenom {
			continue
		} else {
			burnAmts = burnAmts.Add(amt)
		}
	}

	if len(burnAmts) > 0 {
		return k.BaseKeeper.BurnCoins(ctx, moduleName, burnAmts)
	}
	return nil
}
