package keeper

import (
	"github.com/CudoVentures/cudos-node/x/token/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetToken set a specific token in the store from its index
func (k Keeper) SetToken(ctx sdk.Context, token types.Token) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TokenKeyPrefix))
	b := k.cdc.MustMarshal(&token)
	store.Set(types.TokenKey(
		token.Denom,
	), b)
}

// GetToken returns a token from its index
func (k Keeper) GetToken(
	ctx sdk.Context,
	denom string,

) (val types.Token, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TokenKeyPrefix))

	b := store.Get(types.TokenKey(
		denom,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveToken removes a token from the store
func (k Keeper) RemoveToken(
	ctx sdk.Context,
	denom string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TokenKeyPrefix))
	store.Delete(types.TokenKey(
		denom,
	))
}

// GetAllToken returns all token
func (k Keeper) GetAllToken(ctx sdk.Context) (list []types.Token) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TokenKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Token
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
