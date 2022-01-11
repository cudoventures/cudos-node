package keeper

import (
	"github.com/CudoVentures/cudos-node/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) GetApprovedAddresses(ctx sdk.Context, sender sdk.AccAddress) (*types.ApprovedAddressesData, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyApprovedAddresses(sender.String()))
	if bz == nil {
		return &types.ApprovedAddressesData{}, sdkerrors.Wrapf(types.ErrNoApprovedAddresses, "no approved address for user: [%s]", sender.String())
	}
	var approvedAddress types.ApprovedAddresses
	k.cdc.MustUnmarshal(bz, &approvedAddress)

	return approvedAddress.ApprovedAddresses[sender.String()], nil

}

func (k Keeper) IsApprovedOperator(ctx sdk.Context, owner, operator sdk.AccAddress) bool {
	approvedAddressesData, err := k.GetApprovedAddresses(ctx, owner)
	if err != nil {
		return false
	}

	hasPermission := approvedAddressesData.ApprovedAddressesData[operator.String()]
	return hasPermission
}

func (k Keeper) SetApprovedAddress(ctx sdk.Context, sender, operator sdk.AccAddress, approved bool) {

	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyApprovedAddresses(sender.String())) // all types.ApprovedAddresses

	var approvedAddress types.ApprovedAddresses
	var approvedAddressesData types.ApprovedAddressesData

	if bz == nil {
		initApprovedAddresses(approvedAddress, approvedAddressesData, operator,
			sender, approved, k, store)
		return
	}

	k.cdc.MustUnmarshal(bz, &approvedAddress)

	if val, ok := approvedAddress.ApprovedAddresses[sender.String()]; ok {
		val.ApprovedAddressesData[operator.String()] = approved
	} else {
		approvedAddressesData.ApprovedAddressesData = map[string]bool{operator.String(): approved}
		approvedAddress.ApprovedAddresses[sender.String()] = &approvedAddressesData
	}

	bz = k.cdc.MustMarshal(&approvedAddress)
	store.Set(types.KeyApprovedAddresses(sender.String()), bz)
}

func initApprovedAddresses(
	approvedAddress types.ApprovedAddresses,
	approvedAddressesData types.ApprovedAddressesData,
	operator,
	sender sdk.AccAddress,
	approved bool,
	k Keeper,
	store sdk.KVStore) {
	// init inner map
	approvedAddressesData.ApprovedAddressesData = map[string]bool{operator.String(): approved}
	// init main map and set inner map as its value
	approvedAddress.ApprovedAddresses = map[string]*types.ApprovedAddressesData{sender.String(): &approvedAddressesData}
	// marshall and save
	bz := k.cdc.MustMarshal(&approvedAddress)
	store.Set(types.KeyApprovedAddresses(sender.String()), bz)
}
