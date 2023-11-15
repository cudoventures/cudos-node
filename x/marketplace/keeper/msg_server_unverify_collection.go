package keeper

import (
	"context"
	"strconv"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) UnverifyCollection(goCtx context.Context, msg *types.MsgUnverifyCollection) (*types.MsgUnverifyCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.IsAdmin(ctx, msg.Creator); err != nil {
		return &types.MsgUnverifyCollectionResponse{}, err
	}

	verified, err := k.GetCollectionStatus(ctx, msg.Id)
	if err != nil {
		return &types.MsgUnverifyCollectionResponse{}, err
	}

	if verified == false {
		return &types.MsgUnverifyCollectionResponse{}, sdkerrors.Wrapf(types.ErrCollectionAlreadyUnverified, "collection %d is not verified", msg.Id)
	}

	if err := k.SetCollectionStatus(ctx, msg.Id, false); err != nil {
		return &types.MsgUnverifyCollectionResponse{}, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventUnverifyCollectionType,
			sdk.NewAttribute(types.AttributeKeyCollectionID, strconv.FormatUint(msg.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	})

	return &types.MsgUnverifyCollectionResponse{}, nil
}
