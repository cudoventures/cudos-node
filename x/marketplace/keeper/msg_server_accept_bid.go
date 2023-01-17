package keeper

import (
	"context"
	"strconv"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) AcceptBid(
	goCtx context.Context, msg *types.MsgAcceptBid,
) (*types.MsgAcceptBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.Keeper.AcceptBid(ctx, msg); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventAcceptBidType,
			sdk.NewAttribute(types.AttributeAuctionID, strconv.FormatUint(msg.AuctionId, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgAcceptBidResponse{}, nil
}
