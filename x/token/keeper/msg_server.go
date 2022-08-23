package keeper

import (
	"context"

	"github.com/CudoVentures/cudos-node/x/token/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) CreateToken(goCtx context.Context, msg *types.MsgCreateToken) (*types.MsgCreateTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetTokenByDenom(
		ctx,
		msg.Denom,
	)
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "denom already exists")
	}

	var token = types.Token{
		Owner:     msg.Owner,
		Denom:     msg.Denom,
		Name:      msg.Name,
		Decimals:  msg.Decimals,
		MaxSupply: msg.MaxSupply,
	}

	k.SaveToken(
		ctx,
		token,
	)

	if len(msg.InitialBalances) > 0 {
		// todo iterate through the balances and mint properly
		coins := sdk.NewCoins(sdk.NewCoin(msg.Denom, sdk.NewInt(123)))
		if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
			return nil, err
		}

		addr, err := sdk.AccAddressFromBech32(msg.Owner)
		if err != nil {
			return nil, err
		}

		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, coins); err != nil {
			return nil, err
		}
	}

	return &types.MsgCreateTokenResponse{}, nil
}
