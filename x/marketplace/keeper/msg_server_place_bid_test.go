package keeper_test

import (
	"math/rand"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/CudoVentures/cudos-node/testutil/keeper"
	"github.com/CudoVentures/cudos-node/x/marketplace/keeper"
	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func TestMsgServerPlaceBid(t *testing.T) {
	r := rand.New(rand.NewSource(rand.Int63()))
	accs := simtypes.RandomAccounts(r, 2)
	auctionId := uint64(0)
	fund := sdk.NewCoins(sdk.NewCoin("acudos", sdk.NewIntFromUint64(10000)))
	amount := sdk.NewCoin("acudos", sdk.OneInt())

	for _, tc := range []struct {
		desc         string
		arrange      func(msg *types.MsgPlaceBid, msgServer types.MsgServer, ctx sdk.Context)
		addBlockTime time.Duration
		wantErr      error
	}{
		{
			desc:    "valid",
			arrange: func(msg *types.MsgPlaceBid, msgServer types.MsgServer, ctx sdk.Context) {},
		},
		// todo dutch auction
		{
			desc: "insufficient balance",
			arrange: func(msg *types.MsgPlaceBid, msgServer types.MsgServer, ctx sdk.Context) {
				msg.Amount = fund[0].AddAmount(sdk.OneInt())
			},
			wantErr: sdkerrors.ErrInsufficientFunds,
		},
		{
			desc: "invalid bidder",
			arrange: func(msg *types.MsgPlaceBid, msgServer types.MsgServer, ctx sdk.Context) {
				msg.Bidder = "invalid"
			},
			wantErr: sdkerrors.ErrInvalidAddress,
		},
		{
			desc: "bid lower than current bid",
			arrange: func(msg *types.MsgPlaceBid, msgServer types.MsgServer, ctx sdk.Context) {
				_, err := msgServer.PlaceBid(ctx, msg)
				require.NoError(t, err)
			},
			wantErr: types.ErrInvalidPrice,
		},
		{
			desc: "bid lower than min price",
			arrange: func(msg *types.MsgPlaceBid, msgServer types.MsgServer, ctx sdk.Context) {
				msg.Amount = msg.Amount.SubAmount(sdk.OneInt())
			},
			wantErr: types.ErrInvalidPrice,
		},
		{
			desc:         "auction expired",
			arrange:      func(msg *types.MsgPlaceBid, msgServer types.MsgServer, ctx sdk.Context) {},
			addBlockTime: time.Hour * 25,
			wantErr:      types.ErrAuctionExpired,
		},
		{
			desc: "bidder same as creator",
			arrange: func(msg *types.MsgPlaceBid, msgServer types.MsgServer, ctx sdk.Context) {
				msg.Bidder = accs[0].Address.String()
			},
			wantErr: types.ErrCannotBuyOwnNft,
		},
		{
			desc: "invalid auction id",
			arrange: func(msg *types.MsgPlaceBid, msgServer types.MsgServer, ctx sdk.Context) {
				msg.AuctionId++
			},
			wantErr: types.ErrAuctionNotFound,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, nk, bk, ctx := keepertest.MarketplaceKeeper(t)
			msgServer := keeper.NewMsgServerImpl(*k)

			err := bk.MintCoins(ctx, types.ModuleName, fund.Add(fund...))
			require.NoError(t, err)

			err = bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, accs[1].Address, fund)
			require.NoError(t, err)

			err = nk.IssueDenom(ctx, "asd", "asd", "{a:a,b:b}", "asd", "", accs[0].Address.String(), "", "", accs[0].Address)
			require.NoError(t, err)

			_, err = nk.MintNFT(ctx, "asd", "asd", "", "", accs[0].Address, accs[0].Address)
			require.NoError(t, err)

			msgPublishAuction, err := types.NewMsgPublishAuction(accs[0].Address.String(), "asd", "1", time.Hour*24, &types.EnglishAuction{MinPrice: amount})
			require.NoError(t, err)

			_, err = msgServer.PublishAuction(ctx, msgPublishAuction)
			require.NoError(t, err)

			msg := &types.MsgPlaceBid{
				AuctionId: auctionId,
				Amount:    amount,
				Bidder:    accs[1].Address.String(),
			}

			tc.arrange(msg, msgServer, ctx)

			ctx = ctx.WithBlockTime(ctx.BlockTime().Add(tc.addBlockTime))
			_, err = msgServer.PlaceBid(ctx, msg)
			require.ErrorIs(t, err, tc.wantErr)
		})
	}
}
