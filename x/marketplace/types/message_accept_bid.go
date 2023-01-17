package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgAcceptBid = "accept_bid"

var _ sdk.Msg = &MsgAcceptBid{}

func NewMsgAcceptBid(sender string, auctionId uint64) *MsgAcceptBid {
	return &MsgAcceptBid{
		Sender:    sender,
		AuctionId: auctionId,
	}
}

func (msg *MsgAcceptBid) Route() string {
	return RouterKey
}

func (msg *MsgAcceptBid) Type() string {
	return TypeMsgAcceptBid
}

func (msg *MsgAcceptBid) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgAcceptBid) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAcceptBid) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender: %s", err)
	}

	return nil
}
