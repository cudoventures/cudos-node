package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgPlaceBid = "place_bid"

var _ sdk.Msg = &MsgPlaceBid{}

func NewMsgPlaceBid(bidder string, auctionId uint64, amount sdk.Coin) *MsgPlaceBid {
	return &MsgPlaceBid{
		Bidder:    bidder,
		AuctionId: auctionId,
		Amount:    amount,
	}
}

func (msg *MsgPlaceBid) Route() string {
	return RouterKey
}

func (msg *MsgPlaceBid) Type() string {
	return TypeMsgPlaceBid
}

func (msg *MsgPlaceBid) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Bidder)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgPlaceBid) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgPlaceBid) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Bidder); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid bidder: %s", err)
	}

	if err := msg.Amount.Validate(); err != nil {
		return sdkerrors.Wrapf(ErrInvalidPrice, "invalid amount: %s", err)
	}

	if msg.Amount.IsZero() {
		return sdkerrors.Wrapf(ErrInvalidPrice, "amount must be positive")
	}

	return nil
}
