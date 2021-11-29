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
		case nftCustomMsg.MintNft != nil:
			mintNftMsg := nftTypes.MsgMintNFT{
				DenomId:   nftCustomMsg.MintNft.DenomId,
				Name:      nftCustomMsg.MintNft.Name,
				URI:       nftCustomMsg.MintNft.URI,
				Data:      nftCustomMsg.MintNft.Data,
				Sender:    nftCustomMsg.MintNft.Sender,
				Recipient: nftCustomMsg.MintNft.Recipient,
			}
			return []sdk.Msg{&mintNftMsg}, nil
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
		// TODO: Discuss with Megi about the return result for a query
		// - should we a shared infrastructure or define new response specifacally for custom runtime call ?
		// like the below two examples
		case custom.QueryDenomById != nil:
			denom, err := keeper.GetDenom(ctx, custom.QueryDenomById.DenomId)
			if err != nil {
				return nil, err
			}
			result, err := json.Marshal(nftTypes.QueryDenomResponse{Denom: &denom})
			return result, err
		case custom.QueryDenomByIdTest != nil:
			denom, err := keeper.GetDenom(ctx, custom.QueryDenomByIdTest.DenomId)
			if err != nil {
				return nil, err
			}
			result, err := json.Marshal(QueryDenomResponseTest{
				Id:      denom.Id,
				Name:    denom.Name,
				Schema:  denom.Schema,
				Creator: denom.Creator,
			})
			return result, err
		}

		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown Custom query variant")
	}
}

type nftCustomMsg struct {
	IssueDenom *IssueDenomRequest `json:"issue_denom,omitempty"`
	MintNft    *MintNft           `json:"mint_nft,omitempty"`
}

type nftCustomQuery struct {
	QueryDenomById     *QueryDenomById `json:"query_denom_by_id,omitempty"`
	QueryDenomByIdTest *QueryDenomById `json:"query_denom_by_id_test,omitempty"`
}

type IssueDenomRequest struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Schema string `json:"schema,omitempty"`
	Sender string `json:"sender"`
}

type QueryDenomResponseTest struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Schema  string `json:"schema"`
	Creator string `json:"creator"`
}

type MintNft struct {
	DenomId   string `json:"denomId"`
	Name      string `json:"Name"`
	URI       string `json:"uri"`
	Data      string `json:"data"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
}

type QueryDenomById struct {
	DenomId string `json:"denom_id"`
}
