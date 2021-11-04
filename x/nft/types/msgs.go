package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// constant used to indicate that some field should not be updated
const (
	TypeMsgIssueDenom  = "issue_denom"
	TypeMsgTransferNft = "transfer_nft"
	TypeMsgEditNFT     = "edit_nft"
	TypeMsgMintNFT     = "mint_nft"
	TypeMsgBurnNFT     = "burn_nft"
	TypeMsgApproveNft  = "approve_nft"
	TypeMsgRevokeNft   = "revoke_nft"
)

var (
	_ sdk.Msg = &MsgIssueDenom{}
	_ sdk.Msg = &MsgTransferNft{}
	_ sdk.Msg = &MsgEditNFT{}
	_ sdk.Msg = &MsgMintNFT{}
	_ sdk.Msg = &MsgBurnNFT{}
	_ sdk.Msg = &MsgApproveNft{}
	_ sdk.Msg = &MsgRevokeNft{}
)

// NewMsgIssueDenom is a constructor function for MsgSetName
func NewMsgIssueDenom(denomID, denomName, schema, sender string) *MsgIssueDenom {
	return &MsgIssueDenom{
		Sender: sender,
		Id:     denomID,
		Name:   denomName,
		Schema: schema,
	}
}

// Route Implements Msg
func (msg MsgIssueDenom) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgIssueDenom) Type() string { return TypeMsgIssueDenom }

// ValidateBasic Implements Msg.
func (msg MsgIssueDenom) ValidateBasic() error {
	if err := ValidateDenomID(msg.Id); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return ValidateDenomName(msg.Name)
}

// GetSignBytes Implements Msg.
func (msg MsgIssueDenom) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgIssueDenom) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// NewMsgTransferNft is a constructor function for MsgSetName
func NewMsgTransferNft(
	tokenID, denomID, from, to, msgSender string,
) *MsgTransferNft {
	return &MsgTransferNft{
		TokenId: tokenID,
		DenomId: denomID,
		From:    from,
		To:      to,
		Sender:  msgSender,
	}
}

// Route Implements Msg
func (msg MsgTransferNft) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgTransferNft) Type() string { return TypeMsgTransferNft }

// ValidateBasic Implements Msg.
func (msg MsgTransferNft) ValidateBasic() error {
	if err := ValidateDenomID(msg.DenomId); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(msg.From); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.To); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid recipient address (%s)", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid recipient address (%s)", err)
	}

	return ValidateTokenID(msg.TokenId)
}

// GetSignBytes Implements Msg.
func (msg MsgTransferNft) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgTransferNft) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// NewMsgApproveAllNft NewMsgApproveNft is a constructor function for MsgSetName
func NewMsgApproveAllNft(operator, sender string, approved bool,
) *MsgApproveAllNft {
	return &MsgApproveAllNft{
		Operator: operator,
		Sender:   sender,
		Approved: approved,
	}
}

// Route Implements Msg
func (msg MsgApproveAllNft) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgApproveAllNft) Type() string { return TypeMsgApproveNft }

// ValidateBasic Implements Msg.
func (msg MsgApproveAllNft) ValidateBasic() error {

	if _, err := sdk.AccAddressFromBech32(msg.Operator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid operator address (%s)", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid operator address (%s)", err)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgApproveAllNft) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgApproveAllNft) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// NewMsgApproveNft is a constructor function for MsgSetName
func NewMsgApproveNft(tokenID, denomID, sender, approvedAddress string,
) *MsgApproveNft {
	return &MsgApproveNft{
		Id:              tokenID,
		DenomId:         denomID,
		Sender:          sender,
		ApprovedAddress: approvedAddress,
	}
}

// Route Implements Msg
func (msg MsgApproveNft) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgApproveNft) Type() string { return TypeMsgApproveNft }

// ValidateBasic Implements Msg.
func (msg MsgApproveNft) ValidateBasic() error {
	if err := ValidateDenomID(msg.DenomId); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.ApprovedAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	return ValidateTokenID(msg.Id)
}

// GetSignBytes Implements Msg.
func (msg MsgApproveNft) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgApproveNft) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// NewMsgRevokeNft is a constructor function for MsgSetName
func NewMsgRevokeNft(
	addressToApprove, sender, denomId, tokenId string,
) *MsgRevokeNft {
	return &MsgRevokeNft{
		AddressToRevoke: addressToApprove,
		Sender:          sender,
		DenomId:         denomId,
		TokenId:         tokenId,
	}
}

// Route Implements Msg
func (msg MsgRevokeNft) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgRevokeNft) Type() string { return TypeMsgRevokeNft }

// ValidateBasic Implements Msg.
func (msg MsgRevokeNft) ValidateBasic() error {
	if err := ValidateDenomID(msg.DenomId); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.AddressToRevoke); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid AddressToRevoke (%s)", err)
	}

	return ValidateTokenID(msg.TokenId)
}

// GetSignBytes Implements Msg.
func (msg MsgRevokeNft) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgRevokeNft) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// NewMsgEditNFT is a constructor function for MsgSetName
func NewMsgEditNFT(
	tokenID, denomID, tokenName, tokenURI, tokenData, sender string,
) *MsgEditNFT {
	return &MsgEditNFT{
		Id:      tokenID,
		DenomId: denomID,
		Name:    tokenName,
		URI:     tokenURI,
		Data:    tokenData,
		Sender:  sender,
	}
}

// Route Implements Msg
func (msg MsgEditNFT) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgEditNFT) Type() string { return TypeMsgEditNFT }

// ValidateBasic Implements Msg.
func (msg MsgEditNFT) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if err := ValidateDenomID(msg.DenomId); err != nil {
		return err
	}

	if err := ValidateTokenURI(msg.URI); err != nil {
		return err
	}
	return ValidateTokenID(msg.Id)
}

// GetSignBytes Implements Msg.
func (msg MsgEditNFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgEditNFT) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// NewMsgMintNFT is a constructor function for MsgMintNFT
func NewMsgMintNFT(denomID, tokenName, tokenURI, tokenData, sender, recipient string,
) *MsgMintNFT {
	return &MsgMintNFT{
		DenomId:   denomID,
		Name:      tokenName,
		URI:       tokenURI,
		Data:      tokenData,
		Sender:    sender,
		Recipient: recipient,
	}
}

// Route Implements Msg
func (msg MsgMintNFT) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgMintNFT) Type() string { return TypeMsgMintNFT }

// ValidateBasic Implements Msg.
func (msg MsgMintNFT) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if _, err := sdk.AccAddressFromBech32(msg.Recipient); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receipt address (%s)", err)
	}
	if err := ValidateDenomID(msg.DenomId); err != nil {
		return err
	}
	if err := ValidateTokenURI(msg.URI); err != nil {
		return err
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgMintNFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgMintNFT) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// NewMsgBurnNFT is a constructor function for MsgBurnNFT
func NewMsgBurnNFT(sender, tokenID, denomID string) *MsgBurnNFT {
	return &MsgBurnNFT{
		Sender:  sender,
		Id:      tokenID,
		DenomId: denomID,
	}
}

// Route Implements Msg
func (msg MsgBurnNFT) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgBurnNFT) Type() string { return TypeMsgBurnNFT }

// ValidateBasic Implements Msg.
func (msg MsgBurnNFT) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if err := ValidateDenomID(msg.DenomId); err != nil {
		return err
	}
	return ValidateTokenID(msg.Id)
}

// GetSignBytes Implements Msg.
func (msg MsgBurnNFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgBurnNFT) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}
