package keeper

import (
	"context"
	"strconv"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateRoyalties(goCtx context.Context, msg *types.MsgUpdateRoyalties) (*types.MsgUpdateRoyaltiesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.SetCollectionRoyalties(ctx, msg.Creator, msg.Id, msg.MintRoyalties, msg.ResaleRoyalties); err != nil {
		return &types.MsgUpdateRoyaltiesResponse{}, nil
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventUpdateRoyaltiesType,
			sdk.NewAttribute(types.AttributeKeyCollectionID, strconv.FormatUint(msg.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	})

	return &types.MsgUpdateRoyaltiesResponse{}, nil
}
