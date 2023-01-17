package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgPublishCollection{}, "marketplace/PublishCollection", nil)
	cdc.RegisterConcrete(&MsgPublishNft{}, "marketplace/PublishNft", nil)
	cdc.RegisterConcrete(&MsgBuyNft{}, "marketplace/BuyNft", nil)
	cdc.RegisterConcrete(&MsgMintNft{}, "marketplace/MintNft", nil)
	cdc.RegisterConcrete(&MsgRemoveNft{}, "marketplace/RemoveNft", nil)
	cdc.RegisterConcrete(&MsgVerifyCollection{}, "marketplace/VerifyCollection", nil)
	cdc.RegisterConcrete(&MsgUnverifyCollection{}, "marketplace/UnverifyCollection", nil)
	cdc.RegisterConcrete(&MsgCreateCollection{}, "marketplace/CreateCollection", nil)
	cdc.RegisterConcrete(&MsgUpdateRoyalties{}, "marketplace/UpdateRoyalties", nil)
	cdc.RegisterConcrete(&MsgUpdatePrice{}, "marketplace/UpdatePrice", nil)
	cdc.RegisterConcrete(&MsgAddAdmin{}, "marketplace/AddAdmin", nil)
	cdc.RegisterConcrete(&MsgRemoveAdmin{}, "marketplace/RemoveAdmin", nil)
	cdc.RegisterConcrete(&MsgPublishAuction{}, "marketplace/PublishAuction", nil)
	cdc.RegisterConcrete(&MsgPlaceBid{}, "marketplace/Bid", nil)
	cdc.RegisterInterface((*Auction)(nil), nil)
	cdc.RegisterConcrete(&EnglishAuction{}, "marketplace/EnglishAuction", nil)
	cdc.RegisterConcrete(&DutchAuction{}, "marketplace/DutchAuction", nil)
	cdc.RegisterConcrete(&MsgAcceptBid{}, "marketplace/AcceptBid", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgPublishCollection{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgPublishNft{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgBuyNft{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgMintNft{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRemoveNft{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgVerifyCollection{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUnverifyCollection{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateCollection{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateRoyalties{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdatePrice{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddAdmin{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRemoveAdmin{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgPublishAuction{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgPlaceBid{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAcceptBid{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)

	registry.RegisterInterface(
		"cudoventures.cudosnode.marketplace.Auction",
		(*Auction)(nil),
		&EnglishAuction{},
		&DutchAuction{},
	)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
