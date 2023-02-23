package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRemoveAdmin = "remove_admin"

var _ sdk.Msg = &MsgRemoveAdmin{}

func NewMsgRemoveAdmin(creator string, address string) *MsgRemoveAdmin {
	return &MsgRemoveAdmin{
		Creator: creator,
		Address: address,
	}
}

func (msg *MsgRemoveAdmin) Route() string {
	return RouterKey
}

func (msg *MsgRemoveAdmin) Type() string {
	return TypeMsgRemoveAdmin
}

func (msg *MsgRemoveAdmin) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRemoveAdmin) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRemoveAdmin) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if _, err := sdk.AccAddressFromBech32(msg.Address); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid admin address (%s)", err)
	}
	return nil
}
