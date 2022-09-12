package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgVerifyCollection = "verify_collection"

var _ sdk.Msg = &MsgVerifyCollection{}

func NewMsgVerifyCollection(creator string, id uint64) *MsgVerifyCollection {
	return &MsgVerifyCollection{
		Creator: creator,
		Id:      id,
	}
}

func (msg *MsgVerifyCollection) Route() string {
	return RouterKey
}

func (msg *MsgVerifyCollection) Type() string {
	return TypeMsgVerifyCollection
}

func (msg *MsgVerifyCollection) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgVerifyCollection) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgVerifyCollection) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
