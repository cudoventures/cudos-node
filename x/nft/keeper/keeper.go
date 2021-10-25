package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"cudos.org/cudos-node/x/nft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	// this line is used by starport scaffolding # ibc/keeper/import
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type (
	Keeper struct {
		cdc      codec.Codec
		storeKey sdk.StoreKey
		memKey   sdk.StoreKey
		// this line is used by starport scaffolding # ibc/keeper/attribute

	}
)

// NewKeeper creates a new instance of the NFT Keeper
func NewKeeper(
	cdc codec.Codec,
	storeKey,
	memKey sdk.StoreKey,
	// this line is used by starport scaffolding # ibc/keeper/parameter

) *Keeper {
	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		memKey:   memKey,
		// this line is used by starport scaffolding # ibc/keeper/return

	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("cudos.org/cudos-node/%s", types.ModuleName))
}

// IssueDenom issues a denom according to the given params
func (k Keeper) IssueDenom(ctx sdk.Context, id, name, schema string, creator sdk.AccAddress) error {
	return k.SetDenom(ctx, types.NewDenom(id, name, schema, creator))
}

// MintNFTUnverified mints an NFT without verifying if the owner is the creator of denom
// Needed during genesis initialization
func (k Keeper) MintNFTUnverified(ctx sdk.Context, denomID, tokenID, tokenNm, tokenURI, tokenData string, owner sdk.AccAddress) error {
	if !k.HasDenomID(ctx, denomID) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denom ID %s not exists", denomID)
	}

	if k.HasNFT(ctx, denomID, tokenID) {
		return sdkerrors.Wrapf(types.ErrNFTAlreadyExists, "NFT %s already exists in collection %s", tokenID, denomID)
	}

	k.setNFT(
		ctx, denomID,
		types.NewBaseNFT(
			tokenID,
			tokenNm,
			owner,
			tokenURI,
			tokenData,
		),
	)
	k.setOwner(ctx, denomID, tokenID, owner)
	k.increaseSupply(ctx, denomID)

	return nil
}

// MintNFT mints an NFT and manages the NFT's existence within Collections and Owners
func (k Keeper) MintNFT(
	ctx sdk.Context, denomID, tokenID, tokenNm,
	tokenURI, tokenData string, sender, owner sdk.AccAddress,
) error {
	_, err := k.IsDenomCreator(ctx, denomID, sender)
	if err != nil {
		return err
	}

	return k.MintNFTUnverified(ctx, denomID, tokenID, tokenNm, tokenURI, tokenData, owner)
}

// EditNFT updates an already existing NFT
func (k Keeper) EditNFT(
	ctx sdk.Context, denomID, tokenID, tokenNm,
	tokenURI, tokenData string, sender sdk.AccAddress,
) error {
	if !k.HasDenomID(ctx, denomID) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denom ID %s not exists", denomID)
	}

	nft, err := k.GetBaseNFT(ctx, denomID, tokenID)
	if err != nil {
		return err
	}

	if !k.IsOwner(nft, sender) {
		return sdkerrors.Wrapf(types.ErrUnauthorized, "%s is not the owner of %s/%s", sender.String(), denomID, tokenID)
	}

	if types.Modified(tokenNm) {
		nft.Name = tokenNm
	}

	if types.Modified(tokenURI) {
		nft.URI = tokenURI
	}

	if types.Modified(tokenData) {
		nft.Data = tokenData
	}

	k.setNFT(ctx, denomID, nft)

	return nil
}

// TransferOwner transfers the ownership of the given NFT to the new owner
func (k Keeper) TransferOwner(ctx sdk.Context, denomID, tokenID string, from, to, sender sdk.AccAddress) error {
	if !k.HasDenomID(ctx, denomID) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denom ID %s not exists", denomID)
	}

	nft, err := k.GetBaseNFT(ctx, denomID, tokenID)
	if err != nil {
		return err
	}

	if !from.Equals(nft.GetOwner()) {
		return sdkerrors.Wrapf(types.ErrUnauthorized,
			"From [%s] is not the owner of NFT with denomId [%s] / tokenId [%s]. The owner is [%s]", sender.String(), denomID, tokenID, nft.GetOwner())
	}

	if sender.Equals(nft.GetOwner()) || // if the owner is requesting the transfer
		(nft.ApprovedAddresses != nil && nft.ApprovedAddresses[sender.String()]) || // or if the sender is approved for the nft
		k.IsApprovedOperator(ctx, from, sender) { // or if the sender is part of approveAll of user
		transferNFT(ctx, denomID, tokenID, from, to, nft, k)
		return nil
	}

	return sdkerrors.Wrapf(types.ErrUnauthorized,
		"Sender [%s] is neither owner or approved for transfer of denomId [%s] / tokenId [%s]", sender.String(), denomID, tokenID)
}

