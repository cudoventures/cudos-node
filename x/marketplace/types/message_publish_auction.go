package types

import (
	"time"

	nft "github.com/CudoVentures/cudos-node/x/nft/types"
	nfttypes "github.com/CudoVentures/cudos-node/x/nft/types"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgPublishAuction = "publish_auction"
)

var (
	_ sdk.Msg = &MsgPublishAuction{}
)

func NewMsgPublishAuction(creator string, denomId string, tokenId string, duration time.Duration, at AuctionType) (*MsgPublishAuction, error) {
	msg := &MsgPublishAuction{
		Creator:  creator,
		TokenId:  tokenId,
		DenomId:  denomId,
		Duration: duration,
	}

	err := msg.SetAuctionType(at)
	return msg, err
}

func (msg *MsgPublishAuction) GetAuctionType() (AuctionType, error) {
	at, ok := msg.AuctionType.GetCachedValue().(AuctionType)
	if !ok {
		return nil, sdkerrors.ErrInvalidType.Wrapf("expected %T, got %T", (AuctionType)(nil), msg.AuctionType.GetCachedValue())
	}
	return at, nil
}

func (msg *MsgPublishAuction) SetAuctionType(at AuctionType) error {
	any, err := types.NewAnyWithValue(at)
	if err != nil {
		return err
	}
	msg.AuctionType = any
	return nil
}

func (msg MsgPublishAuction) UnpackInterfaces(unpacker types.AnyUnpacker) error {
	var at AuctionType
	return unpacker.UnpackAny(msg.AuctionType, &at)
}

func (msg *MsgPublishAuction) Route() string {
	return RouterKey
}

func (msg *MsgPublishAuction) Type() string {
	return TypeMsgPublishAuction
}

func (msg *MsgPublishAuction) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgPublishAuction) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgPublishAuction) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.ErrInvalidAddress
	}

	if nft.ValidateDenomID(msg.DenomId) != nil || nft.ValidateTokenID(msg.TokenId) != nil {
		return nfttypes.ErrInvalidNFT
	}

	if msg.Duration < time.Hour*24 {
		return ErrInvalidAuctionDuration
	}

	at, err := msg.GetAuctionType()
	if err != nil {
		return err
	}

	return at.ValidateBasic()
}
