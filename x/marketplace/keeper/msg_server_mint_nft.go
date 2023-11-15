package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
)

func (k msgServer) MintNft(goCtx context.Context, msg *types.MsgMintNft) (*types.MsgMintNftResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return &types.MsgMintNftResponse{}, err
	}

	recipient, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return &types.MsgMintNftResponse{}, err
	}

	tokenId, err := k.Keeper.MintNFT(ctx, msg.DenomId, msg.Name, msg.Uri, msg.Data, msg.Price, recipient, sender)
	if err != nil {
		return &types.MsgMintNftResponse{}, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventMintNftType,
			sdk.NewAttribute(types.AttributeKeyDenomID, msg.DenomId),
			sdk.NewAttribute(types.AttributeKeyTokenID, tokenId),
			sdk.NewAttribute(types.AttributeKeyBuyer, msg.Recipient),
			sdk.NewAttribute(types.AttributeKeyUID, msg.Uid),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	})

	return &types.MsgMintNftResponse{}, nil
}
