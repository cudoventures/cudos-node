package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateAddress = "create_address"
	TypeMsgUpdateAddress = "update_address"
	TypeMsgDeleteAddress = "delete_address"
)

var _ sdk.Msg = &MsgCreateAddress{}

func NewMsgCreateAddress(creator, network, label, value string) *MsgCreateAddress {
	return &MsgCreateAddress{
		Creator: creator,
		Network: network,
		Label:   label,
		Value:   value,
	}
}

func (msg *MsgCreateAddress) Route() string {
	return RouterKey
}

func (msg *MsgCreateAddress) Type() string {
	return TypeMsgCreateAddress
}

func (msg *MsgCreateAddress) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateAddress) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateAddress) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateAddress{}

func NewMsgUpdateAddress(creator, network, label, value string) *MsgUpdateAddress {
	return &MsgUpdateAddress{
		Creator: creator,
		Network: network,
		Label:   label,
		Value:   value,
	}
}

func (msg *MsgUpdateAddress) Route() string {
	return RouterKey
}

func (msg *MsgUpdateAddress) Type() string {
	return TypeMsgUpdateAddress
}

func (msg *MsgUpdateAddress) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateAddress) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateAddress) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeleteAddress{}

func NewMsgDeleteAddress(creator, network, label string) *MsgDeleteAddress {
	return &MsgDeleteAddress{
		Creator: creator,
		Network: network,
		Label:   label,
	}
}
func (msg *MsgDeleteAddress) Route() string {
	return RouterKey
}

func (msg *MsgDeleteAddress) Type() string {
	return TypeMsgDeleteAddress
}

func (msg *MsgDeleteAddress) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteAddress) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteAddress) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}
