package keeper

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"

	"cudos.org/cudos-node/x/nft/exported"
	"cudos.org/cudos-node/x/nft/types"
)

// SetGenesisCollection saves all NFTs and returns an error if there already exists or any one of the owner's bech32
// account address is invalid
func (k Keeper) SetGenesisCollection(ctx sdk.Context, collection types.Collection) error {
	for _, nft := range collection.NFTs {
		if _, err := k.MintNFTUnverified(
			ctx,
			collection.Denom.Id,
			nft.GetName(),
			nft.GetURI(),
			nft.GetData(),
			nft.GetOwner(),
		); err != nil {
			return err
		}
	}
	return nil
}

// SetCollection saves all NFTs and returns an error if there already exists or any one of the owner's bech32 account
// address is invalid or any NFT's owner is not the creator of denomination
func (k Keeper) SetCollection(ctx sdk.Context, collection types.Collection, sender sdk.AccAddress) error {
	for _, nft := range collection.NFTs {
		if _, err := k.MintNFT(
			ctx,
			collection.Denom.Id,
			nft.GetName(),
			nft.GetURI(),
			nft.GetData(),
			sender,
			nft.GetOwner(),
		); err != nil {
			return err
		}
	}
	return nil
}

// GetCollection returns the collection by the specified denom ID
func (k Keeper) GetCollection(ctx sdk.Context, denomID string) (types.Collection, error) {
	denom, err := k.GetDenom(ctx, denomID)
	if err != nil {
		return types.Collection{}, sdkerrors.Wrapf(types.ErrInvalidDenom, "denomID %s not existed ", denomID)
	}

	nfts := k.GetNFTs(ctx, denomID)
	return types.NewCollection(denom, nfts), nil
}

// GetPaginateCollection returns the collection by the specified denom ID
func (k Keeper) GetPaginateCollection(ctx sdk.Context, request *types.QueryCollectionRequest, denomID string) (types.Collection, *query.PageResponse, error) {
	denom, err := k.GetDenom(ctx, denomID)
	if err != nil {
		return types.Collection{}, nil, sdkerrors.Wrapf(types.ErrInvalidDenom, "denomID %s not existed ", denomID)
	}
	var nfts []exported.NFT
	store := ctx.KVStore(k.storeKey)
	nftStore := prefix.NewStore(store, types.KeyNFT(denomID, ""))
	pageRes, err := query.Paginate(nftStore, request.Pagination, func(key []byte, value []byte) error {
		var baseNFT types.BaseNFT
		k.cdc.MustUnmarshal(value, &baseNFT)
		nfts = append(nfts, baseNFT)
		return nil
	})
	if err != nil {
		return types.Collection{}, nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}
	return types.NewCollection(denom, nfts), pageRes, nil
}

// GetCollections returns all the collections
func (k Keeper) GetCollections(ctx sdk.Context) (cs []types.Collection) {
	for _, denom := range k.GetDenoms(ctx) {
		nfts := k.GetNFTs(ctx, denom.Id)
		cs = append(cs, types.NewCollection(denom, nfts))
	}
	return cs
}

// GetTotalSupply returns the number of NFTs by the specified denom ID
func (k Keeper) GetTotalSupply(ctx sdk.Context, denomID string) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyCollection(denomID))
	if len(bz) == 0 {
		return 0
	}
	return types.MustUnMarshalSupply(k.cdc, bz)
}

// GetTotalSupplyOfOwner returns the amount of NFTs by the specified conditions
func (k Keeper) GetTotalSupplyOfOwner(ctx sdk.Context, id string, owner sdk.AccAddress) (supply uint64) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyOwner(owner, id, ""))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		supply++
	}
	return supply
}

func (k Keeper) increaseSupply(ctx sdk.Context, denomID string) {
	supply := k.GetTotalSupply(ctx, denomID)
	supply++

	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalSupply(k.cdc, supply)
	store.Set(types.KeyCollection(denomID), bz)
}

func (k Keeper) decreaseSupply(ctx sdk.Context, denomID string) {
	supply := k.GetTotalSupply(ctx, denomID)
	supply--

	store := ctx.KVStore(k.storeKey)
	if supply == 0 {
		store.Delete(types.KeyCollection(denomID))
		return
	}

	bz := types.MustMarshalSupply(k.cdc, supply)
	store.Set(types.KeyCollection(denomID), bz)
}

// GetNftTotalCountForCollection returns the count of all minted nfts ( including burned )
func (k Keeper) GetNftTotalCountForCollection(ctx sdk.Context, denomID string) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyCollectionTotalNfts(denomID))
	if len(bz) == 0 {
		return 0
	}
	return types.MustUnMarshalTotalNftCountForCollection(k.cdc, bz)
}

// IncrementTotalCounterForCollection increments the count of all minted nfts for a collection
func (k Keeper) IncrementTotalCounterForCollection(ctx sdk.Context, denomID string) {
	totalCount := k.GetNftTotalCountForCollection(ctx, denomID)
	totalCount++

	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshallTotalCountForCollection(k.cdc, totalCount)
	store.Set(types.KeyCollectionTotalNfts(denomID), bz)
}
