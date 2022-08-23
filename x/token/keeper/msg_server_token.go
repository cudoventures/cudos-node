package keeper

import (
	"context"

	"github.com/CudoVentures/cudos-node/x/token/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateToken(goCtx context.Context, msg *types.MsgCreateToken) (*types.MsgCreateTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetToken(
		ctx,
		msg.Denom,
	)
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var token = types.Token{
		Owner:    msg.Owner,
		Denom:    msg.Denom,
		Name:     msg.Name,
		Decimals: msg.Decimals,
		// todo handle initial balances (mint)
		// InitialBalances: msg.InitialBalances,
		MaxSupply:  msg.MaxSupply,
		Allowances: []*types.Allowances{},
	}

	k.SetToken(
		ctx,
		token,
	)
	return &types.MsgCreateTokenResponse{}, nil
}

func (k msgServer) UpdateToken(goCtx context.Context, msg *types.MsgUpdateToken) (*types.MsgUpdateTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetToken(
		ctx,
		msg.Denom,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg owner is the same as the current owner
	if msg.Owner != valFound.Owner {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	var token = types.Token{
		Owner:    msg.Owner,
		Denom:    msg.Denom,
		Name:     msg.Name,
		Decimals: msg.Decimals,
		// InitialBalances: msg.InitialBalances,
		// MaxSupply:       msg.MaxSupply,
		// Allowances:      msg.Allowances,
	}

	k.SetToken(ctx, token)

	return &types.MsgUpdateTokenResponse{}, nil
}

func (k msgServer) DeleteToken(goCtx context.Context, msg *types.MsgDeleteToken) (*types.MsgDeleteTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetToken(
		ctx,
		msg.Denom,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg owner is the same as the current owner
	if msg.Owner != valFound.Owner {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveToken(
		ctx,
		msg.Denom,
	)

	return &types.MsgDeleteTokenResponse{}, nil
}
