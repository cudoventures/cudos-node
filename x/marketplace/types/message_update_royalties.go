package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateRoyalties = "update_royalties"

var _ sdk.Msg = &MsgUpdateRoyalties{}

func NewMsgUpdateRoyalties(creator string, id uint64, mintRoyalties, resaleRoyalties []Royalty) *MsgUpdateRoyalties {
	return &MsgUpdateRoyalties{
		Creator:         creator,
		Id:              id,
		MintRoyalties:   mintRoyalties,
		ResaleRoyalties: resaleRoyalties,
	}
}

func (msg *MsgUpdateRoyalties) Route() string {
	return RouterKey
}

func (msg *MsgUpdateRoyalties) Type() string {
	return TypeMsgUpdateRoyalties
}

func (msg *MsgUpdateRoyalties) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateRoyalties) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateRoyalties) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if err := ValidateMintRoyalties(msg.MintRoyalties); err != nil {
		return err
	}

	return ValidateResaleRoyalties(msg.ResaleRoyalties)
}
