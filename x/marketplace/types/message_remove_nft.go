package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRemoveNft = "remove_nft"

var _ sdk.Msg = &MsgRemoveNft{}

func NewMsgRemoveNft(creator string, id uint64) *MsgRemoveNft {
	return &MsgRemoveNft{
		Creator: creator,
		Id:      id,
	}
}

func (msg *MsgRemoveNft) Route() string {
	return RouterKey
}

func (msg *MsgRemoveNft) Type() string {
	return TypeMsgRemoveNft
}

func (msg *MsgRemoveNft) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRemoveNft) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRemoveNft) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
