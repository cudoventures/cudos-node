package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"

	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	authzcodec "github.com/cosmos/cosmos-sdk/x/authz/codec"
	govcodec "github.com/cosmos/cosmos-sdk/x/gov/codec"
	groupcodec "github.com/cosmos/cosmos-sdk/x/group/codec"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
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
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)

func init() {
	RegisterLegacyAminoCodec(Amino)
	cryptocodec.RegisterCrypto(Amino)
	sdk.RegisterLegacyAminoCodec(Amino)

	// Register all Amino interfaces and concrete types on the authz and gov Amino codec so that this can later be
	// used to properly serialize MsgGrant, MsgExec and MsgSubmitProposal instances
	RegisterLegacyAminoCodec(authzcodec.Amino)
	RegisterLegacyAminoCodec(govcodec.Amino)
	RegisterLegacyAminoCodec(groupcodec.Amino)
}
