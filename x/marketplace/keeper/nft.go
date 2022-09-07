package keeper

import (
	"encoding/binary"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetNftCount get the total number of nft
func (k Keeper) GetNftCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.NftCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// SetNftCount set the total number of nft
func (k Keeper) SetNftCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.NftCountKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendNft appends a nft in the store with a new id and update the count
func (k Keeper) AppendNft(
	ctx sdk.Context,
	nft types.Nft,
) uint64 {
	// Create the nft
	count := k.GetNftCount(ctx)

	// Set the ID of the appended value
	nft.Id = count

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NftKey))
	appendedValue := k.cdc.MustMarshal(&nft)
	store.Set(GetNftIDBytes(nft.Id), appendedValue)

	// Update nft count
	k.SetNftCount(ctx, count+1)

	return count
}

// SetNft set a specific nft in the store
func (k Keeper) SetNft(ctx sdk.Context, nft types.Nft) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NftKey))
	b := k.cdc.MustMarshal(&nft)
	store.Set(GetNftIDBytes(nft.Id), b)
}

// GetNft returns a nft from its id
func (k Keeper) GetNft(ctx sdk.Context, id uint64) (val types.Nft, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NftKey))
	b := store.Get(GetNftIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveNft removes a nft from the store
func (k Keeper) RemoveNft(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NftKey))
	store.Delete(GetNftIDBytes(id))

	count := k.GetNftCount(ctx)
	k.SetNftCount(ctx, count-1)
}

// GetAllNft returns all nft
func (k Keeper) GetAllNft(ctx sdk.Context) (list []types.Nft) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NftKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Nft
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetNftIDBytes returns the byte representation of the ID
func GetNftIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetNftIDFromBytes returns ID in uint64 format from a byte array
func GetNftIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
