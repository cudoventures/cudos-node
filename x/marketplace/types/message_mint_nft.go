package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgMintNft = "mint_nft"

var _ sdk.Msg = &MsgMintNft{}

func NewMsgMintNft(creator, denomId, recipient, name, uri, data string, price sdk.Coin) *MsgMintNft {
	return &MsgMintNft{
		Creator:   creator,
		DenomId:   denomId,
		Recipient: recipient,
		Price:     price,
		Name:      name,
		Uri:       uri,
		Data:      data,
	}
}

func (msg *MsgMintNft) Route() string {
	return RouterKey
}

func (msg *MsgMintNft) Type() string {
	return TypeMsgMintNft
}

func (msg *MsgMintNft) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgMintNft) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgMintNft) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if _, err := sdk.AccAddressFromBech32(msg.Recipient); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid recipient address (%s)", err)
	}
	if msg.DenomId == "" {
		return sdkerrors.Wrap(ErrEmptyDenomID, "empty denom id")
	}
	if msg.Price.Amount.Equal(sdk.NewInt(0)) {
		return sdkerrors.Wrapf(ErrInvalidPrice, "invalid price (%+v)", msg.Price)
	}

	return nil
}
