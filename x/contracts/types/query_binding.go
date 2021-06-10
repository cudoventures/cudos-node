package types

import (
	"encoding/json"
	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/CosmWasm/wasmvm/types"
	codec2 "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strings"
)

type maskCustomMsg struct {
	Debug string `json:"debug,omitempty"`
	Raw   []byte `json:"raw,omitempty"`
}

// toMaskRawMsg encodes an sdk msg using amino json encoding.
// Then wraps it as an opaque message
func toMaskRawMsg(cdc *codec2.LegacyAmino, msg sdk.Msg) (types.CosmosMsg, error) {
	rawBz, err := cdc.MarshalJSON(msg)
	if err != nil {
		return types.CosmosMsg{}, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	customMsg, err := json.Marshal(maskCustomMsg{
		Raw: rawBz,
	})
	res := types.CosmosMsg{
		Custom: customMsg,
	}
	return res, nil
}

// maskEncoders needs to be registered in test setup to handle custom message callbacks
func maskEncoders(cdc *codec2.LegacyAmino) *wasm.MessageEncoders {
	return &wasm.MessageEncoders{
		Custom: fromMaskRawMsg(cdc),
	}
}

// fromMaskRawMsg decodes msg.Data to an sdk.Msg using amino json encoding.
// this needs to be registered on the Encoders
func fromMaskRawMsg(cdc *codec2.LegacyAmino) wasm.CustomEncoder {
	return func(_sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, error) {
		var custom maskCustomMsg
		err := json.Unmarshal(msg, &custom)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
		}
		if custom.Raw != nil {
			var sdkMsg sdk.Msg
			err := cdc.UnmarshalJSON(custom.Raw, &sdkMsg)
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
			}
			return []sdk.Msg{sdkMsg}, nil
		}
		if custom.Debug != "" {
			return nil, sdkerrors.Wrapf(wasm.ErrInvalidMsg, "Custom Debug: %s", custom.Debug)
		}
		return nil, sdkerrors.Wrap(wasm.ErrInvalidMsg, "Unknown Custom message variant")
	}
}



type Text struct {
	Text string `json:"text"`
}

type maskCustomQuery struct {
	Ping    *struct{} `json:"ping,omitempty"`
	Capital *Text     `json:"capital,omitempty"`
}

type customQueryResponse struct {
	Msg string `json:"msg"`
}

// maskPlugins needs to be registered in test setup to handle custom query callbacks
func MaskPlugins() *wasm.QueryPlugins {
	return &wasm.QueryPlugins{
		Custom: performCustomQuery,
	}
}

func performCustomQuery(_ sdk.Context, request json.RawMessage) ([]byte, error) {
	var custom maskCustomQuery
	err := json.Unmarshal(request, &custom)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	if custom.Capital != nil {
		msg := strings.ToUpper(custom.Capital.Text)
		return json.Marshal(customQueryResponse{Msg: msg})
	}
	if custom.Ping != nil {
		return json.Marshal(customQueryResponse{Msg: "pong"})
	}
	return nil, sdkerrors.Wrap(wasm.ErrInvalidMsg, "Unknown Custom query variant")
}
