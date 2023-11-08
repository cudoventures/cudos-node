package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
)

func (m msgServer) PublishCollection(goCtx context.Context, msg *types.MsgPublishCollection) (*types.MsgPublishCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	collectionID, err := m.Keeper.PublishCollection(ctx, types.NewCollection(msg.DenomId, msg.MintRoyalties, msg.ResaleRoyalties, msg.Creator, false))
	if err != nil {
		return &types.MsgPublishCollectionResponse{}, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventPublishCollectionType,
			sdk.NewAttribute(types.AttributeKeyCollectionID, strconv.FormatUint(collectionID, 10)),
			sdk.NewAttribute(types.AttributeKeyDenomID, msg.DenomId),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	})

	return &types.MsgPublishCollectionResponse{}, nil
}
