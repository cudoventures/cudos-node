package keeper

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/CudoVentures/cudos-node/x/nft/types"
)

// HasDenomID returns whether the specified denom ID exists
func (k Keeper) HasDenomID(ctx sdk.Context, id string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.KeyDenomID(id))
}

// HasDenomNm returns whether the specified denom name exists
func (k Keeper) HasDenomNm(ctx sdk.Context, name string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.KeyDenomName(name))
}

// HasDenomSymbol returns whether the specified denom symbol exists
func (k Keeper) HasDenomSymbol(ctx sdk.Context, name string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.KeyDenomSymbol(name))
}

// SetDenom is responsible for saving the definition of denom
func (k Keeper) SetDenom(ctx sdk.Context, denom types.Denom) error {
	if k.HasDenomID(ctx, denom.Id) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denomID %s has already exists", denom.Id)
	}

	if k.HasDenomNm(ctx, denom.Name) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denomName %s has already exists", denom.Name)
	}

	if k.HasDenomSymbol(ctx, denom.Symbol) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denomSymbol %s has already exists", denom.Symbol)
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&denom)
	keyByDenomID := types.KeyDenomID(denom.Id)
	store.Set(keyByDenomID, bz)
	store.Set(types.KeyDenomName(denom.Name), keyByDenomID)
	store.Set(types.KeyDenomSymbol(denom.Symbol), keyByDenomID)
	return nil
}

// GetDenom returns the denom by id
func (k Keeper) GetDenom(ctx sdk.Context, id string) (denom types.Denom, err error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyDenomID(id))
	if len(bz) == 0 {
		return denom, sdkerrors.Wrapf(types.ErrInvalidDenom, "not found denomID: %s", id)
	}

	k.cdc.MustUnmarshal(bz, &denom)
	return denom, nil
}

// GetDenom returns the denom by name
func (k Keeper) GetDenomByName(ctx sdk.Context, name string) (denom types.Denom, err error) {
	store := ctx.KVStore(k.storeKey)

	keyDenomID := store.Get(types.KeyDenomName(name))
	if len(keyDenomID) == 0 {
		return denom, sdkerrors.Wrapf(types.ErrInvalidDenom, "not found denom name: %s", name)
	}

	bz := store.Get(keyDenomID)
	if len(bz) == 0 {
		return denom, sdkerrors.Wrapf(types.ErrInvalidDenom, "not found denom by denom id key: %s", string(keyDenomID))
	}

	k.cdc.MustUnmarshal(bz, &denom)
	return denom, nil
}

// GetDenomBySymbol returns the denom by symbol
func (k Keeper) GetDenomBySymbol(ctx sdk.Context, symbol string) (denom types.Denom, err error) {
	store := ctx.KVStore(k.storeKey)

	keyDenomID := store.Get(types.KeyDenomSymbol(symbol))
	if len(keyDenomID) == 0 {
		return denom, sdkerrors.Wrapf(types.ErrInvalidDenom, "not found denom symbol: %s", symbol)
	}

	bz := store.Get(keyDenomID)
	if len(bz) == 0 {
		return denom, sdkerrors.Wrapf(types.ErrInvalidDenom, "not found denom by denom id key: %s", string(keyDenomID))
	}

	k.cdc.MustUnmarshal(bz, &denom)
	return denom, nil
}

// GetDenoms returns all the denoms
func (k Keeper) GetDenoms(ctx sdk.Context) (denoms []types.Denom) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyDenomID(""))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var denom types.Denom
		k.cdc.MustUnmarshal(iterator.Value(), &denom)
		denoms = append(denoms, denom)
	}
	return denoms
}

// IsDenomCreator checks if address is the creator of Denom
// Return the Denom if true, an error otherwise
func (k Keeper) IsDenomCreator(ctx sdk.Context, denomID string, address sdk.AccAddress) (types.Denom, error) {
	denom, err := k.GetDenom(ctx, denomID)
	if err != nil {
		return types.Denom{}, err
	}

	creator, err := sdk.AccAddressFromBech32(denom.Creator)
	if err != nil {
		return types.Denom{}, err
	}

	if !creator.Equals(address) {
		return denom, sdkerrors.Wrapf(types.ErrUnauthorized, "%s is not the creator of %s", address, denomID)
	}

	return denom, nil
}

func (k Keeper) IsDenomMinter(denom types.Denom, address sdk.AccAddress) error {
	if denom.Minter == "" {
		return nil
	}

	minter, err := sdk.AccAddressFromBech32(denom.Minter)
	if err != nil {
		return err
	}

	if !minter.Equals(address) {
		return sdkerrors.Wrapf(types.ErrUnauthorized, "%s is not the minter of %s", address, denom.Id)
	}

	return nil
}

func (k Keeper) IsEditable(ctx sdk.Context, denomID string) error {
	denom, err := k.GetDenom(ctx, denomID)
	if err != nil {
		return err
	}

	traits := strings.Split(denom.Traits, ",")
	for _, trait := range traits {
		if types.DenomTraitsMapStrToType[trait] == types.NotEditable {
			return sdkerrors.Wrapf(types.ErrNotEditable, "denom '%s' has not editable trait", denomID)
		}
	}

	return nil
}

// UpdateDenom is responsible for updating the definition of denom
func (k Keeper) UpdateDenom(ctx sdk.Context, denom types.Denom) error {
	if !k.HasDenomID(ctx, denom.Id) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denomID %s not exists", denom.Id)
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&denom)
	store.Set(types.KeyDenomID(denom.Id), bz)
	return nil
}
