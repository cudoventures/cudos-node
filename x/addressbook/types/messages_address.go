package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateAddress = "create_address"
	TypeMsgUpdateAddress = "update_address"
	TypeMsgDeleteAddress = "delete_address"

	MinNetworkLength = 1
	MaxNetworkLength = 256
	MinLabelLength   = 1
	MaxLabelLength   = 256
	MinValueLength   = 1
	MaxValueLength   = 256
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

func (*MsgCreateAddress) Route() string {
	return RouterKey
}

func (*MsgCreateAddress) Type() string {
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
	if err := validateNetwork(msg.Network); err != nil {
		return err
	}
	if err := validateLabel(msg.Label); err != nil {
		return err
	}
	if err := validateValue(msg.Value); err != nil {
		return err
	}
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

func (*MsgUpdateAddress) Route() string {
	return RouterKey
}

func (*MsgUpdateAddress) Type() string {
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
	if err := validateNetwork(msg.Network); err != nil {
		return err
	}
	if err := validateLabel(msg.Label); err != nil {
		return err
	}
	if err := validateValue(msg.Value); err != nil {
		return err
	}
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
	if err := validateNetwork(msg.Network); err != nil {
		return err
	}
	if err := validateLabel(msg.Label); err != nil {
		return err
	}
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}

func validateNetwork(network string) error {
	if len(network) < MinNetworkLength || len(network) > MaxNetworkLength {
		return sdkerrors.Wrapf(ErrInvalidNetwork, "the length of network(%s) only accepts value [%d, %d]", network, MinNetworkLength, MaxNetworkLength)
	}
	return nil
}

func validateLabel(label string) error {
	if len(label) < MinLabelLength || len(label) > MaxLabelLength {
		return sdkerrors.Wrapf(ErrInvalidNetwork, "the length of label(%s) only accepts value [%d, %d]", label, MinLabelLength, MaxLabelLength)
	}
	return nil
}

func validateValue(value string) error {
	if len(value) < MinValueLength || len(value) > MaxValueLength {
		return sdkerrors.Wrapf(ErrInvalidNetwork, "the length of value(%s) only accepts value [%d, %d]", value, MinValueLength, MaxValueLength)
	}
	return nil
}
