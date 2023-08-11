package custom_bindings

import (
	"encoding/json"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/CudoVentures/cudos-node/x/addressbook/keeper"
	"github.com/CudoVentures/cudos-node/x/addressbook/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
)

func PerformCustomAddressbookQuery(keeper keeper.Keeper) wasmkeeper.CustomQuerier {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		var custom adressbookCustomQuery
		err := json.Unmarshal(request, &custom)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
		}

		switch {
		case custom.QueryAllAddresses != nil:
			addresses, err := keeper.AddressAll(ctx, &types.QueryAllAddressRequest{Pagination: custom.QueryAllAddresses.Pagination})
			if err != nil {
				return nil, err
			}
			return json.Marshal(addresses)
		case custom.QueryAddress != nil:
			address, err := keeper.Address(ctx, &types.QueryGetAddressRequest{
				Creator: custom.QueryAddress.Creator,
				Network: custom.QueryAddress.Network,
				Label:   custom.QueryAddress.Label,
			})
			if err != nil {
				return nil, err
			}
			return json.Marshal(address)
		case custom.QueryParams != nil:
			params, err := keeper.Params(ctx, &types.QueryParamsRequest{})
			if err != nil {
				return nil, err
			}
			return json.Marshal(params)
		}
		return nil, sdkerrors.Wrap(wasmtypes.ErrInvalidMsg, "Unknown Custom query variant")
	}
}

type adressbookCustomQuery struct {
	QueryAllAddresses *QueryAllAddresses `json:"query_all_addresses,omitempty"`
	QueryAddress      *QueryAddress      `json:"query_address,omitempty"`
	QueryParams       *QueryParams       `json:"query_params,omitempty"`
}

type QueryAllAddresses struct {
	Pagination *query.PageRequest `json:"pagination,omitempty"`
}

type QueryAddress struct {
	Creator string `json:"creator"`
	Network string `json:"network"`
	Label   string `json:"label"`
}

type QueryParams struct {
}
