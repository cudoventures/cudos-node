package types

import (
	nfttypes "github.com/CudoVentures/cudos-node/x/nft/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgPublishCollection = "publish_collection"

var _ sdk.Msg = &MsgPublishCollection{}

func NewMsgPublishCollection(creator, denomId string, mintRoyalties, resaleRoyalties []Royalty) *MsgPublishCollection {
	return &MsgPublishCollection{
		Creator:         creator,
		DenomId:         denomId,
		MintRoyalties:   mintRoyalties,
		ResaleRoyalties: resaleRoyalties,
	}
}

func (msg *MsgPublishCollection) Route() string {
	return RouterKey
}

func (msg *MsgPublishCollection) Type() string {
	return TypeMsgPublishCollection
}

func (msg *MsgPublishCollection) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgPublishCollection) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgPublishCollection) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if err := nfttypes.ValidateDenomID(msg.DenomId); err != nil {
		return err
	}

	if err := ValidateMintRoyalties(msg.MintRoyalties); err != nil {
		return err
	}

	return ValidateResaleRoyalties(msg.ResaleRoyalties)
}
