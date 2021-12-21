package app

import (
	nftCustomBindings "cudos.org/cudos-node/x/nft/custom-bindings"
	nftKeeper "cudos.org/cudos-node/x/nft/keeper"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmKeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
)

func GetCustomMsgEncodersOptions() []wasmKeeper.Option {
	nftEncodingOptions := wasmKeeper.WithMessageEncoders(nftEncoders())
	return []wasm.Option{nftEncodingOptions}
}

func GetCustomMsgQueryOptions(keeper nftKeeper.Keeper) []wasmKeeper.Option {
	nftQueryOptions := wasmKeeper.WithQueryPlugins(nftQueryPlugins(keeper))
	return []wasm.Option{nftQueryOptions}
}

func nftEncoders() *wasmKeeper.MessageEncoders {
	return &wasmKeeper.MessageEncoders{
		Custom: nftCustomBindings.EncodeNftMessage(),
	}
}

// nftQueryPlugins needs to be registered in test setup to handle custom query callbacks
func nftQueryPlugins(keeper nftKeeper.Keeper) *wasmKeeper.QueryPlugins {
	return &wasmKeeper.QueryPlugins{
		Custom: nftCustomBindings.PerformCustomNftQuery(keeper),
	}
}
