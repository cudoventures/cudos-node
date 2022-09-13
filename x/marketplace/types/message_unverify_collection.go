package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUnverifyCollection = "unverify_collection"

var _ sdk.Msg = &MsgUnverifyCollection{}

func NewMsgUnverifyCollection(creator string, id uint64) *MsgUnverifyCollection {
	return &MsgUnverifyCollection{
		Creator: creator,
		Id:      id,
	}
}

func (msg *MsgUnverifyCollection) Route() string {
	return RouterKey
}

func (msg *MsgUnverifyCollection) Type() string {
	return TypeMsgUnverifyCollection
}

func (msg *MsgUnverifyCollection) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUnverifyCollection) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUnverifyCollection) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
