package keeper

import (
	"context"
	"strconv"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) PublishAuction(goCtx context.Context, msg *types.MsgPublishAuction) (*types.MsgPublishAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	auctionType, err := msg.GetAuctionType()
	if err != nil {
		return nil, err
	}

	auction, err := types.NewAuction(msg.Creator, msg.DenomId, msg.TokenId, ctx.BlockTime().Add(msg.Duration), auctionType)
	if err != nil {
		return nil, err
	}

	auctionID, err := k.Keeper.PublishAuction(ctx, auction)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventPublishAuctionType,
			sdk.NewAttribute(types.AttributeAuctionID, strconv.FormatUint(auctionID, 10)),
			sdk.NewAttribute(types.AttributeKeyTokenID, msg.TokenId),
			sdk.NewAttribute(types.AttributeKeyDenomID, msg.DenomId),
			sdk.NewAttribute(types.AttributeAuctionType, auctionType.String()),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	})

	return &types.MsgPublishAuctionResponse{}, err
}
