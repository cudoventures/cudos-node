package v2

import (
	"fmt"
	"testing"

	"github.com/CudoVentures/cudos-node/x/nft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMigrateStore(t *testing.T) {
	encCfg := simappparams.MakeTestEncodingConfig()
	storeKey := sdk.NewKVStoreKey("nft")
	ctx := testutil.DefaultContext(storeKey, sdk.NewTransientStoreKey("transient_test"))

	denom := types.Denom{
		Id:     "testid",
		Name:   "testname",
		Symbol: "testsymbol",
	}

	// Setup store state as it was in the older version
	store := ctx.KVStore(storeKey)
	bz := encCfg.Marshaler.MustMarshal(&denom)
	keyByDenomID := types.KeyDenomID(denom.Id)
	store.Set(keyByDenomID, bz)
	store.Set(types.KeyDenomName(denom.Name), bz)
	store.Set(types.KeyDenomSymbol(denom.Symbol), bz)

	_, err := getDenomByName(ctx, storeKey, encCfg.Marshaler, denom.Name)
	require.Error(t, err)

	_, err = getDenomBySymbol(ctx, storeKey, encCfg.Marshaler, denom.Symbol)
	require.Error(t, err)

	require.NoError(t, MigrateStore(ctx, storeKey, encCfg.Marshaler))

	denomByName, err := getDenomByName(ctx, storeKey, encCfg.Marshaler, denom.Name)
	require.NoError(t, err)
	require.Equal(t, denom, denomByName)

	denomBySymbol, err := getDenomBySymbol(ctx, storeKey, encCfg.Marshaler, denom.Symbol)
	require.NoError(t, err)
	require.Equal(t, denom, denomBySymbol)
}

func getDenomByName(ctx sdk.Context, storeKey *sdk.KVStoreKey, cdc codec.BinaryCodec, name string) (denom types.Denom, err error) {
	store := ctx.KVStore(storeKey)

	keyDenomID := store.Get(types.KeyDenomName(name))
	if len(keyDenomID) == 0 {
		return denom, fmt.Errorf("not found denom name: %s", name)
	}

	bz := store.Get(keyDenomID)
	if len(bz) == 0 {
		return denom, fmt.Errorf("not found denom by denom id key: %s", string(keyDenomID))
	}

	cdc.MustUnmarshal(bz, &denom)
	return denom, nil
}

func getDenomBySymbol(ctx sdk.Context, storeKey *sdk.KVStoreKey, cdc codec.BinaryCodec, symbol string) (denom types.Denom, err error) {
	store := ctx.KVStore(storeKey)

	keyDenomID := store.Get(types.KeyDenomSymbol(symbol))
	if len(keyDenomID) == 0 {
		return denom, fmt.Errorf("not found denom symbol: %s", symbol)
	}

	bz := store.Get(keyDenomID)
	if len(bz) == 0 {
		return denom, fmt.Errorf("not found denom by denom id key: %s", string(keyDenomID))
	}

	cdc.MustUnmarshal(bz, &denom)
	return denom, nil
}
