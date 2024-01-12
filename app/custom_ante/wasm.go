package custom_ante

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// TODO: Replace with the correct height
const DenyMsgStoreCodeAfter = 1000000

func DenyMsgStoreCode(next sdk.AnteHandler) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (_ sdk.Context, _ error) {
		if ctx.BlockHeight() > DenyMsgStoreCodeAfter {
			for _, msg := range tx.GetMsgs() {
				if _, ok := msg.(*wasmtypes.MsgStoreCode); ok {
					return ctx, sdkerrors.Wrap(
						sdkerrors.ErrInvalidRequest,
						"MsgStoreCode is disabled",
					)
				}
			}
		}
		return next(ctx, tx, simulate)
	}
}
