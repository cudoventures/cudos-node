package app

import (
	"encoding/json"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmKeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	marketplaceCustomBindings "github.com/CudoVentures/cudos-node/x/marketplace/custom-bindings"
	marketplacekeeper "github.com/CudoVentures/cudos-node/x/marketplace/keeper"
	nftCustomBindings "github.com/CudoVentures/cudos-node/x/nft/custom-bindings"
	nftkeeper "github.com/CudoVentures/cudos-node/x/nft/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetCustomPlugins(nftKeeper nftkeeper.Keeper, marketplaceKeeper marketplacekeeper.Keeper) []wasmKeeper.Option {
	nftEncodingOpt := wasmKeeper.WithMessageEncoders(&wasmKeeper.MessageEncoders{
		Custom: nftCustomBindings.EncodeNftMessage(),
	})

	queryPluginsOpt := wasmKeeper.WithQueryPlugins(&wasmKeeper.QueryPlugins{
		Custom: func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
			bz, err := nftCustomBindings.PerformCustomNftQuery(nftKeeper)(ctx, request)
			if err == nil {
				return bz, nil
			}

			return marketplaceCustomBindings.PerformCustomMarketplaceQuery(marketplaceKeeper)(ctx, request)
		},
	})

	marketplaceMessengerDecoratorOpt := wasmKeeper.WithMessageHandlerDecorator(
		marketplaceCustomBindings.NewMarketplaceWasmMessageDecorator(marketplaceKeeper),
	)

	return []wasm.Option{nftEncodingOpt, queryPluginsOpt, marketplaceMessengerDecoratorOpt}
}
