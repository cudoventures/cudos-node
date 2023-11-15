package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgPublishNft = "publish_nft"

var _ sdk.Msg = &MsgPublishNft{}

func NewMsgPublishNft(creator, tokenId, denomId string, price sdk.Coin) *MsgPublishNft {
	return &MsgPublishNft{
		Creator: creator,
		TokenId: tokenId,
		DenomId: denomId,
		Price:   price,
	}
}

func (*MsgPublishNft) Route() string {
	return RouterKey
}

func (*MsgPublishNft) Type() string {
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

	if msg.Price.Amount.Equal(sdk.NewInt(0)) {
		return sdkerrors.Wrapf(ErrInvalidPrice, "invalid price (%s)", msg.Price)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
