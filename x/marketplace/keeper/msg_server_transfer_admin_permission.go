package keeper

import (
	"context"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) TransferAdminPermission(goCtx context.Context, msg *types.MsgTransferAdminPermission) (*types.MsgTransferAdminPermissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.IsAdmin(ctx, msg.Creator); err != nil {
		return &types.MsgTransferAdminPermissionResponse{}, err
	}

	if err := k.transferAdminPermission(ctx, msg.Creator, msg.NewAdmin); err != nil {
		return &types.MsgTransferAdminPermissionResponse{}, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTransferAdminPermissionType,
			sdk.NewAttribute(types.AttributeKeyNewAdmin, msg.NewAdmin),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	})

	return &types.MsgTransferAdminPermissionResponse{}, nil
}
