package custom_bindings

import (
	"encoding/json"
	wasmKeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmd/x/wasm/types"
	nftTypes "github.com/CudoVentures/cudos-node/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// fromReflectRawMsg decodes msg.Data to an sdk.Msg using proto Any and json encoding.
// this needs to be registered on the Encoders
func EncodeNftMessage() wasmKeeper.CustomEncoder {
	return func(_sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, error) {

		var nftCustomMsg nftCustomMsg
		err := json.Unmarshal(msg, &nftCustomMsg)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
		}

		switch {
		case nftCustomMsg.IssueDenomMsg != nil:
			issueDenomMsg := nftTypes.NewMsgIssueDenom(
				nftCustomMsg.IssueDenomMsg.Id,
				nftCustomMsg.IssueDenomMsg.Name,
				nftCustomMsg.IssueDenomMsg.Schema,
				nftCustomMsg.IssueDenomMsg.Sender,
				nftCustomMsg.IssueDenomMsg.ContractAddressSigner,
				nftCustomMsg.IssueDenomMsg.Symbol,
			)
			return []sdk.Msg{issueDenomMsg}, nil
		case nftCustomMsg.MintNftMsg != nil:
			mintNftMsg := nftTypes.NewMsgMintNFT(
				nftCustomMsg.MintNftMsg.DenomId,
				nftCustomMsg.MintNftMsg.Name,
				nftCustomMsg.MintNftMsg.URI,
				nftCustomMsg.MintNftMsg.Data,
				nftCustomMsg.MintNftMsg.Sender,
				nftCustomMsg.MintNftMsg.Recipient,
				nftCustomMsg.MintNftMsg.ContractAddressSigner)
			return []sdk.Msg{mintNftMsg}, nil
		case nftCustomMsg.EditNftMsg != nil:
			editNftMsg := nftTypes.NewMsgEditNFT(
				nftCustomMsg.EditNftMsg.TokenId,
				nftCustomMsg.EditNftMsg.DenomId,
				nftCustomMsg.EditNftMsg.Name,
				nftCustomMsg.EditNftMsg.URI,
				nftCustomMsg.EditNftMsg.Data,
				nftCustomMsg.EditNftMsg.Sender,
				nftCustomMsg.EditNftMsg.ContractAddressSigner)
			return []sdk.Msg{editNftMsg}, nil
		case nftCustomMsg.TransferNftMsg != nil:
			transferNftMsg := nftTypes.NewMsgTransferNft(
				nftCustomMsg.TransferNftMsg.DenomId,
				nftCustomMsg.TransferNftMsg.TokenId,
				nftCustomMsg.TransferNftMsg.From,
				nftCustomMsg.TransferNftMsg.To,
				nftCustomMsg.TransferNftMsg.Sender,
				nftCustomMsg.TransferNftMsg.ContractAddressSigner)
			return []sdk.Msg{transferNftMsg}, nil
		case nftCustomMsg.BurnNftMsg != nil:
			burnNftMsg := nftTypes.NewMsgBurnNFT(
				nftCustomMsg.BurnNftMsg.Sender,
				nftCustomMsg.BurnNftMsg.TokenId,
				nftCustomMsg.BurnNftMsg.DenomId,
				nftCustomMsg.BurnNftMsg.ContractAddressSigner)
			return []sdk.Msg{burnNftMsg}, nil
		case nftCustomMsg.ApproveNftMsg != nil:
			approveNftMsg := nftTypes.NewMsgApproveNft(
				nftCustomMsg.ApproveNftMsg.TokenId,
				nftCustomMsg.ApproveNftMsg.DenomId,
				nftCustomMsg.ApproveNftMsg.Sender,
				nftCustomMsg.ApproveNftMsg.ApprovedAddress,
				nftCustomMsg.ApproveNftMsg.ContractAddressSigner)
			return []sdk.Msg{approveNftMsg}, nil
		case nftCustomMsg.ApproveAllMsg != nil:
			approveNftMsg := nftTypes.NewMsgApproveAllNft(
				nftCustomMsg.ApproveAllMsg.ApprovedOperator,
				nftCustomMsg.ApproveAllMsg.Sender,
				nftCustomMsg.ApproveAllMsg.ContractAddressSigner,
				nftCustomMsg.ApproveAllMsg.Approved)
			return []sdk.Msg{approveNftMsg}, nil
		case nftCustomMsg.RevokeApprovalMsg != nil:
			approveNftMsg := nftTypes.NewMsgRevokeNft(
				nftCustomMsg.RevokeApprovalMsg.AddressToRevoke,
				nftCustomMsg.RevokeApprovalMsg.Sender,
				nftCustomMsg.RevokeApprovalMsg.DenomId,
				nftCustomMsg.RevokeApprovalMsg.TokenId,
				nftCustomMsg.RevokeApprovalMsg.ContractAddressSigner)
			return []sdk.Msg{approveNftMsg}, nil
		default:
			return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown custom nft message variant")
		}
	}
}

