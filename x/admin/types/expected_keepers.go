package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type DistributionKeeper interface {
	DistributeFromFeePool(ctx sdk.Context, amount sdk.Coins, receiveAddr sdk.AccAddress) error
}

type BankKeeper interface {
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
}
