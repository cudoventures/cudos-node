package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgTransferAdminPermission = "transfer_admin_permission"

var _ sdk.Msg = &MsgTransferAdminPermission{}

func NewMsgTransferAdminPermission(creator string, newAdmin string) *MsgTransferAdminPermission {
	return &MsgTransferAdminPermission{
		Creator:  creator,
		NewAdmin: newAdmin,
	}
}

func (msg *MsgTransferAdminPermission) Route() string {
	return RouterKey
}

func (msg *MsgTransferAdminPermission) Type() string {
	return TypeMsgTransferAdminPermission
}

func (msg *MsgTransferAdminPermission) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgTransferAdminPermission) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgTransferAdminPermission) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if _, err := sdk.AccAddressFromBech32(msg.NewAdmin); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid new admin address (%s)", err)
	}
	if msg.Creator == msg.NewAdmin {
		return sdkerrors.Wrap(ErrAlreadyAdmin, "cannot transfer admin permission to yourself")
	}
	return nil
}
