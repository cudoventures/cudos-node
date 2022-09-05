package keeper

import (
	"context"
	"strconv"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) RemoveNft(goCtx context.Context, msg *types.MsgRemoveNft) (*types.MsgRemoveNftResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	if err := k.RemoveNFT(ctx, msg.Id, owner); err != nil {
		return &types.MsgRemoveNftResponse{}, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventRemoveNftType,
			sdk.NewAttribute(types.AttributeKeyNftID, strconv.FormatUint(msg.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Creator),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	})

	return &types.MsgRemoveNftResponse{}, nil
}
