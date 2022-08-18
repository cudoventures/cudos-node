package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
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

	store := ctx.KVStore(k.storeKey)

	collectionID := k.AppendCollection(ctx, collection)
	store.Set(types.KeyCollectionDenomID(collection.DenomId), types.Uint64ToBytes(collectionID))

	return collectionID, nil
}
