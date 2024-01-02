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
var burnDenom = "acudos"

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
	burnAmts := sdk.Coins{}
	for _, amt := range amounts {
		if k.dkSet && amt.Denom == burnDenom {
			if err := k.dk.FundCommunityPool(ctx, sdk.NewCoins(amt), k.ak.GetModuleAddress(moduleName)); err != nil {
				return err
			}
		} else {
			burnAmts = burnAmts.Add(amt)
		}
	}

	if len(burnAmts) > 0 {
		return k.BaseKeeper.BurnCoins(ctx, moduleName, burnAmts)
	}
	return nil
	// acc := k.ak.GetModuleAccount(ctx, moduleName)
	// if acc == nil {
	// 	panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", moduleName))
	// }

	// if !acc.HasPermission(authtypes.Burner) {
	// 	panic(sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "module account %s does not have permissions to burn tokens", moduleName))
	// }

	// err := k.subUnlockedCoins(ctx, acc.GetAddress(), amounts)
	// if err != nil {
	// 	return err
	// }

	// for _, amount := range amounts {
	// 	if amount.Denom == burnDenom {
	// 		// transfer collected fees to the distribution module account
	// 		err := k.addCoins(ctx, k.dk.GetDistributionAccount(ctx).GetAddress(), sdk.Coins{amount})
	// 		if err != nil {
	// 			panic(err)
	// 		}

	// 		fp := k.dk.GetFeePool(ctx)
	// 		fp.CommunityPool = fp.CommunityPool.Add(sdk.NewDecCoinFromCoin(amount))
	// 		k.dk.SetFeePool(ctx, fp)
	// 		// resultString = "moved tokens from module account to community pool"
	// 	} else {
	// 		supply := k.GetSupply(ctx, amount.GetDenom())
	// 		supply = supply.Sub(amount)
	// 		k.setSupply(ctx, supply)
	// 		// resultString = "burned tokens from module account"
	// 	}
	// }

	// // emit burn event
	// ctx.EventManager().EmitEvent(
	// 	types.NewCoinBurnEvent(acc.GetAddress(), amounts),
	// )

	// return nil
}
