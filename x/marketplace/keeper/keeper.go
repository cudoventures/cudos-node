package keeper

import (
	"fmt"
	"strings"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	"github.com/CudoVentures/cudos-node/x/nft/exported"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   sdk.StoreKey
		memKey     sdk.StoreKey
		paramstore paramtypes.Subspace

		bankKeeper types.BankKeeper
		nftKeeper  types.NftKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,

	bankKeeper types.BankKeeper, nftKeeper types.NftKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{

		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,
		bankKeeper: bankKeeper, nftKeeper: nftKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) isCollectionPublished(ctx sdk.Context, denomID string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.KeyCollectionDenomID(denomID))
}

func (k Keeper) PublishCollection(ctx sdk.Context, collection types.Collection) (uint64, error) {
	denom, err := k.nftKeeper.GetDenom(ctx, collection.DenomId)
	if err != nil {
		return 0, err
	}

	if denom.Creator != collection.Owner {
		return 0, sdkerrors.Wrapf(types.ErrNotDenomOwner, "Owner of denom %s is %s", collection.DenomId, collection.Owner)
	}

	if k.isCollectionPublished(ctx, collection.DenomId) {
		return 0, sdkerrors.Wrapf(types.ErrCollectionAlreadyPublished, "Collection for denom %s is already published", collection.DenomId)
	}

	collectionID := k.AppendCollection(ctx, collection)

	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyCollectionDenomID(collection.DenomId), types.Uint64ToBytes(collectionID))

	return collectionID, nil
}

func (k Keeper) isApprovedNftAddress(nftVal exported.NFT, owner string) bool {
	approvedAddresses := nftVal.GetApprovedAddresses()
	for _, addr := range approvedAddresses {
		if addr == owner {
			return true
		}
	}
	return false
}

func (k Keeper) PublishNFT(ctx sdk.Context, nft types.Nft) (uint64, error) {
	if _, err := k.nftKeeper.GetDenom(ctx, nft.DenomId); err != nil {
		return 0, err
	}

	nftVal, err := k.nftKeeper.GetNFT(ctx, nft.DenomId, nft.TokenId)
	if err != nil {
		return 0, err
	}

	if nftVal.GetOwner().String() == nft.Owner ||
		k.nftKeeper.IsApprovedOperator(ctx, nftVal.GetOwner(), sdk.AccAddress(nft.Owner)) ||
		k.isApprovedNftAddress(nftVal, nft.Owner) {

		if err := k.nftKeeper.SoftLockNFT(ctx, types.ModuleName, nft.DenomId, nft.TokenId); err != nil {
			return 0, err
		}

		nftID := k.AppendNft(ctx, nft)

		store := ctx.KVStore(k.storeKey)
		store.Set(types.KeyNftDenomTokenID(nft.DenomId, nft.TokenId), types.Uint64ToBytes(nftID))

		return nftID, nil
	}

	return 0, nil
}

func (k Keeper) BuyNFT(ctx sdk.Context, nftID uint64, buyer sdk.AccAddress) error {
	nft, found := k.GetNft(ctx, nftID)
	if !found {
		return sdkerrors.Wrapf(types.ErrNftNotFound, "nft with id (%d) is not found for sale", nftID)
	}

	if nft.Owner == buyer.String() {
		return sdkerrors.Wrap(types.ErrCannotBuyOwnNft, "cannot buy own nft")
	}

	price, err := sdk.ParseCoinNormalized(nft.Price)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrInvalidPrice, "invalid price (%s)", nft.Price)
	}

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, buyer, types.ModuleName, sdk.NewCoins(price)); err != nil {
		return err
	}

	collection, found := k.getCollectionByDenomID(ctx, nft.DenomId)
	if !found || collection.ResaleRoyalties == "" {
		return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, buyer, sdk.NewCoins(price))
	}

	if err := k.distributeRoyalties(ctx, price, nft.Owner, collection.ResaleRoyalties); err != nil {
		return err
	}

	baseNft, err := k.nftKeeper.GetBaseNFT(ctx, nft.DenomId, nft.TokenId)
	if err != nil {
		return err
	}

	if err := k.nftKeeper.SoftUnlockNFT(ctx, types.ModuleName, nft.DenomId, nft.TokenId); err != nil {
		return err
	}

	k.nftKeeper.TransferNftInternal(ctx, nft.DenomId, nft.TokenId, sdk.AccAddress(nft.Owner), buyer, baseNft)

	return nil
}

func (k Keeper) MintNFT(ctx sdk.Context, denomID, priceStr, name, uri, data string, recipient sdk.AccAddress, sender sdk.AccAddress) (string, error) {
	denom, err := k.nftKeeper.GetDenom(ctx, denomID)
	if err != nil {
		return "", err
	}

	collection, found := k.getCollectionByDenomID(ctx, denomID)
	if !found {
		return "", sdkerrors.Wrapf(types.ErrCollectionNotFound, "collection %s not published for sale", denomID)
	}

	price, err := sdk.ParseCoinNormalized(priceStr)
	if err != nil {
		return "", sdkerrors.Wrapf(types.ErrInvalidPrice, "invalid price (%s)", priceStr)
	}

	if err := k.distributeRoyalties(ctx, price, denom.Creator, collection.MintRoyalties); err != nil {
		return "", err
	}

	return k.nftKeeper.MintNFT(ctx, denomID, name, uri, data, sender, recipient)
}

func (k Keeper) getCollectionByDenomID(ctx sdk.Context, denomID string) (types.Collection, bool) {
	store := ctx.KVStore(k.storeKey)
	collectionID := types.BytesToUint64(store.Get(types.KeyCollectionDenomID(denomID)))
	return k.GetCollection(ctx, collectionID)
}

func getProportion(totalCoin sdk.Coin, ratio sdk.Dec) sdk.Coin {
	return sdk.NewCoin(totalCoin.Denom, totalCoin.Amount.ToDec().Mul(ratio).TruncateInt())
}

func (k Keeper) distributeRoyalties(ctx sdk.Context, price sdk.Coin, seller, royalties string) error {
	if royalties == "" {
		return nil
	}

	var totalPercentPaid float64

	splitFn := func(c rune) bool {
		return c == ','
	}

	royaltiesList := strings.FieldsFunc(royalties, splitFn)

	for _, royalty := range royaltiesList {
		royaltyParts := strings.Split(royalty, ":")

		royaltyReceiver, err := sdk.AccAddressFromBech32(royaltyParts[0])
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid royalty address (%s): %s", royaltyParts[0], err)
		}

		royaltyPercent, err := sdk.NewDecFromStr(royaltyParts[1])
		if err != nil {
			return sdkerrors.Wrapf(types.ErrInvalidRoyaltyPercent, "invalid royalty percent (%s): %s", royaltyParts[1], err)
		}

		portion := getProportion(price, royaltyPercent)

		return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, royaltyReceiver, sdk.NewCoins(portion))
	}

	if totalPercentPaid < 100.0 {
		sellerAddr, err := sdk.AccAddressFromBech32(seller)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid seller address (%s): %s", seller, err)
		}

		sellerPercent := 100.0 - totalPercentPaid
		sellerLeftAmount, err := sdk.NewDecFromStr(fmt.Sprintf("%.2f", sellerPercent))
		if err != nil {
			return sdkerrors.Wrapf(types.ErrInvalidRoyaltyPercent, "invalid seller royalty percent (%f): %s", sellerPercent, err)
		}

		portion := getProportion(price, sellerLeftAmount)

		return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sellerAddr, sdk.NewCoins(portion))
	}

	return nil
}
