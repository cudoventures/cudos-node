package keeper

import (
	"context"

	"github.com/CudoVentures/cudos-node/x/addressbook/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateAddress(goCtx context.Context, msg *types.MsgCreateAddress) (*types.MsgCreateAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetAddress(
		ctx,
		msg.Creator,
		msg.Network,
		msg.Label,
	)
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	address := types.Address{
		Creator: msg.Creator,
		Network: msg.Network,
		Label:   msg.Label,
		Value:   msg.Value,
	}

	k.CreateNewAddress(ctx, address)

	return &types.MsgCreateAddressResponse{}, nil
}

func (k msgServer) UpdateAddress(goCtx context.Context, msg *types.MsgUpdateAddress) (*types.MsgUpdateAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetAddress(
		ctx,
		msg.Creator,
		msg.Network,
		msg.Label,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != valFound.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	address := types.Address{
		Creator: msg.Creator,
		Network: msg.Network,
		Label:   msg.Label,
		Value:   msg.Value,
	}

	k.UpdateExistingAddress(ctx, address)

	return &types.MsgUpdateAddressResponse{}, nil
}

func (k msgServer) DeleteAddress(goCtx context.Context, msg *types.MsgDeleteAddress) (*types.MsgDeleteAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetAddress(
		ctx,
		msg.Creator,
		msg.Network,
		msg.Label,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != valFound.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveAddress(ctx, msg.Creator, msg.Network, msg.Label)

	return &types.MsgDeleteAddressResponse{}, nil
}
