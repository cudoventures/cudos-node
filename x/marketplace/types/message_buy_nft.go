package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgBuyNft = "buy_nft"

var _ sdk.Msg = &MsgBuyNft{}

func NewMsgBuyNft(creator string, id uint64) *MsgBuyNft {
	return &MsgBuyNft{
		Creator: creator,
		Id:      id,
	}
}

func (msg *MsgBuyNft) Route() string {
	return RouterKey
}

func (msg *MsgBuyNft) Type() string {
	return TypeMsgBuyNft
}

func (msg *MsgBuyNft) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgBuyNft) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgBuyNft) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