func transferNFT(ctx sdk.Context, denomID string, tokenID string, from sdk.AccAddress, to sdk.AccAddress, nft types.BaseNFT, k Keeper) {
	nft.ApprovedAddresses = nil
	nft.Owner = to.String()
	k.setNFT(ctx, denomID, nft)
	k.swapOwner(ctx, denomID, tokenID, from, to)
}

// BurnNFT deletes a specified NFT
func (k Keeper) BurnNFT(ctx sdk.Context, denomID, tokenID string, owner sdk.AccAddress) error {
	if !k.HasDenomID(ctx, denomID) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denom ID %s not exists", denomID)
	}

	nft, err := k.GetBaseNFT(ctx, denomID, tokenID)
	if err != nil {
		return err
	}

	if !k.IsOwner(nft, owner) {
		return sdkerrors.Wrapf(types.ErrUnauthorized, "%s is not the owner of %s/%s", owner.String(), denomID, tokenID)
	}

	k.deleteNFT(ctx, denomID, nft)
	k.deleteOwner(ctx, denomID, tokenID, owner)
	k.decreaseSupply(ctx, denomID)

	return nil
}

func (k Keeper) AddApproval(ctx sdk.Context, denomID, tokenID string, sender sdk.AccAddress, approvedAddress sdk.AccAddress) error {

	nft, err := k.GetBaseNFT(ctx, denomID, tokenID)
	if err != nil {
		return err
	}

	if nft.GetOwner().Equals(sender) || k.IsApprovedOperator(ctx, nft.GetOwner(), sender) {
		k.ApproveNFT(ctx, nft, approvedAddress, denomID)
		return nil
	}

	return sdkerrors.Wrapf(types.ErrUnauthorized,
		"Approve failed - could not authorize (%s)! Sender address (%s) is neither owner or approved for denomId (%s) / tokenId (%s)! ", approvedAddress, sender, denomID, tokenID)
}

// Todo: check if we need this to be private. For example, right now its defined in the keeper
// if it is accessible from there - it means all the check are bypassed ??
func approveNFT(nft types.BaseNFT, approvedAddress sdk.AccAddress) {
	if nft.ApprovedAddresses == nil {
		nft.ApprovedAddresses = map[string]bool{approvedAddress.String(): true}
	} else {
		nft.ApprovedAddresses[approvedAddress.String()] = true
	}
}

func (k Keeper) AddApprovalForAll(ctx sdk.Context, sender sdk.AccAddress, operatorAddressToBeAdded sdk.AccAddress, approved bool) error {
	if sender.Equals(operatorAddressToBeAdded) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "ApproveAll failed! Sender address (%s) is the same as operator (%s)! ", sender, operatorAddressToBeAdded)
	}
	k.SetApprovedAddress(ctx, sender, operatorAddressToBeAdded, approved)
	return nil
}

func (k Keeper) RevokeApproval(ctx sdk.Context, denomID, tokenID string, sender, addressToRevoke sdk.AccAddress) error {
	nft, err := k.GetBaseNFT(ctx, denomID, tokenID)
	if err != nil {
		return err
	}

	if nft.GetOwner().Equals(sender) || k.IsApprovedOperator(ctx, nft.GetOwner(), sender) {
		err := k.RevokeApprovalNFT(ctx, nft, addressToRevoke, denomID)
		if err != nil {
			return err
		}
		return nil
	}

	return sdkerrors.Wrapf(types.ErrUnauthorized,
		"Approve failed - could not revoke access for (%s)! Sender address (%s) is neither owner or approved for denomId (%s) / tokenId (%s)! ", addressToRevoke, sender, denomID, tokenID)
}

// Todo: check if we need this to be private. For example, right now its defined in the keeper
// if it is accessible from there - it means all the check are bypassed ??
func revokeApprovalNFT(nft types.BaseNFT, approvedAddress sdk.AccAddress, denomID string) error {
	if nft.ApprovedAddresses == nil {
		return sdkerrors.Wrapf(types.ErrNoApprovedAddresses, "No approved address (%s) for nft with denomId (%s) / tokenId (%s)", approvedAddress.String(), denomID, nft.GetID())
	}

	_, ok := nft.ApprovedAddresses[approvedAddress.String()]
	if ok {
		delete(nft.ApprovedAddresses, approvedAddress.String())
	} else {
		return sdkerrors.Wrapf(types.ErrNoApprovedAddresses, "No approved address (%s) for nft with denomId (%s) / tokenId (%s)", approvedAddress.String(), denomID, nft.GetID())
	}
	return nil
}
