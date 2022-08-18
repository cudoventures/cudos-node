package keeper

import (
	"context"
	"strconv"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m msgServer) PublishCollection(goCtx context.Context, msg *types.MsgPublishCollection) (*types.MsgPublishCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	collectionID, err := m.Keeper.PublishCollection(ctx, types.NewCollection(msg.DenomId, msg.FirstSaleRoyalties, msg.ResaleRoyalties, creator.String(), false))
	if err != nil {
		return &types.MsgPublishCollectionResponse{}, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventPublishCollectionType,
			sdk.NewAttribute(types.AttributeKeyCollectionID, strconv.FormatUint(collectionID, 10)),
			sdk.NewAttribute(types.AttributeKeyDenomID, msg.DenomId),
			sdk.NewAttribute(types.AttributeKeyCreator, creator.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, creator.String()),
		),
	})

	return &types.MsgPublishCollectionResponse{}, nil
}
