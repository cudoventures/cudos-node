package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
)

func (k msgServer) PublishNft(goCtx context.Context, msg *types.MsgPublishNft) (*types.MsgPublishNftResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return nil, err
	}

	nftID, err := k.PublishNFT(ctx, types.NewNft(msg.TokenId, msg.DenomId, msg.Creator, msg.Price))
	if err != nil {
		return &types.MsgPublishNftResponse{}, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventPublishNftType,
			sdk.NewAttribute(types.AttributeKeyNftID, strconv.FormatUint(nftID, 10)),
			sdk.NewAttribute(types.AttributeKeyTokenID, msg.TokenId),
			sdk.NewAttribute(types.AttributeKeyDenomID, msg.DenomId),
			sdk.NewAttribute(types.AttributeKeyPrice, msg.Price.String()),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	})

	return &types.MsgPublishNftResponse{}, nil
}
