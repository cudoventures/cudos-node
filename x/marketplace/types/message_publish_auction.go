package types

import (
	"time"

	nfttypes "github.com/CudoVentures/cudos-node/x/nft/types"
	"github.com/cosmos/cosmos-sdk/codec/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgPublishAuction = "publish_auction"
)

var (
	_ sdk.Msg                            = &MsgPublishAuction{}
	_ codectypes.UnpackInterfacesMessage = MsgPublishAuction{}
)

func NewMsgPublishAuction(
	creator string,
	denomId string, tokenId string,
	duration time.Duration,
	a Auction,
) (*MsgPublishAuction, error) {
	msg := &MsgPublishAuction{
		Creator:  creator,
		TokenId:  tokenId,
		DenomId:  denomId,
		Duration: duration,
	}

	err := msg.SetAuction(a)
	return msg, err
}

func (msg *MsgPublishAuction) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator %s", err)
	}

	if err := nfttypes.ValidateDenomID(msg.DenomId); err != nil {
		return nfttypes.ErrInvalidDenom
	}

	if err := nfttypes.ValidateTokenID(msg.TokenId); err != nil {
		return nfttypes.ErrInvalidTokenID
	}

	if msg.Duration < time.Hour*24 {
		return sdkerrors.Wrap(ErrInvalidAuctionDuration, "duration is less than 24 hours")
	}

	a, err := msg.GetAuction()
	if err != nil {
		return err
	}

	return a.ValidateBasic()
}

func (msg *MsgPublishAuction) GetAuction() (Auction, error) {
	return UnpackAuction(msg.Auction)
}

func (msg *MsgPublishAuction) SetAuction(a Auction) error {
	any, err := PackAuction(a)
	if err != nil {
		return err
	}

	msg.Auction = any
	return nil
}

func (msg MsgPublishAuction) UnpackInterfaces(unpacker types.AnyUnpacker) error {
	var at Auction
	return unpacker.UnpackAny(msg.Auction, &at)
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
