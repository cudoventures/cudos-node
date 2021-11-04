package types

// DONTCOVER

import (
	gogotypes "github.com/gogo/protobuf/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"

	"cudos.org/cudos-node/x/nft/exported"
)

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}

// RegisterLegacyAminoCodec concrete types on codec
// (Amino is still needed for Ledger at the moment)
// nolint: staticcheck
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgIssueDenom{}, "cudos.org/cudos-node/nft/MsgIssueDenom", nil)
	cdc.RegisterConcrete(&MsgTransferNft{}, "cudos.org/cudos-node/nft/MsgTransferNft", nil)
	cdc.RegisterConcrete(&MsgApproveNft{}, "cudos.org/cudos-node/nft/MsgApproveNft", nil)
	cdc.RegisterConcrete(&MsgRevokeNft{}, "cudos.org/cudos-node/nft/MsgRevokeNft", nil)
	cdc.RegisterConcrete(&MsgEditNFT{}, "cudos.org/cudos-node/nft/MsgEditNFT", nil)
	cdc.RegisterConcrete(&MsgMintNFT{}, "cudos.org/cudos-node/nft/MsgMintNFT", nil)
	cdc.RegisterConcrete(&MsgBurnNFT{}, "cudos.org/cudos-node/nft/MsgBurnNFT", nil)

	cdc.RegisterInterface((*exported.NFT)(nil), nil)
	cdc.RegisterConcrete(&BaseNFT{}, "cudos.org/cudos-node/nft/BaseNFT", nil)
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
