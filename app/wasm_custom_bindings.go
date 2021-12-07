package app

import (
	nftKeeper "cudos.org/cudos-node/x/nft/keeper"
	nftTypes "cudos.org/cudos-node/x/nft/types"
	"encoding/json"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmKeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmd/x/wasm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func GetCustomMsgEncodersOptions() []wasmKeeper.Option {
	nftEncodingOptions := wasmKeeper.WithMessageEncoders(nftEncoders())
	return []wasm.Option{nftEncodingOptions}
}

func GetCustomMsgQueryOptions(keeper nftKeeper.Keeper) []wasmKeeper.Option {
	nftQueryOptions := wasmKeeper.WithQueryPlugins(nftQueryPlugins(keeper))
	return []wasm.Option{nftQueryOptions}
}

func nftEncoders() *wasmKeeper.MessageEncoders {
	return &wasmKeeper.MessageEncoders{
		Custom: encodeNftMessage(),
	}
}

// nftQueryPlugins needs to be registered in test setup to handle custom query callbacks
func nftQueryPlugins(keeper nftKeeper.Keeper) *wasmKeeper.QueryPlugins {
	return &wasmKeeper.QueryPlugins{
		Custom: performCustomNftQuery(keeper),
	}
}

// fromReflectRawMsg decodes msg.Data to an sdk.Msg using proto Any and json encoding.
// this needs to be registered on the Encoders
func encodeNftMessage() wasmKeeper.CustomEncoder {
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
				nftCustomMsg.IssueDenomMsg.ContractAddressSigner)
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

func performCustomNftQuery(keeper nftKeeper.Keeper) wasmKeeper.CustomQuerier {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		var custom nftCustomQuery
		err := json.Unmarshal(request, &custom)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
		}

		switch {
		case custom.QueryDenomById != nil:
			denom, err := keeper.GetDenom(ctx, custom.QueryDenomById.DenomId)
			if err != nil {
				return nil, err
			}
			return json.Marshal(nftTypes.QueryDenomResponse{Denom: &denom})
		case custom.QueryDenomByName != nil:
			denom, err := keeper.GetDenomByName(ctx, custom.QueryDenomByName.DenomName)
			if err != nil {
				return nil, err
			}
			return json.Marshal(nftTypes.QueryDenomResponse{Denom: &denom})
		case custom.QueryDenoms != nil:
			denoms := keeper.GetDenoms(ctx)
			return json.Marshal(nftTypes.QueryDenomsResponse{Denoms: denoms}.Pagination.NextKey)
		case custom.QueryCollection != nil:
			collection, err := keeper.GetCollection(ctx, custom.QueryCollection.DenomId)
			if err != nil {
				return nil, err
			}
			return json.Marshal(nftTypes.QueryCollectionResponse{Collection: &collection})
		case custom.QuerySupply != nil:
			totalSupply := keeper.GetTotalSupply(ctx, custom.QueryCollection.DenomId)
			return json.Marshal(nftTypes.QuerySupplyResponse{Amount: totalSupply})
		case custom.QueryOwner != nil:
			if len(custom.QueryOwner.Address) > 0 {
				ownerAddress, err := sdk.AccAddressFromBech32(custom.QueryOwner.Address)
				if err != nil {
					return nil, err
				}
				owner, err := keeper.GetOwner(ctx, ownerAddress, custom.QueryOwner.DenomId)
				if err != nil {
					return nil, err
				}
				return json.Marshal(nftTypes.QueryOwnerResponse{Owner: &owner})
			}
			return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Owner address is empty!")
		case custom.QueryToken != nil:
			nft, err := keeper.GetBaseNFT(ctx, custom.QueryToken.DenomId, custom.QueryToken.TokenId)
			if err != nil {
				return nil, err
			}
			return json.Marshal(nftTypes.QueryNFTResponse{NFT: &nft})
		case custom.QueryApprovals != nil:
			approvedAddressesForNft, err := keeper.GetNFTApprovedAddresses(ctx, custom.QueryApprovals.DenomId, custom.QueryApprovals.TokenId)
			if err != nil {
				return nil, err
			}
			return json.Marshal(nftTypes.QueryApprovalsNFTResponse{ApprovedAddresses: approvedAddressesForNft})
		case custom.QueryApprovedForAll != nil:
			if len(custom.QueryApprovedForAll.OwnerAddress) > 0 && len(custom.QueryApprovedForAll.OperatorAddress) > 0 {
				ownerAddress, err := sdk.AccAddressFromBech32(custom.QueryApprovedForAll.OwnerAddress)
				if err != nil {
					return nil, err
				}

				operatorAddress, err := sdk.AccAddressFromBech32(custom.QueryApprovedForAll.OperatorAddress)
				if err != nil {
					return nil, err
				}

				isApproved := keeper.IsApprovedOperator(ctx, ownerAddress, operatorAddress)
				return json.Marshal(nftTypes.QueryApprovalsIsApprovedForAllResponse{IsApproved: isApproved})
			}

		}
		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown Custom query variant")
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

type nftCustomQuery struct {
	QueryDenomById      *QueryDenomById      `json:"query_denom_by_id,omitempty"`
	QueryDenomByName    *QueryDenomByName    `json:"query_denom_by_name,omitempty"`
	QueryDenoms         *QueryAllDenoms      `json:"query_denoms,omitempty"`
	QueryCollection     *QueryCollection     `json:"query_collection,omitempty"`
	QuerySupply         *QuerySupply         `json:"query_supply,omitempty"`
	QueryOwner          *QueryOwner          `json:"query_owner,omitempty"`
	QueryToken          *QueryToken          `json:"query_token,omitempty"`
	QueryApprovals      *QueryApprovals      `json:"query_approvals,omitempty"`
	QueryApprovedForAll *QueryApprovedForAll `json:"query_approved_for_all,omitempty"`
}

type IssueDenomRequest struct {
	Id                    string `json:"id"`
	Name                  string `json:"name"`
	Schema                string `json:"schema"`
	Sender                string `json:"sender"`
	ContractAddressSigner string `json:"contract_address_signer"`
}

type MintNftRequest struct {
	DenomId               string `json:"denom_id"`
	Name                  string `json:"name"`
	URI                   string `json:"uri"`
	Data                  string `json:"data"`
	Sender                string `json:"sender"`
	Recipient             string `json:"recipient"`
	ContractAddressSigner string `json:"contract_address_signer"`
}

type EditNftRequest struct {
	DenomId               string `json:"denom_id"`
	TokenId               string `json:"token_id"`
	Name                  string `json:"name"`
	URI                   string `json:"uri"`
	Data                  string `json:"data"`
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

type QueryDenomById struct {
	DenomId string `json:"denom_id"`
}

type QueryDenomByName struct {
	DenomName string `json:"denom_name"`
}

type QueryAllDenoms struct {
}

type QueryCollection struct {
	DenomId string `json:"denom_id"`
}

type QuerySupply struct {
	DenomId string `json:"denom_id"`
}

type QueryOwner struct {
	Address string `json:"address"`
	DenomId string `json:"denom_id"`
}

type QueryToken struct {
	DenomId string `json:"denom_id"`
	TokenId string `json:"token_id"`
}

type QueryApprovals struct {
	DenomId string `json:"denom_id"`
	TokenId string `json:"token_id"`
}

type QueryApprovedForAll struct {
	OwnerAddress    string `json:"owner_address"`
	OperatorAddress string `json:"operator_address"`
}
