package keeper_test

import (
	"strconv"
	"testing"

	keepertest "github.com/CudoVentures/cudos-node/testutil/keeper"
	"github.com/CudoVentures/cudos-node/testutil/nullify"
	"github.com/CudoVentures/cudos-node/x/token/keeper"
	"github.com/CudoVentures/cudos-node/x/token/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// func NewTestTokenKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
// 	storeKey := sdk.NewKVStoreKey(types.StoreKey)
// 	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

// 	db := tmdb.NewMemDB()
// 	stateStore := store.NewCommitMultiStore(db)
// 	stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
// 	stateStore.MountStoreWithDB(memStoreKey, sdk.StoreTypeMemory, nil)
// 	require.NoError(t, stateStore.LoadLatestVersion())

// 	registry := codectypes.NewInterfaceRegistry()
// 	cdc := codec.NewProtoCodec(registry)

// 	paramsSubspace := typesparams.NewSubspace(cdc,
// 		types.Amino,
// 		storeKey,
// 		memStoreKey,
// 		"TokenParams",
// 	)
// 	k := keeper.NewKeeper(
// 		cdc,
// 		storeKey,
// 		memStoreKey,
// 		paramsSubspace,
// 		nil,
// 	)

// 	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

// 	return k, ctx
// }

// Prevent strconv unused error
var _ = strconv.IntSize

func createNToken(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Token {
	items := make([]types.Token, n)
	for i := range items {
		items[i].Denom = strconv.Itoa(i)

		keeper.SaveToken(ctx, items[i])
	}
	return items
}

func TestGetTokenByDenom(t *testing.T) {
	keeper, ctx := keepertest.TestTokenKeeper(t)
	items := createNToken(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetTokenByDenom(ctx,
			item.Denom,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}

func TestGetAllTokens(t *testing.T) {
	keeper, ctx := keepertest.TestTokenKeeper(t)
	items := createNToken(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllTokens(ctx)),
	)
}
