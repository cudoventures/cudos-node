package keeper

import (
	"cudos.org/cudos-node/x/nft/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetNFT set a specific NFT in the store from its index
func (k Keeper) SetNFT(ctx sdk.Context, nft types.NFT) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NFTKey))
	b := k.cdc.MustMarshal(&nft)
	store.Set(types.KeyPrefix(nft.Index), b)
}

// GetNFT returns a Nft from its index
func (k Keeper) GetNFT(ctx sdk.Context, index string) (val types.NFT, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NFTKey))

	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}

	k.cdc.Unmarshal(b, &val)
	return val, true
}

// RemoveNFT removes a Nft from the store
func (k Keeper) RemoveNFT(ctx sdk.Context, index string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NFTKey))
	store.Delete(types.KeyPrefix(index))
}

// GetAllNFT returns all Nft
func (k Keeper) GetAllNFT(ctx sdk.Context) (list []types.NFT) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NFTKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.NFT
		k.cdc.Unmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
