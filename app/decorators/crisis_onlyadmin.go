package decorators

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
)

// Copied from https://github.com/CudoVentures/cosmos-sdk/blob/3816012a2d4ea5c9bbb3d8e6174d3b96ff91a039/x/crisis/keeper/msg_server.go#L11
var (
	ErrAdminOnly    = errors.New("sender has no admin tokens")
	adminTokenDenom = "cudosAdmin"
)

type OnlyAdminVerifyInvariantDecorator struct {
	bankKeeper bankkeeper.ViewKeeper
}

func NewOnlyAdminVerifyInvariantDecorator(bankKeeper bankkeeper.ViewKeeper) *OnlyAdminVerifyInvariantDecorator {
	return &OnlyAdminVerifyInvariantDecorator{
		bankKeeper: bankKeeper,
	}
}

func (od *OnlyAdminVerifyInvariantDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool, next sdk.AnteHandler,
) (newCtx sdk.Context, err error) {
	if err := od.checkMsgs(ctx, tx.GetMsgs()); err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate)
}

func (od *OnlyAdminVerifyInvariantDecorator) checkMsgs(ctx sdk.Context, msgs []sdk.Msg) error {
	for _, msg := range msgs {
		if msg, ok := msg.(*crisistypes.MsgVerifyInvariant); ok {
			if err := od.checkVerifyInvariant(ctx, msg); err != nil {
				return err
			}
			continue
		}
	}

	return nil
}

// Copied from https://github.com/CudoVentures/cosmos-sdk/blob/3816012a2d4ea5c9bbb3d8e6174d3b96ff91a039/x/crisis/keeper/msg_server.go#L24
func (od *OnlyAdminVerifyInvariantDecorator) checkVerifyInvariant(ctx sdk.Context, msg *crisistypes.MsgVerifyInvariant) error {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return err
	}

	adminCoin := od.bankKeeper.GetBalance(ctx, sender, adminTokenDenom)

	if adminCoin.Amount.Equal(sdk.ZeroInt()) {
		return ErrAdminOnly
	}
	return nil
}
