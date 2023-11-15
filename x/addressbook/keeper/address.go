package keeper

import (
	"github.com/CudoVentures/cudos-node/x/addressbook/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// CreateNewAddress set a specific address in the store from its index
func (k Keeper) CreateNewAddress(ctx sdk.Context, address types.Address) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AddressKeyPrefix))

	key := types.AddressKey(address.Creator, address.Network, address.Label)

	if b := store.Get(key); b != nil {
		return sdkerrors.Wrapf(types.ErrKeyAlreadyExists, "Key (%s) already exists", key)
	}

	b := k.cdc.MustMarshal(&address)
	store.Set(types.AddressKey(address.Creator, address.Network, address.Label), b)

	return nil
}

// UpdateExistingAddress set a specific address in the store from its index
func (k Keeper) UpdateExistingAddress(ctx sdk.Context, address types.Address) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AddressKeyPrefix))

	key := types.AddressKey(address.Creator, address.Network, address.Label)

	if b := store.Get(key); b == nil {
		return sdkerrors.Wrapf(types.ErrKeyNotFound, "Key (%s) not found", key)
	}

	b := k.cdc.MustMarshal(&address)
	store.Set(types.AddressKey(address.Creator, address.Network, address.Label), b)

	return nil
}

// GetAddress returns a address from its index
func (k Keeper) GetAddress(ctx sdk.Context, creator, network, label string) (val types.Address, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AddressKeyPrefix))

	b := store.Get(types.AddressKey(creator, network, label))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveAddress removes a address from the store
func (k Keeper) RemoveAddress(ctx sdk.Context, creator, network, label string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AddressKeyPrefix))
	store.Delete(types.AddressKey(creator, network, label))
}

// GetAllAddress returns all address
func (k Keeper) GetAllAddress(ctx sdk.Context) (list []types.Address) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AddressKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Address
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
