package v1_1_1

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CreateUpgradeHandler creates an SDK upgrade handler for v1

func UpdateWasmParams(ctx sdk.Context, gk wasm.Keeper) {
	params := gk.GetParams(ctx)
	params.CodeUploadAccess.Permission = wasmtypes.AccessTypeNobody
	gk.SetParams(ctx, params)
}
