package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgPublishNft = "publish_nft"

var _ sdk.Msg = &MsgPublishNft{}

func NewMsgPublishNft(creator, tokenId, denomId, price string) *MsgPublishNft {
	return &MsgPublishNft{
		Creator: creator,
		TokenId: tokenId,
		DenomId: denomId,
		Price:   price,
	}
}

func (msg *MsgPublishNft) Route() string {
	return RouterKey
}

func (msg *MsgPublishNft) Type() string {
	return TypeMsgPublishNft
}

func (msg *MsgPublishNft) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgPublishNft) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgPublishNft) ValidateBasic() error {
	if msg.TokenId == "" {
		return sdkerrors.Wrap(ErrEmptyNftID, "empty nft id")
	}

	if msg.DenomId == "" {
		return sdkerrors.Wrap(ErrEmptyDenomID, "empty denom id")
	}

	if _, err := sdk.ParseCoinNormalized(msg.Price); err != nil {
		return sdkerrors.Wrapf(ErrInvalidPrice, "invalid price (%s)", msg.Price)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
