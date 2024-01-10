package v1_1_1

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func UpdateWasmParams(ctx sdk.Context, gk wasm.Keeper) {
	params := gk.GetParams(ctx)
	params.CodeUploadAccess.Permission = wasmtypes.AccessTypeNobody
	params.InstantiateDefaultPermission = wasmtypes.AccessTypeNobody
	gk.SetParams(ctx, params)
}
