package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// constant used to indicate that some field should not be updated
const (
	TypeMsgIssueDenom    = "issue_denom"
	TypeMsgTransferNft   = "transfer_nft"
	TypeMsgEditNFT       = "edit_nft"
	TypeMsgMintNFT       = "mint_nft"
	TypeMsgBurnNFT       = "burn_nft"
	TypeMsgApproveNft    = "approve_nft"
	TypeMsgRevokeNft     = "revoke_nft"
	TypeMsgApproveAllNft = "approve_all"
)

var (
	_ sdk.Msg = &MsgIssueDenom{}
	_ sdk.Msg = &MsgTransferNft{}
	_ sdk.Msg = &MsgEditNFT{}
	_ sdk.Msg = &MsgMintNFT{}
	_ sdk.Msg = &MsgBurnNFT{}
	_ sdk.Msg = &MsgApproveNft{}
	_ sdk.Msg = &MsgRevokeNft{}
	_ sdk.Msg = &MsgApproveAllNft{}
)

// NewMsgIssueDenom is a constructor function for MsgIssueDenom
func NewMsgIssueDenom(denomID, denomName, schema, sender, contractAddressSigner, symbol string) *MsgIssueDenom {
	return &MsgIssueDenom{
		Sender:                sender,
		Id:                    denomID,
		Name:                  denomName,
		Schema:                schema,
		ContractAddressSigner: contractAddressSigner, // field is only populated when the request is coming from a contract, in other cases its empty string
		Symbol:                symbol,
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

	if err := ValidateDenomName(msg.Name); err != nil {
		return err
	}

	if err := ValidateDenomSymbol(msg.Symbol); err != nil {
		return err
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgIssueDenom) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgIssueDenom) GetSigners() []sdk.AccAddress {
	var signers []sdk.AccAddress

	if msg.ContractAddressSigner != "" {
		contractAddressSigner, err := sdk.AccAddressFromBech32(msg.ContractAddressSigner)
		if err != nil {
			panic(err)
		}
		signers = append(signers, contractAddressSigner)
		return signers
	}

	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	signers = append(signers, from)
	return signers

}

// NewMsgTransferNft is a constructor function for MsgTransferNft
func NewMsgTransferNft(
	denomID, tokenID, from, to, msgSender, contractAddressSigner string,
) *MsgTransferNft {
	return &MsgTransferNft{
		DenomId:               denomID,
		TokenId:               tokenID,
		From:                  from,
		To:                    to,
		Sender:                msgSender,
		ContractAddressSigner: contractAddressSigner,
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
	var signers []sdk.AccAddress

	if msg.ContractAddressSigner != "" {
		contractAddressSigner, err := sdk.AccAddressFromBech32(msg.ContractAddressSigner)
		if err != nil {
			panic(err)
		}
		signers = append(signers, contractAddressSigner)
		return signers
	}

	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	signers = append(signers, from)
	return signers
}

// NewMsgApproveAllNft NewMsgApproveNft is a constructor function for MsgApproveAllNft
func NewMsgApproveAllNft(operator, sender, contractAddressSigner string, approved bool,
) *MsgApproveAllNft {
	return &MsgApproveAllNft{
		Operator:              operator,
		Sender:                sender,
		Approved:              approved,
		ContractAddressSigner: contractAddressSigner,
	}
}

// Route Implements Msg
func (msg MsgApproveAllNft) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgApproveAllNft) Type() string { return TypeMsgApproveAllNft }

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
	var signers []sdk.AccAddress

	if msg.ContractAddressSigner != "" {
		contractAddressSigner, err := sdk.AccAddressFromBech32(msg.ContractAddressSigner)
		if err != nil {
			panic(err)
		}
		signers = append(signers, contractAddressSigner)
		return signers
	}

	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	signers = append(signers, from)
	return signers
}

