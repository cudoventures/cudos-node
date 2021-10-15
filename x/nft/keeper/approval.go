package keeper

import (
	"cudos.org/cudos-node/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetApprovedAddress(ctx sdk.Context, sender, addressToBeApproved sdk.AccAddress, approved bool) {

	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyApprovedAddresses(sender.String())) // all types.ApprovedAddresses

	var approvedAddress types.ApprovedAddresses
	var approvedAddressesData types.ApprovedAddressesData

	if bz == nil {
		initApprovedAddresses(approvedAddress, approvedAddressesData, addressToBeApproved,
			sender, approved, k, store)
		return
	}

	k.cdc.MustUnmarshal(bz, &approvedAddress)
	if val, ok := approvedAddress.ApprovedAdresses[sender.String()]; ok {
		// if key is in map - update it
		val.ApprovedAddresses[addressToBeApproved.String()] = approved
	} else {
		// if not - init map and add it
		approvedAddressesData.ApprovedAddresses = map[string]bool{addressToBeApproved.String(): approved}
		approvedAddress.ApprovedAdresses[sender.String()] = &approvedAddressesData
	}

	bz = k.cdc.MustMarshal(&approvedAddress)
	store.Set(types.KeyApprovedAddresses(sender.String()), bz)
}

func initApprovedAddresses(
	approvedAddress types.ApprovedAddresses,
	approvedAddressesData types.ApprovedAddressesData,
	addressToBeApproved,
	sender sdk.AccAddress,
	approved bool,
	k Keeper,
	store sdk.KVStore) {
	// init inner map
	approvedAddressesData.ApprovedAddresses = map[string]bool{addressToBeApproved.String(): approved}
	// init main map and set inner map as its value
	approvedAddress.ApprovedAdresses = map[string]*types.ApprovedAddressesData{sender.String(): &approvedAddressesData}
	// marshall and save
	bz := k.cdc.MustMarshal(&approvedAddress)
	store.Set(types.KeyApprovedAddresses(sender.String()), bz)
}
