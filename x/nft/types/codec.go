package types

// DONTCOVER

import (
	gogotypes "github.com/cosmos/gogoproto/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"

	"github.com/CudoVentures/cudos-node/x/nft/exported"

	authzcodec "github.com/cosmos/cosmos-sdk/x/authz/codec"
	govcodec "github.com/cosmos/cosmos-sdk/x/gov/codec"
	groupcodec "github.com/cosmos/cosmos-sdk/x/group/codec"
)

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()

	// Register all Amino interfaces and concrete types on the authz and gov Amino codec so that this can later be
	// used to properly serialize MsgGrant, MsgExec and MsgSubmitProposal instances
	RegisterLegacyAminoCodec(authzcodec.Amino)
	RegisterLegacyAminoCodec(govcodec.Amino)
	RegisterLegacyAminoCodec(groupcodec.Amino)
}

// RegisterLegacyAminoCodec concrete types on codec
// (Amino is still needed for Ledger at the moment)
// nolint: staticcheck
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgIssueDenom{}, "github.com/CudoVentures/cudos-node/nft/MsgIssueDenom", nil)
	cdc.RegisterConcrete(&MsgTransferNft{}, "github.com/CudoVentures/cudos-node/nft/MsgTransferNft", nil)
	cdc.RegisterConcrete(&MsgApproveNft{}, "github.com/CudoVentures/cudos-node/nft/MsgApproveNft", nil)
	cdc.RegisterConcrete(&MsgApproveAllNft{}, "github.com/CudoVentures/cudos-node/nft/MsgApproveAllNft", nil)
	cdc.RegisterConcrete(&MsgRevokeNft{}, "github.com/CudoVentures/cudos-node/nft/MsgRevokeNft", nil)
	cdc.RegisterConcrete(&MsgEditNFT{}, "github.com/CudoVentures/cudos-node/nft/MsgEditNFT", nil)
	cdc.RegisterConcrete(&MsgMintNFT{}, "github.com/CudoVentures/cudos-node/nft/MsgMintNFT", nil)
	cdc.RegisterConcrete(&MsgBurnNFT{}, "github.com/CudoVentures/cudos-node/nft/MsgBurnNFT", nil)
	cdc.RegisterConcrete(&MsgTransferDenom{}, "github.com/CudoVentures/cudos-node/nft/MsgTransferDenom", nil)

	cdc.RegisterInterface((*exported.NFT)(nil), nil)
	cdc.RegisterConcrete(&BaseNFT{}, "github.com/CudoVentures/cudos-node/nft/BaseNFT", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgIssueDenom{},
		&MsgTransferNft{},
		&MsgApproveNft{},
		&MsgRevokeNft{},
		&MsgEditNFT{},
		&MsgMintNFT{},
		&MsgBurnNFT{},
		&MsgTransferDenom{},
		&MsgApproveAllNft{},
	)

	registry.RegisterImplementations((*exported.NFT)(nil),
		&BaseNFT{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// return supply protobuf code
func MustMarshalSupply(cdc codec.Codec, supply uint64) []byte {
	supplyWrap := gogotypes.UInt64Value{Value: supply}
	return cdc.MustMarshal(&supplyWrap)
}

// return th supply
func MustUnMarshalSupply(cdc codec.Codec, value []byte) uint64 {
	var supplyWrap gogotypes.UInt64Value
	cdc.MustUnmarshal(value, &supplyWrap)
	return supplyWrap.Value
}

// return the tokenID protobuf code
func MustMarshalTokenID(cdc codec.Codec, tokenID string) []byte {
	tokenIDWrap := gogotypes.StringValue{Value: tokenID}
	return cdc.MustMarshal(&tokenIDWrap)
}

// return th tokenID
func MustUnMarshalTokenID(cdc codec.Codec, value []byte) string {
	var tokenIDWrap gogotypes.StringValue
	cdc.MustUnmarshal(value, &tokenIDWrap)
	return tokenIDWrap.Value
}

func MustUnMarshalTotalNftCountForCollection(cdc codec.Codec, value []byte) uint64 {
	var totalCountWrap gogotypes.UInt64Value
	cdc.MustUnmarshal(value, &totalCountWrap)
	return totalCountWrap.Value
}

func MustMarshallTotalCountForCollection(cdc codec.Codec, supply uint64) []byte {
	totalCountWrap := gogotypes.UInt64Value{Value: supply}
	return cdc.MustMarshal(&totalCountWrap)
}
