package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgSend{}

// Sends message when NFT ownership is changed
func NewMsgSend(creator string, id string, receiver string) *MsgSend {
	return &MsgSend{
		Creator:  creator,
		Id:       id,
		Receiver: receiver,
	}
}

func (msg *MsgSend) Route() string {
	return RouterKey
}

func (msg *MsgSend) Type() string {
	return "Send"
}

func (msg *MsgSend) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSend) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSend) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
