package keeper

import (
	"context"
	"strconv"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateCollection(goCtx context.Context, msg *types.MsgCreateCollection) (*types.MsgCreateCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.Verified {
		if err := k.IsAdmin(ctx, msg.Creator); err != nil {
			return &types.MsgCreateCollectionResponse{}, err
		}
	}

	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	collectionID, err := k.Keeper.CreateCollection(ctx, creator, msg.Id, msg.Name, msg.Schema, msg.Symbol, msg.Traits, msg.Description,
		msg.Minter, msg.Data, msg.MintRoyalties, msg.ResaleRoyalties, msg.Verified)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventCreateCollectionType,
			sdk.NewAttribute(types.AttributeKeyCollectionID, strconv.FormatUint(collectionID, 10)),
			sdk.NewAttribute(types.AttributeKeyDenomID, msg.Id),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	})

	return &types.MsgCreateCollectionResponse{}, nil
}
