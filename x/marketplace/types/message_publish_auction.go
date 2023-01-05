package types

import (
	"time"

	nft "github.com/CudoVentures/cudos-node/x/nft/types"
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

func NewMsgPublishAuction(creator string, denomId string, tokenId string, duration time.Duration, auctionType AuctionType) (*MsgPublishAuction, error) {
	msg := &MsgPublishAuction{
		Creator:  creator,
		TokenId:  tokenId,
		DenomId:  denomId,
		Duration: duration,
	}

	err := msg.SetAuctionType(auctionType)
	return msg, err
}

func (msg *MsgPublishAuction) GetAuctionType() (AuctionType, error) {
	auctionType, ok := msg.AuctionType.GetCachedValue().(AuctionType)
	if !ok {
		return nil, sdkerrors.ErrInvalidType.Wrapf("expected %T, got %T", (AuctionType)(nil), msg.AuctionType.GetCachedValue())
	}
	return auctionType, nil
}

func (msg *MsgPublishAuction) SetAuctionType(auctionType AuctionType) error {
	any, err := types.NewAnyWithValue(auctionType)
	if err != nil {
		return err
	}
	msg.AuctionType = any
	return nil
}

func (msg MsgPublishAuction) UnpackInterfaces(unpacker types.AnyUnpacker) error {
	var auctionType AuctionType
	return unpacker.UnpackAny(msg.AuctionType, &auctionType)
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
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address")
	}

	if nft.ValidateDenomID(msg.DenomId) != nil || nft.ValidateTokenID(msg.TokenId) != nil {
		// todo add errors here and errytwhere
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid nft")
	}

	if msg.Duration < time.Hour*24 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "duration must be atleast 24 hours")
	}

	auctionType, err := msg.GetAuctionType()
	if err != nil {
		return sdkerrors.Wrap(err, "invalid auction type")
	}

	if err := auctionType.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "invalid auction type")
	}

	return nil
}
