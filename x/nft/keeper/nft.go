package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/CudoVentures/cudos-node/x/nft/exported"
	"github.com/CudoVentures/cudos-node/x/nft/types"
)

func (Keeper) IsApprovedAddress(nft *types.BaseNFT, sender string) bool {
	for _, address := range nft.ApprovedAddresses {
		if sender == address {
			return true
		}
	}
	return false
}

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
		return types.BaseNFT{}, sdkerrors.Wrapf(types.ErrNotFoundNFT, "not found NFT: denomId: %s, tokenId: %s", denomID, tokenID)
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
func (Keeper) IsOwner(nft types.BaseNFT, owner sdk.AccAddress) bool {
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
		nft.ApprovedAddresses = []string{approvedAddress.String()}
	} else {
		nft.ApprovedAddresses = append(nft.ApprovedAddresses, approvedAddress.String())
	}
	k.setNFT(ctx, denomID, nft)
}

func (k Keeper) RevokeApprovalNFT(ctx sdk.Context, nft types.BaseNFT, addressToRevoke sdk.AccAddress, denomID string) error {
	if nft.ApprovedAddresses == nil {
		return sdkerrors.Wrapf(types.ErrNoApprovedAddresses, "No approved address (%s) for nft with denomId (%s) / tokenId (%s)", addressToRevoke.String(), denomID, nft.GetID())
	}

	// Searching for the given address and removing it if found by reslicing the array (shifts all elements at the right of the deleted index by one to the left )
	for i, address := range nft.ApprovedAddresses {
		if address == addressToRevoke.String() {
			nft.ApprovedAddresses = append(nft.ApprovedAddresses[:i], nft.ApprovedAddresses[i+1:]...)
			k.setNFT(ctx, denomID, nft)
			return nil
		}
	}
	return sdkerrors.Wrapf(types.ErrNoApprovedAddresses, "No approved address (%s) for nft with denomId (%s) / tokenId (%s)", addressToRevoke.String(), denomID, nft.GetID())
}

// deleteNFT deletes an existing NFT from store
func (k Keeper) deleteNFT(ctx sdk.Context, denomID string, nft exported.NFT) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyNFT(denomID, nft.GetID()))
}

// GetNFTApprovedAddresses returns the approved addresses for the nft
func (k Keeper) GetNFTApprovedAddresses(ctx sdk.Context, denomID, tokenID string) (approvedAddresses []string, err error) {
	nft, err := k.GetBaseNFT(ctx, denomID, tokenID)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrUnknownNFT, "invalid NFT %s from collection %s", tokenID, denomID)
	}

	if nft.ApprovedAddresses == nil {
		return nil,
			sdkerrors.Wrapf(types.ErrNoApprovedAddresses, "No approved addresses for NFT %s from collection %s", tokenID, denomID)
	}

	return nft.ApprovedAddresses, nil
}

func (k Keeper) getSoftLockOwner(ctx sdk.Context, denomID, tokenID string) string {
	store := ctx.KVStore(k.storeKey)
	bLockOwnerID := store.Get(types.KeyNFTLockOwner(denomID, tokenID))
	if bLockOwnerID == nil {
		return ""
	}
	return string(bLockOwnerID)
}

func (k Keeper) SoftLockNFT(ctx sdk.Context, lockOwner, denomID, tokenID string) error {
	currentLockOwner := k.getSoftLockOwner(ctx, denomID, tokenID)
	if currentLockOwner != "" {
		return sdkerrors.Wrapf(types.ErrAlreadySoftLocked, "Failed to acquire soft lock on Denom %s NFT %s for %s because already acquired by %s",
			denomID, tokenID, lockOwner, currentLockOwner)
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyNFTLockOwner(denomID, tokenID), []byte(lockOwner))

	return nil
}

func (k Keeper) SoftUnlockNFT(ctx sdk.Context, lockOwner, denomID, tokenID string) error {
	currentLockOwner := k.getSoftLockOwner(ctx, denomID, tokenID)
	if currentLockOwner == "" {
		return sdkerrors.Wrapf(types.ErrNotSoftLocked, "Failed to release soft lock because Denom %s NFT %s is not locked", denomID, tokenID)
	}

	if currentLockOwner != lockOwner {
		return sdkerrors.Wrapf(types.ErrNotOwnerOfSoftLock, "Failed to release soft lock on Denom %s NFT %s for %s because its acquired by %s",
			denomID, tokenID, lockOwner, currentLockOwner)
	}

	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyNFTLockOwner(denomID, tokenID))

	return nil
}

func (k Keeper) IsSoftLocked(ctx sdk.Context, denomID, tokenID string) error {
	if currentLockOwner := k.getSoftLockOwner(ctx, denomID, tokenID); currentLockOwner != "" {
		return sdkerrors.Wrapf(types.ErrSoftLocked, "token id %s from denom with id %s is soft locked by %s", tokenID, denomID, currentLockOwner)
	}
	return nil
}
