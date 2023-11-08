package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	nfttypes "github.com/CudoVentures/cudos-node/x/nft/types"
)

const TypeMsgCreateCollection = "create_collection"

var _ sdk.Msg = &MsgCreateCollection{}

func NewMsgCreateCollection(creator, id, name, schema, symbol, traits, description, minter, data string, mintRoyalties, resaleRoyalties []Royalty, verified bool) *MsgCreateCollection {
	return &MsgCreateCollection{
		Creator:         creator,
		Id:              id,
		Name:            name,
		Schema:          schema,
		Symbol:          symbol,
		Traits:          traits,
		Description:     description,
		Minter:          minter,
		Data:            data,
		MintRoyalties:   mintRoyalties,
		ResaleRoyalties: resaleRoyalties,
		Verified:        verified,
	}
}

func (msg *MsgCreateCollection) Route() string {
	return RouterKey
}

func (msg *MsgCreateCollection) Type() string {
	return TypeMsgCreateCollection
}

func (msg *MsgCreateCollection) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateCollection) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateCollection) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if err := nfttypes.ValidateDenomID(msg.Id); err != nil {
		return err
	}
	if err := nfttypes.ValidateDenomName(msg.Name); err != nil {
		return err
	}
	if err := nfttypes.ValidateSchema(msg.Schema); err != nil {
		return err
	}
	if err := nfttypes.ValidateDenomSymbol(msg.Symbol); err != nil {
		return err
	}
	if err := nfttypes.ValidateDenomTraits(msg.Traits); err != nil {
		return err
	}
	if err := nfttypes.ValidateDescription(msg.Description); err != nil {
		return err
	}
	if err := nfttypes.ValidateMinter(msg.Minter); err != nil {
		return err
	}
	if err := nfttypes.ValidateDenomData(msg.Data); err != nil {
		return err
	}

	if err := ValidateMintRoyalties(msg.MintRoyalties); err != nil {
		return err
	}

	return ValidateResaleRoyalties(msg.ResaleRoyalties)
}
