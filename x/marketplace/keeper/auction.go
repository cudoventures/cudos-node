package keeper

import (
	"encoding/binary"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetAuctionCount get the total number of auction
func (k Keeper) GetAuctionCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.AuctionCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// SetAuctionCount set the total number of auction
func (k Keeper) SetAuctionCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.AuctionCountKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendAuction appends a auction in the store with a new id and update the count
func (k Keeper) AppendAuction(ctx sdk.Context, auction types.Auction) (uint64, error) {
	auction.Id = k.GetAuctionCount(ctx)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AuctionKey))
	b, err := k.cdc.Marshal(&auction)
	if err != nil {
		return 0, err
	}

	store.Set(GetAuctionIDBytes(auction.Id), b)
	k.SetAuctionCount(ctx, auction.Id+1)
	return auction.Id, nil
}

// SetAuction set a specific auction in the store
func (k Keeper) SetAuction(ctx sdk.Context, auction types.Auction) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AuctionKey))
	b, err := k.cdc.Marshal(&auction)
	if err != nil {
		return err
	}

	store.Set(GetAuctionIDBytes(auction.Id), b)

	return nil
}

// GetAuction returns a auction from its id
func (k Keeper) GetAuction(ctx sdk.Context, id uint64) (val types.Auction, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AuctionKey))
	bz := store.Get(GetAuctionIDBytes(id))
	if bz == nil {
		return val, false
	}

	if err := k.cdc.Unmarshal(bz, &val); err != nil {
		return val, false
	}

	return val, true
}

// RemoveAuction removes a auction from the store
func (k Keeper) RemoveAuction(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AuctionKey))
	store.Delete(GetAuctionIDBytes(id))
}

// GetAllAuction returns all auction
func (k Keeper) GetAllAuction(ctx sdk.Context) (list []types.Auction) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AuctionKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Auction
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetAuctionIDBytes returns the byte representation of the ID
func GetAuctionIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}
