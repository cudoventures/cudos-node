package keeper

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"cudos.org/cudos-node/x/nft/types"
)

func createNNFT(keeper *Keeper, ctx sdk.Context, n int) []types.NFT {
	items := make([]types.NFT, n)
	for i := range items {
		items[i].Creator = "any"
		items[i].Index = fmt.Sprintf("%d", i)
		keeper.SetNFT(ctx, items[i])
	}
	return items
}

func TestNFTGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNNFT(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetNFT(ctx, item.Index)
		assert.True(t, found)
		assert.Equal(t, item, rst)
	}
}
func TestNFTRemove(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNNFT(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveNFT(ctx, item.Index)
		_, found := keeper.GetNFT(ctx, item.Index)
		assert.False(t, found)
	}
}

func TestNFTGetAll(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNNFT(keeper, ctx, 10)
	assert.Equal(t, items, keeper.GetAllNFT(ctx))
}
