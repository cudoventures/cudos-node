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
		case custom.QueryDenom != nil:
			msg := nftTypes.QueryDenomRequest{
				DenomId: custom.QueryDenom.DenomId,
			}
			denom, err := keeper.GetDenom(ctx, msg.DenomId)
			if err != nil {
				return nil, err
			}
			return json.Marshal(nftTypes.QueryDenomResponse{Denom: &denom})
		}
		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown Custom query variant")
	}
}

type nftCustomMsg struct {
	IssueDenom *Denom   `json:"issue_denom,omitempty"`
	MintNft    *MintNft `json:"mint_nft,omitempty"`
}

type nftCustomQuery struct {
	Ping        *struct{}   `json:"ping,omitempty"`
	Capitalized *Text       `json:"capitalized,omitempty"`
	QueryDenom  *QueryDenom `json:"queryDenom,omitempty"`
}

type Text struct {
	Text string `json:"text"`
}

// this is from the go code back to the contract (capitalized or ping)
type customQueryResponse struct {
	Msg string `json:"msg"`
}

type Denom struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Schema string `json:"schema"`
	Sender string `json:"sender"`
}

type MintNft struct {
	DenomId   string `json:"denomId"`
	Name      string `json:"Name"`
	URI       string `json:"uri"`
	Data      string `json:"data"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
}

type QueryDenom struct {
	DenomId string `json:"denomId"`
}