// NewMsgApproveNft is a constructor function for MsgApproveNft
func NewMsgApproveNft(tokenID, denomID, sender, approvedAddress, contractAddressSigner string,
) *MsgApproveNft {
	return &MsgApproveNft{
		Id:                    tokenID,
		DenomId:               denomID,
		Sender:                sender,
		ApprovedAddress:       approvedAddress,
		ContractAddressSigner: contractAddressSigner,
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
	var signers []sdk.AccAddress

	if msg.ContractAddressSigner != "" {
		contractAddressSigner, err := sdk.AccAddressFromBech32(msg.ContractAddressSigner)
		if err != nil {
			panic(err)
		}
		signers = append(signers, contractAddressSigner)
		return signers
	}

	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	signers = append(signers, from)
	return signers
}

// NewMsgRevokeNft is a constructor function for MsgRevokeNft
func NewMsgRevokeNft(
	addressToRevoke, sender, denomId, tokenId, contractAddressSigner string,
) *MsgRevokeNft {
	return &MsgRevokeNft{
		AddressToRevoke:       addressToRevoke,
		Sender:                sender,
		DenomId:               denomId,
		TokenId:               tokenId,
		ContractAddressSigner: contractAddressSigner,
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
	var signers []sdk.AccAddress

	if msg.ContractAddressSigner != "" {
		contractAddressSigner, err := sdk.AccAddressFromBech32(msg.ContractAddressSigner)
		if err != nil {
			panic(err)
		}
		signers = append(signers, contractAddressSigner)
		return signers
	}

	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	signers = append(signers, from)
	return signers
}

// NewMsgEditNFT is a constructor function for MsgSetName
func NewMsgEditNFT(
	tokenID, denomID, tokenName, tokenURI, tokenData, sender, contractAddressSigner string,
) *MsgEditNFT {
	return &MsgEditNFT{
		Id:                    tokenID,
		DenomId:               denomID,
		Name:                  tokenName,
		URI:                   tokenURI,
		Data:                  tokenData,
		Sender:                sender,
		ContractAddressSigner: contractAddressSigner, // field is only populated when the request is coming from and is signed by a contract, in other cases its empty string
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
	var signers []sdk.AccAddress

	// this order is important
	// always check first if the request is signed by a contract
	// the wasmd module requires that only the contract address is the signer
	if msg.ContractAddressSigner != "" {
		contractAddressSigner, err := sdk.AccAddressFromBech32(msg.ContractAddressSigner)
		if err != nil {
			panic(err)
		}
		signers = append(signers, contractAddressSigner)
		return signers
	}

	// if no contract address - it was issued directly by the user via CLI/REST
	// so the sender is the signer
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	signers = append(signers, from)
	return signers
}

// NewMsgMintNFT is a constructor function for MsgMintNFT
func NewMsgMintNFT(denomID, tokenName, tokenURI, tokenData, sender, recipient, contractAddressSigner string,
) *MsgMintNFT {
	return &MsgMintNFT{
		DenomId:               denomID,
		Name:                  tokenName,
		URI:                   tokenURI,
		Data:                  tokenData,
		Sender:                sender,
		Recipient:             recipient,
		ContractAddressSigner: contractAddressSigner, // field is only populated when the request is coming from a contract, in other cases it should be an empty string
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
	var signers []sdk.AccAddress

	if msg.ContractAddressSigner != "" {
		contractAddressSigner, err := sdk.AccAddressFromBech32(msg.ContractAddressSigner)
		if err != nil {
			panic(err)
		}
		signers = append(signers, contractAddressSigner)
		return signers
	}

	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	signers = append(signers, from)
	return signers
}

// NewMsgBurnNFT is a constructor function for MsgBurnNFT
func NewMsgBurnNFT(sender, tokenID, denomID, contractAddressSigner string) *MsgBurnNFT {
	return &MsgBurnNFT{
		Sender:                sender,
		Id:                    tokenID,
		DenomId:               denomID,
		ContractAddressSigner: contractAddressSigner,
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
	var signers []sdk.AccAddress

	if msg.ContractAddressSigner != "" {
		contractAddressSigner, err := sdk.AccAddressFromBech32(msg.ContractAddressSigner)
		if err != nil {
			panic(err)
		}
		signers = append(signers, contractAddressSigner)
		return signers
	}

	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	signers = append(signers, from)
	return signers
}
