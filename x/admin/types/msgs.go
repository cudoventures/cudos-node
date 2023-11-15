package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// bank message types
const (
	TypeMsgSend      = "adminSpendCommunityPool"
	TypeMsgMultiSend = "multisend"
)

var _ sdk.Msg = &MsgAdminSpendCommunityPool{}

// NewMsgSend - construct a msg to send coins from one account to another.
//
//nolint:interfacer
func NewMsgAdminSpendCommunityPool(fromAddr, toAddr sdk.AccAddress, amount sdk.Coins) *MsgAdminSpendCommunityPool {
	return &MsgAdminSpendCommunityPool{Initiator: fromAddr.String(), ToAddress: toAddr.String(), Coins: amount}
}

// Route Implements Msg.
func (MsgAdminSpendCommunityPool) Route() string { return RouterKey }

// Type Implements Msg.
func (MsgAdminSpendCommunityPool) Type() string { return TypeMsgSend }

// ValidateBasic Implements Msg.
func (msg MsgAdminSpendCommunityPool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Initiator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.ToAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid recipient address (%s)", err)
	}

	if !msg.Coins.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Coins.String())
	}

	if !msg.Coins.IsAllPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Coins.String())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgAdminSpendCommunityPool) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgAdminSpendCommunityPool) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Initiator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}
