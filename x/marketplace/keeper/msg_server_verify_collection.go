package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
)

func (k msgServer) VerifyCollection(goCtx context.Context, msg *types.MsgVerifyCollection) (*types.MsgVerifyCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.IsAdmin(ctx, msg.Creator); err != nil {
		return &types.MsgVerifyCollectionResponse{}, err
	}

	verified, err := k.GetCollectionStatus(ctx, msg.Id)
	if err != nil {
		return &types.MsgVerifyCollectionResponse{}, err
	}

	if verified == true {
		return &types.MsgVerifyCollectionResponse{}, sdkerrors.Wrapf(types.ErrCollectionAlreadyVerified, "collection %d is verified", msg.Id)
	}

	if err := k.SetCollectionStatus(ctx, msg.Id, true); err != nil {
		return &types.MsgVerifyCollectionResponse{}, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventVerifyCollectionType,
			sdk.NewAttribute(types.AttributeKeyCollectionID, strconv.FormatUint(msg.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	})

	return &types.MsgVerifyCollectionResponse{}, nil
}
