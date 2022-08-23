package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateToken = "create_token"
	TypeMsgUpdateToken = "update_token"
	TypeMsgDeleteToken = "delete_token"
)

var _ sdk.Msg = &MsgCreateToken{}

func NewMsgCreateToken(
	owner string,
	denom string,
	name string,
	decimals uint64,
	initialBalances []*Balance,
	maxSupply uint64,

) *MsgCreateToken {
	return &MsgCreateToken{
		Owner:           owner,
		Denom:           denom,
		Name:            name,
		Decimals:        decimals,
		InitialBalances: initialBalances,
		MaxSupply:       maxSupply,
	}
}

func (msg *MsgCreateToken) Route() string {
	return RouterKey
}

func (msg *MsgCreateToken) Type() string {
	return TypeMsgCreateToken
}

func (msg *MsgCreateToken) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

func (msg *MsgCreateToken) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateToken) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}
	return nil
}
