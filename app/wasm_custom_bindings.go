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
		case nftCustomMsg.IssueDenom != nil:
			issueDenomMsg := nftTypes.MsgIssueDenom{
				Id:     nftCustomMsg.IssueDenom.Id,
				Name:   nftCustomMsg.IssueDenom.Name,
				Schema: nftCustomMsg.IssueDenom.Schema,
				Sender: nftCustomMsg.IssueDenom.Sender,
			}
			return []sdk.Msg{&issueDenomMsg}, nil
		case nftCustomMsg.MintNFT != nil:
			mintNftMsg := nftTypes.MsgMintNFT{
				DenomId:   nftCustomMsg.MintNFT.DenomId,
				Name:      nftCustomMsg.MintNFT.Name,
				URI:       nftCustomMsg.MintNFT.URI,
				Data:      nftCustomMsg.MintNFT.Data,
				Sender:    nftCustomMsg.MintNFT.Sender,
				Recipient: nftCustomMsg.MintNFT.Recipient,
			}
			return []sdk.Msg{&mintNftMsg}, nil
		case nftCustomMsg.EditNFT != nil:
			editNftMsg := nftTypes.MsgEditNFT{
				DenomId: nftCustomMsg.EditNFT.DenomId,
				Name:    nftCustomMsg.EditNFT.Name,
				URI:     nftCustomMsg.EditNFT.URI,
				Data:    nftCustomMsg.EditNFT.Data,
				Sender:  nftCustomMsg.EditNFT.Sender,
			}
			return []sdk.Msg{&editNftMsg}, nil
		case nftCustomMsg.TransferNFT != nil:
			transferNftMsg := nftTypes.MsgTransferNft{
				TokenId: nftCustomMsg.TransferNFT.TokenId,
				DenomId: nftCustomMsg.TransferNFT.DenomId,
				From:    nftCustomMsg.TransferNFT.From,
				To:      nftCustomMsg.TransferNFT.To,
				Sender:  nftCustomMsg.TransferNFT.Sender,
			}
			return []sdk.Msg{&transferNftMsg}, nil
		case nftCustomMsg.BurnNft != nil:
			burnNftMsg := nftTypes.MsgBurnNFT{
				Id:      nftCustomMsg.BurnNft.TokenId,
				DenomId: nftCustomMsg.BurnNft.DenomId,
				Sender:  nftCustomMsg.BurnNft.Sender,
			}
			return []sdk.Msg{&burnNftMsg}, nil
		case nftCustomMsg.ApproveNft != nil:
			approveNftMsg := nftTypes.MsgApproveNft{
				Id:              nftCustomMsg.ApproveNft.TokenId,
				DenomId:         nftCustomMsg.ApproveNft.DenomId,
				Sender:          nftCustomMsg.ApproveNft.Sender,
				ApprovedAddress: nftCustomMsg.ApproveNft.ApprovedAddress,
			}
			return []sdk.Msg{&approveNftMsg}, nil
		case nftCustomMsg.ApproveAll != nil:
			approveNftMsg := nftTypes.MsgApproveAllNft{
				Operator: nftCustomMsg.ApproveAll.ApprovedOperator,
				Sender:   nftCustomMsg.ApproveAll.Sender,
				Approved: nftCustomMsg.ApproveAll.Approved,
			}
			return []sdk.Msg{&approveNftMsg}, nil
		case nftCustomMsg.RevokeApproval != nil:
			approveNftMsg := nftTypes.MsgRevokeNft{
				AddressToRevoke: nftCustomMsg.RevokeApproval.AddressToRevoke,
				DenomId:         nftCustomMsg.RevokeApproval.DenomId,
				TokenId:         nftCustomMsg.RevokeApproval.TokenId,
				Sender:          nftCustomMsg.RevokeApproval.Sender,
			}
			return []sdk.Msg{&approveNftMsg}, nil
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
			return json.Marshal(nftTypes.QueryDenomsResponse{Denoms: denoms})
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
	IssueDenom     *IssueDenomRequest     `json:"issue_denom,omitempty"`
	MintNFT        *MintNftRequest        `json:"mint_nft,omitempty"`
	EditNFT        *EditNftRequest        `json:"edit_nft,omitempty"`
	TransferNFT    *TransferNftRequest    `json:"transfer_nft,omitempty"`
	BurnNft        *BurnNftRequest        `json:"burn_nft,omitempty"`
	ApproveNft     *ApproveNftRequest     `json:"approve_nft,omitempty"`
	ApproveAll     *ApproveAllRequest     `json:"approve_all,omitempty"`
	RevokeApproval *RevokeApprovalRequest `json:"revoke_approval,omitempty"`
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
	Id     string `json:"id"`
	Name   string `json:"name"`
	Schema string `json:"schema,omitempty"`
	Sender string `json:"sender"`
}

type MintNftRequest struct {
	DenomId   string `json:"denomId"`
	Name      string `json:"Name"`
	URI       string `json:"uri"`
	Data      string `json:"data"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
}

type EditNftRequest struct {
	DenomId string `json:"denomId"`
	Name    string `json:"Name"`
	URI     string `json:"uri"`
	Data    string `json:"data"`
	Sender  string `json:"sender"`
}

type TransferNftRequest struct {
	TokenId string `json:"token_id"`
	DenomId string `json:"denom_id"`
	From    string `json:"from"`
	To      string `json:"to"`
	Sender  string `json:"sender"`
}

type BurnNftRequest struct {
	DenomId string `json:"denom_id"`
	TokenId string `json:"token_id"`
	Sender  string `json:"sender"`
}

type ApproveNftRequest struct {
	TokenId         string `json:"token_id"`
	DenomId         string `json:"denom_id"`
	ApprovedAddress string `json:"approved_address"`
	Sender          string `json:"sender"`
}

type ApproveAllRequest struct {
	ApprovedOperator string `json:"approved_operator"`
	Approved         bool   `json:"approved"`
	Sender           string `json:"sender"`
}

type RevokeApprovalRequest struct {
	AddressToRevoke string `json:"address_to_revoke"`
	DenomId         string `json:"denom_id"`
	TokenId         string `json:"token_id"`
	Sender          string `json:"sender"`
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
