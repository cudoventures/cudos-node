package simulation

import (
	"cosmossdk.io/math"
	cudominttypes "github.com/CudoVentures/cudos-node/x/cudoMint/types"
	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func fundAcc(ctx sdk.Context, bk types.BankKeeper, recipient sdk.AccAddress, amount sdk.Coin) {
	// minting 4 times more than requested to ensure that there is enough funds for tx fees
	amount = sdk.NewCoin(amount.Denom, amount.Amount.Mul(math.NewInt(4)))
	bk.MintCoins(ctx, cudominttypes.ModuleName, sdk.NewCoins(amount))
	bk.SendCoinsFromModuleToAccount(ctx, cudominttypes.ModuleName, recipient, sdk.NewCoins(amount))
}
