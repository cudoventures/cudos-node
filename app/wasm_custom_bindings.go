package app

import (
	"encoding/json"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmKeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	addressbookCustomBindings "github.com/CudoVentures/cudos-node/x/addressbook/custom-bindings"
	addressbookkeeper "github.com/CudoVentures/cudos-node/x/addressbook/keeper"
	marketplaceCustomBindings "github.com/CudoVentures/cudos-node/x/marketplace/custom-bindings"
	marketplacekeeper "github.com/CudoVentures/cudos-node/x/marketplace/keeper"
	nftCustomBindings "github.com/CudoVentures/cudos-node/x/nft/custom-bindings"
	nftkeeper "github.com/CudoVentures/cudos-node/x/nft/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func GetCustomPlugins(
	nftKeeper nftkeeper.Keeper,
	marketplaceKeeper marketplacekeeper.Keeper,
	addressbookKeeper addressbookkeeper.Keeper,
) []wasmKeeper.Option {

	queryHandlers := []func(ctx sdk.Context, request json.RawMessage) ([]byte, error){
		nftCustomBindings.PerformCustomNftQuery(nftKeeper),
		addressbookCustomBindings.PerformCustomAddressbookQuery(addressbookKeeper),
		marketplaceCustomBindings.PerformCustomMarketplaceQuery(marketplaceKeeper),
	}

	queryPluginsOpt := wasmKeeper.WithQueryPlugins(&wasmKeeper.QueryPlugins{
		Custom: func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
			for _, handler := range queryHandlers {
				if bz, err := handler(ctx, request); err == nil {
					return bz, nil
				}
			}
			return nil, sdkerrors.Wrap(wasmtypes.ErrInvalidMsg, "No custom query handler was able to process the request")
		},
	})

	nftEncodingOpt := wasmKeeper.WithMessageEncoders(&wasmKeeper.MessageEncoders{
		Custom: nftCustomBindings.EncodeNftMessage(),
	})

	marketplaceMessengerDecoratorOpt := wasmKeeper.WithMessageHandlerDecorator(
		marketplaceCustomBindings.NewMarketplaceWasmMessageDecorator(marketplaceKeeper),
	)

	addressbookMessengerDecoratorOpt := wasmKeeper.WithMessageHandlerDecorator(
		addressbookCustomBindings.NewAddressbookWasmMessageDecorator(addressbookKeeper),
	)

	return []wasm.Option{
		nftEncodingOpt,
		queryPluginsOpt,
		marketplaceMessengerDecoratorOpt,
		addressbookMessengerDecoratorOpt,
	}
}
