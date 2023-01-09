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
	if bz == nil {
		return 0
	}

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
func (k Keeper) AppendAuction(ctx sdk.Context, a types.Auction) (uint64, error) {
	a.SetId(k.GetAuctionCount(ctx))
	bz, err := k.cdc.MarshalInterface(a)
	if err != nil {
		return 0, err
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AuctionKey))
	store.Set(GetAuctionIDBytes(a.GetId()), bz)
	k.SetAuctionCount(ctx, a.GetId()+1)
	return a.GetId(), nil
}

// SetAuction set a specific auction in the store
func (k Keeper) SetAuction(ctx sdk.Context, a types.Auction) error {
	bz, err := k.cdc.MarshalInterface(a)
	if err != nil {
		return err
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AuctionKey))
	store.Set(GetAuctionIDBytes(a.GetId()), bz)
	return nil
}

// GetAuction returns a auction from its id
func (k Keeper) GetAuction(ctx sdk.Context, id uint64) (types.Auction, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AuctionKey))
	bz := store.Get(GetAuctionIDBytes(id))
	if bz == nil {
		return nil, types.ErrAuctionNotFound
	}

	var a types.Auction
	err := k.cdc.UnmarshalInterface(bz, &a)
	return a, err
}

// RemoveAuction removes a auction from the store
func (k Keeper) RemoveAuction(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AuctionKey))
	store.Delete(GetAuctionIDBytes(id))
}

// GetAllAuction returns all auction
func (k Keeper) GetAllAuction(ctx sdk.Context) ([]types.Auction, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AuctionKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	var auctions []types.Auction
	for ; iterator.Valid(); iterator.Next() {
		var a types.Auction
		if err := k.cdc.UnmarshalInterface(iterator.Value(), &a); err != nil {
			return auctions, err
		}
		auctions = append(auctions, a)
	}

	return auctions, nil
}

// GetAuctionIDBytes returns the byte representation of the ID
func GetAuctionIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}
