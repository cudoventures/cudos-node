package keeper

import (
	"context"
	"strconv"
	"time"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) PublishAuction(
	goCtx context.Context, msg *types.MsgPublishAuction,
) (*types.MsgPublishAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	a, err := types.AuctionFromMsgPublishAuction(msg, ctx.BlockTime())
	if err != nil {
		return nil, err
	}

	auctionId, err := k.Keeper.PublishAuction(ctx, a)
	if err != nil {
		return nil, err
	}

	auctionInfo, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventPublishAuctionType,
			sdk.NewAttribute(types.AttributeAuctionID, strconv.FormatUint(auctionId, 10)),
			sdk.NewAttribute(types.AttributeKeyTokenID, msg.TokenId),
			sdk.NewAttribute(types.AttributeKeyDenomID, msg.DenomId),
			sdk.NewAttribute(types.AttributeAuctionInfo, string(auctionInfo)),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeStartTime, ctx.BlockTime().Format(time.RFC3339)),
			sdk.NewAttribute(types.AttributeEndTime, a.GetEndTime().Format(time.RFC3339)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	})

	return &types.MsgPublishAuctionResponse{}, nil
}
