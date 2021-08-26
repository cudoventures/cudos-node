package nft

import (
	"cudos.org/cudos-node/x/nft/keeper"
	"cudos.org/cudos-node/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	// Set all the Nft
	for _, elem := range genState.NFTList {
		k.SetNFT(ctx, *elem)
	}

	// this line is used by starport scaffolding # ibc/genesis/init
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// this line is used by starport scaffolding # genesis/module/export
	// Get all Nft
	nFTList := k.GetAllNFT(ctx)
	for _, elem := range nFTList {
		elem := elem
		genesis.NFTList = append(genesis.NFTList, &elem)
	}

	// this line is used by starport scaffolding # ibc/genesis/export

	return genesis
}
