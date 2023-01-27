package keeper

import (
	"context"
	"fmt"
	"strconv"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) BuyNft(goCtx context.Context, msg *types.MsgBuyNft) (*types.MsgBuyNftResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	buyer, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	nftBefore, found := k.GetNft(ctx, msg.Id)
	if !found {
		return nil, fmt.Errorf("NFT not found. ID: {%d}", msg.Id)
	}

	nft, err := k.Keeper.BuyNFT(ctx, msg.Id, buyer)

	if err != nil {
		return &types.MsgBuyNftResponse{}, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventBuyNftType,
			sdk.NewAttribute(types.AttributeKeyDenomID, nft.DenomId),
			sdk.NewAttribute(types.AttributeKeyTokenID, nft.TokenId),
			sdk.NewAttribute(types.AttributeKeyNftID, strconv.FormatUint(msg.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyPrice, nftBefore.Price.String()),
			sdk.NewAttribute(types.AttributeKeyOwner, nftBefore.Owner),
			sdk.NewAttribute(types.AttributeKeyBuyer, msg.Creator),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	})

	return &types.MsgBuyNftResponse{}, nil
}
