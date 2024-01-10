package decorators

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Copied from https://github.com/CudoVentures/cosmos-sdk/blob/3816012a2d4ea5c9bbb3d8e6174d3b96ff91a039/x/staking/types/msg.go#L20
const MinSelfDelegation = "2000000000000000000000000"

type MinSelfDelegationDecorator struct{}

func NewMinSelfDelegationDecorator() *MinSelfDelegationDecorator {
	return &MinSelfDelegationDecorator{}
}

func (m *MinSelfDelegationDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool, next sdk.AnteHandler,
) (newCtx sdk.Context, err error) {
	if err := m.checkMsgs(ctx, tx.GetMsgs()); err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate)
}

func (m *MinSelfDelegationDecorator) checkMsgs(ctx sdk.Context, msgs []sdk.Msg) error {
	for _, msg := range msgs {
		if msg, ok := msg.(*stakingtypes.MsgCreateValidator); ok {
			if err := m.checkCreateValidator(ctx, msg); err != nil {
				return err
			}
			continue
		}
	}

	return nil
}

func (m *MinSelfDelegationDecorator) checkCreateValidator(ctx sdk.Context, msg *stakingtypes.MsgCreateValidator) error {
	msd, _ := sdk.NewIntFromString(MinSelfDelegation)

	if msg.MinSelfDelegation.LT(msd) {
		// Copied from https://github.com/CudoVentures/cosmos-sdk/blob/3816012a2d4ea5c9bbb3d8e6174d3b96ff91a039/x/staking/types/msg.go#L143
		return sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("minimum self delegation must be more than %v", msd),
		)
	}
	return nil
}
