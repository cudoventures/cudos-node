package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"cudos.org/cudos-node/x/nft/exported"
	"cudos.org/cudos-node/x/nft/types"
)

// GetNFT set a specific NFT in the store from its index
func (k Keeper) GetNFT(ctx sdk.Context, denomID, tokenID string) (nft exported.NFT, err error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyNFT(denomID, tokenID))
	if bz == nil {
		return nil, sdkerrors.Wrapf(types.ErrUnknownCollection, "not found NFT: %s", denomID)
	}

	var baseNFT types.BaseNFT
	k.cdc.MustUnmarshal(bz, &baseNFT)

	return baseNFT, nil
}

// GetNFT set a specific NFT in the store from its index
func (k Keeper) GetBaseNFT(ctx sdk.Context, denomID, tokenID string) (nft types.BaseNFT, err error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyNFT(denomID, tokenID))
	if bz == nil {
		return types.BaseNFT{}, sdkerrors.Wrapf(types.ErrNotFoundNFT, "not found NFT: %s", denomID)
	}

	var baseNFT types.BaseNFT
	k.cdc.MustUnmarshal(bz, &baseNFT)

	return baseNFT, nil
}

// GetNFTs returns all NFTs by the specified denom ID
func (k Keeper) GetNFTs(ctx sdk.Context, denom string) (nfts []exported.NFT) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.KeyNFT(denom, ""))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var baseNFT types.BaseNFT
		k.cdc.MustUnmarshal(iterator.Value(), &baseNFT)
		nfts = append(nfts, baseNFT)
	}

	return nfts
}

// IsOwner checks if the sender is the owner of the given NFT
// Return the NFT if true, an error otherwise
func (k Keeper) IsOwner(nft types.BaseNFT, owner sdk.AccAddress) bool {
	if !owner.Equals(nft.GetOwner()) {
		return false
	}
	return true
}

// HasNFT checks if the specified NFT exists
func (k Keeper) HasNFT(ctx sdk.Context, denomID, tokenID string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.KeyNFT(denomID, tokenID))
}

func (k Keeper) setNFT(ctx sdk.Context, denomID string, nft types.BaseNFT) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&nft)
	store.Set(types.KeyNFT(denomID, nft.GetID()), bz)
}

func (k Keeper) ApproveNFT(ctx sdk.Context, nft types.BaseNFT, approvedAddress sdk.AccAddress, denomID string) {
	if nft.ApprovedAddresses == nil {
		nft.ApprovedAddresses = map[string]bool{approvedAddress.String(): true}
	} else {
		nft.ApprovedAddresses[approvedAddress.String()] = true
	}
	k.setNFT(ctx, denomID, nft)
}

func (k Keeper) RevokeApprovalNFT(ctx sdk.Context, nft types.BaseNFT, addressToRevoke sdk.AccAddress, denomID string) error {

	if nft.ApprovedAddresses == nil {
		return sdkerrors.Wrapf(types.ErrNoApprovedAddresses, "No approved address (%s) for nft with denomId (%s) / tokenId (%s)", addressToRevoke.String(), denomID, nft.GetID())
	}

	_, ok := nft.ApprovedAddresses[addressToRevoke.String()]
	if ok {
		delete(nft.ApprovedAddresses, addressToRevoke.String())
	} else {
		return sdkerrors.Wrapf(types.ErrNoApprovedAddresses, "No approved address (%s) for nft with denomId (%s) / tokenId (%s)", addressToRevoke.String(), denomID, nft.GetID())
	}

	k.setNFT(ctx, denomID, nft)
	return nil
}

// deleteNFT deletes an existing NFT from store
func (k Keeper) deleteNFT(ctx sdk.Context, denomID string, nft exported.NFT) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyNFT(denomID, nft.GetID()))
}
