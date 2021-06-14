package types

import (
	"encoding/json"
	"fmt"
	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/CosmWasm/wasmvm/types"
	"github.com/cosmos/cosmos-sdk/codec"
	types2 "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"math/rand"
)

type Text struct {
	Text string `json:"text"`
}

/**** Code to support custom messages *****/

type reflectCustomMsg struct {
	Debug string `json:"debug,omitempty"`
	Raw   []byte `json:"raw,omitempty"`
}

// toReflectRawMsg encodes an sdk msg using any type with json encoding.
// Then wraps it as an opaque message
func toReflectRawMsg(cdc codec.Marshaler, msg sdk.Msg) (types.CosmosMsg, error) {
	any, err := types2.NewAnyWithValue(msg)
	if err != nil {
		return types.CosmosMsg{}, err
	}
	rawBz, err := cdc.MarshalJSON(any)
	if err != nil {
		return types.CosmosMsg{}, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	customMsg, err := json.Marshal(reflectCustomMsg{
		Raw: rawBz,
	})
	res := types.CosmosMsg{
		Custom: customMsg,
	}
	return res, nil
}

// reflectEncoders needs to be registered in test setup to handle custom message callbacks
func ReflectEncoders(cdc codec.Marshaler) *wasm.MessageEncoders {
	return &wasm.MessageEncoders{
		Custom: fromReflectRawMsg(cdc),
	}
}

// fromReflectRawMsg decodes msg.Data to an sdk.Msg using proto Any and json encoding.
// this needs to be registered on the Encoders
func fromReflectRawMsg(cdc codec.Marshaler) wasm.CustomEncoder {
	return func(_sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, error) {
		var custom reflectCustomMsg
		err := json.Unmarshal(msg, &custom)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
		}
		if custom.Raw != nil {
			var any types2.Any
			if err := cdc.UnmarshalJSON(custom.Raw, &any); err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
			}
			var msg sdk.Msg
			if err := cdc.UnpackAny(&any, &msg); err != nil {
				return nil, err
			}
			return []sdk.Msg{msg}, nil
		}
		if custom.Debug != "" {
			return nil, sdkerrors.Wrapf(wasm.ErrInvalidMsg, "Custom Debug: %s", custom.Debug)
		}
		return nil, sdkerrors.Wrap(wasm.ErrInvalidMsg, "Unknown Custom message variant")
	}
}

type reflectCustomQueryWrapper struct {
	Route string `json:"route,omitempty""`
	QueryData reflectCustomQuery `json:"query_data,omitempty"`
}

type reflectCustomQuery struct {
	Ping        *struct{} `json:"ping,omitempty"`
	Capitalized *Text     `json:"capital,omitempty"`
}

// this is from the go code back to the contract (capitalized or ping)
type customQueryResponse struct {
	Msg uint64 `json:"msg"`
}

// these are the return values from contract -> go depending on type of query
type ownerResponse struct {
	Owner string `json:"owner"`
}

type capitalizedResponse struct {
	Text string `json:"text"`
}

type chainResponse struct {
	Data []byte `json:"data"`
}



// reflectPlugins needs to be registered in test setup to handle custom query callbacks
func ReflectPlugins() *wasm.QueryPlugins {
	return &wasm.QueryPlugins{
		Custom: performCustomQuery,
	}
}

func performCustomQuery(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
	var custom reflectCustomQueryWrapper
	err := json.Unmarshal(request, &custom)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	randMsg := rand.Uint64() % 10_000

	fmt.Printf("Here wi %s", request)
	ctx.Logger().Info(fmt.Sprintf("here with %v %v", custom.QueryData.Capitalized, custom.QueryData.Ping))
	if custom.QueryData.Capitalized != nil {
		//msg := strings.ToUpper(custom.QueryData.Capitalized.Text)
		return json.Marshal(customQueryResponse{Msg: randMsg})
	}
	if custom.QueryData.Capitalized != nil {
		return json.Marshal(customQueryResponse{Msg: randMsg})
	}
	return nil, sdkerrors.Wrap(wasm.ErrInvalidMsg, "Unknown Custom query variant")
}
