package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgSendMessage = "send_message"
)

type MsgSendMessage struct {
	FromAddress sdk.AccAddress `json:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address"`
	Subject     string         `json:"subject"`
	Body        string         `json:"body"`
}

// ProtoMessage implements proto.Message.
func (msg *MsgSendMessage) ProtoMessage() {
	panic("unimplemented")
}

// Reset implements proto.Message.
func (msg *MsgSendMessage) Reset() {
	panic("unimplemented")
}

// String implements proto.Message.
func (msg *MsgSendMessage) String() string {
	panic("unimplemented")
}

// NewMsgSendMessage is a constructor function for MsgSendMessage
func NewMsgSendMessage(fromAddr sdk.AccAddress, toAddr sdk.AccAddress, subject, body string) MsgSendMessage {
	return MsgSendMessage{
		FromAddress: fromAddr,
		ToAddress:   toAddr,
		Subject:     subject,
		Body:        body,
	}
}

func (msg MsgSendMessage) Route() string { return RouterKey }
func (msg MsgSendMessage) Type() string  { return TypeMsgSendMessage }
func (msg MsgSendMessage) ValidateBasic() error {
	if msg.FromAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing from address")
	}
	if msg.ToAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing to address")
	}
	if len(msg.Subject) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "missing subject")
	}
	if len(msg.Body) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "missing body")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSendMessage) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners defines whose signature is required
func (msg MsgSendMessage) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}
