package keeper

import (
	"context"
	"strconv"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdatePrice(goCtx context.Context, msg *types.MsgUpdatePrice) (*types.MsgUpdatePriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	nft, err := k.SetNftPrice(ctx, msg.Creator, msg.Id, msg.Price)
	if err != nil {
		return &types.MsgUpdatePriceResponse{}, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventUpdatePriceType,
			sdk.NewAttribute(types.AttributeKeyDenomID, nft.DenomId),
			sdk.NewAttribute(types.AttributeKeyTokenID, nft.TokenId),
			sdk.NewAttribute(types.AttributeKeyNftID, strconv.FormatUint(msg.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	})

	return &types.MsgUpdatePriceResponse{}, nil
}