type nftCustomMsg struct {
	IssueDenomMsg     *IssueDenomRequest     `json:"issue_denom_msg,omitempty"`
	MintNftMsg        *MintNftRequest        `json:"mint_nft_msg,omitempty"`
	EditNftMsg        *EditNftRequest        `json:"edit_nft_msg,omitempty"`
	TransferNftMsg    *TransferNftRequest    `json:"transfer_nft_msg,omitempty"`
	BurnNftMsg        *BurnNftRequest        `json:"burn_nft_msg,omitempty"`
	ApproveNftMsg     *ApproveNftRequest     `json:"approve_nft_msg,omitempty"`
	ApproveAllMsg     *ApproveAllRequest     `json:"approve_all_msg,omitempty"`
	RevokeApprovalMsg *RevokeApprovalRequest `json:"revoke_approval_msg,omitempty"`
}

type IssueDenomRequest struct {
	Id                    string `json:"id"`
	Name                  string `json:"name"`
	Schema                string `json:"schema,omitempty"`
	Sender                string `json:"sender"`
	ContractAddressSigner string `json:"contract_address_signer"`
	Symbol                string `json:"symbol"`
}

type MintNftRequest struct {
	DenomId               string `json:"denom_id"`
	Name                  string `json:"name,omitempty"`
	URI                   string `json:"uri,omitempty"`
	Data                  string `json:"data,omitempty"`
	Sender                string `json:"sender"`
	Recipient             string `json:"recipient,omitempty"`
	ContractAddressSigner string `json:"contract_address_signer"`
}

type EditNftRequest struct {
	DenomId               string `json:"denom_id"`
	TokenId               string `json:"token_id"`
	Name                  string `json:"name,omitempty"`
	URI                   string `json:"uri,omitempty"`
	Data                  string `json:"data,omitempty"`
	Sender                string `json:"sender"`
	ContractAddressSigner string `json:"contract_address_signer"`
}

type TransferNftRequest struct {
	TokenId               string `json:"token_id"`
	DenomId               string `json:"denom_id"`
	From                  string `json:"from"`
	To                    string `json:"to"`
	Sender                string `json:"sender"`
	ContractAddressSigner string `json:"contract_address_signer"`
}

type BurnNftRequest struct {
	DenomId               string `json:"denom_id"`
	TokenId               string `json:"token_id"`
	Sender                string `json:"sender"`
	ContractAddressSigner string `json:"contract_address_signer"`
}

type ApproveNftRequest struct {
	TokenId               string `json:"token_id"`
	DenomId               string `json:"denom_id"`
	ApprovedAddress       string `json:"approved_address"`
	Sender                string `json:"sender"`
	ContractAddressSigner string `json:"contract_address_signer"`
}

type ApproveAllRequest struct {
	ApprovedOperator      string `json:"approved_operator"`
	Approved              bool   `json:"approved"`
	Sender                string `json:"sender"`
	ContractAddressSigner string `json:"contract_address_signer"`
}

type RevokeApprovalRequest struct {
	AddressToRevoke       string `json:"address_to_revoke"`
	DenomId               string `json:"denom_id"`
	TokenId               string `json:"token_id"`
	Sender                string `json:"sender"`
	ContractAddressSigner string `json:"contract_address_signer"`
}
