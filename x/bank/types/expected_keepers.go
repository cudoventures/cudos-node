package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BankHooks interface {
	AfterBurnCoinHook(ctx sdk.Context, moduleName sdk.AccAddress, amount sdk.Coins) error
}
