package keeper

import (
	"encoding/binary"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetCollectionCount get the total number of collection
func (k Keeper) GetCollectionCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.CollectionCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// SetCollectionCount set the total number of collection
func (k Keeper) SetCollectionCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.CollectionCountKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendCollection appends a collection in the store with a new id and update the count
func (k Keeper) AppendCollection(
	ctx sdk.Context,
	collection types.Collection,
) uint64 {
	// Create the collection
	count := k.GetCollectionCount(ctx)

	// Set the ID of the appended value
	collection.Id = count

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CollectionKey))
	appendedValue := k.cdc.MustMarshal(&collection)
	store.Set(GetCollectionIDBytes(collection.Id), appendedValue)

	// Update collection count
	k.SetCollectionCount(ctx, count+1)

	return count
}

// SetCollection set a specific collection in the store
func (k Keeper) SetCollection(ctx sdk.Context, collection types.Collection) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CollectionKey))
	b := k.cdc.MustMarshal(&collection)
	store.Set(GetCollectionIDBytes(collection.Id), b)
}

// GetCollection returns a collection from its id
func (k Keeper) GetCollection(ctx sdk.Context, id uint64) (val types.Collection, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CollectionKey))
	b := store.Get(GetCollectionIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveCollection removes a collection from the store
func (k Keeper) RemoveCollection(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CollectionKey))
	store.Delete(GetCollectionIDBytes(id))
}

// GetAllCollection returns all collection
func (k Keeper) GetAllCollection(ctx sdk.Context) (list []types.Collection) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CollectionKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Collection
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetCollectionIDBytes returns the byte representation of the ID
func GetCollectionIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetCollectionIDFromBytes returns ID in uint64 format from a byte array
func GetCollectionIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
