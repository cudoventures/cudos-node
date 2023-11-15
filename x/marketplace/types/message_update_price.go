package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdatePrice = "update_price"

var _ sdk.Msg = &MsgUpdatePrice{}

func NewMsgUpdatePrice(creator string, id uint64, price sdk.Coin) *MsgUpdatePrice {
	return &MsgUpdatePrice{
		Creator: creator,
		Id:      id,
		Price:   price,
	}
}

func (*MsgUpdatePrice) Route() string {
	return RouterKey
}

func (*MsgUpdatePrice) Type() string {
	return TypeMsgUpdatePrice
}

func (msg *MsgUpdatePrice) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdatePrice) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdatePrice) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Price.Amount.Equal(sdk.NewInt(0)) {
		return sdkerrors.Wrapf(ErrInvalidPrice, "invalid price (%s)", msg.Price)
	}

	return nil
}
