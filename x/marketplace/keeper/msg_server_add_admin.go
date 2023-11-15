package keeper

import (
	"context"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) AddAdmin(goCtx context.Context, msg *types.MsgAddAdmin) (*types.MsgAddAdminResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.Keeper.AddAdmin(ctx, msg.Address, msg.Creator); err != nil {
		return &types.MsgAddAdminResponse{}, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventAddAdminType,
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Address),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	})

	return &types.MsgAddAdminResponse{}, nil
}
