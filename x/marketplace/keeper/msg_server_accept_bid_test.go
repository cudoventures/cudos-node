package keeper_test

import (
	"math/rand"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/stretchr/testify/require"

	keepertest "github.com/CudoVentures/cudos-node/testutil/keeper"
	"github.com/CudoVentures/cudos-node/x/marketplace/keeper"
	"github.com/CudoVentures/cudos-node/x/marketplace/types"
)

func TestMsgServerAcceptBid(t *testing.T) {
	accs = simtypes.RandomAccounts(rand.New(rand.NewSource(rand.Int63())), 2)
	acc1, acc2 = accs[0].Address, accs[1].Address
	fund = sdk.NewCoins(sdk.NewCoin("acudos", sdk.NewIntFromUint64(100)))
	amount = sdk.NewCoin("acudos", sdk.OneInt())

	for _, tc := range []struct {
		desc    string
		arrange func(
			msg *types.MsgAcceptBid,
			a types.Auction,
			k *keeper.Keeper,
			bk types.BankKeeper,
			ctx sdk.Context,
		)
		addBlockTime time.Duration
		wantErr      error
	}{
		{
			desc: "valid",
			arrange: func(
				msg *types.MsgAcceptBid,
				a types.Auction,
				k *keeper.Keeper,
				bk types.BankKeeper,
				ctx sdk.Context,
			) {
				auctionId, err := k.PublishAuction(ctx, a)
				require.NoError(t, err)

				err = k.PlaceBid(ctx, auctionId, types.Bid{
					Amount: amount,
					Bidder: acc2.String(),
				})
				require.NoError(t, err)
			},
		},
		{
			desc: "err doTrade",
			arrange: func(
				msg *types.MsgAcceptBid,
				a types.Auction,
				k *keeper.Keeper,
				bk types.BankKeeper,
				ctx sdk.Context,
			) {
				auctionId, err := k.PublishAuction(ctx, a)
				require.NoError(t, err)

				err = k.PlaceBid(ctx, auctionId, types.Bid{
					Amount: amount,
					Bidder: acc2.String(),
				})
				require.NoError(t, err)

				err = bk.SendCoinsFromModuleToAccount(
					ctx, types.ModuleName, acc2, fund.Add(amount),
				)
				require.NoError(t, err)
			},
			wantErr: sdkerrors.ErrInsufficientFunds,
		},
		{
			desc: "no current bid",
			arrange: func(
				msg *types.MsgAcceptBid,
				a types.Auction,
				k *keeper.Keeper,
				bk types.BankKeeper,
				ctx sdk.Context,
			) {
				_, err := k.PublishAuction(ctx, a)
				require.NoError(t, err)
			},
			wantErr: sdkerrors.ErrInvalidRequest,
		},
		{
			desc: "not english auction",
			arrange: func(
				msg *types.MsgAcceptBid,
				a types.Auction,
				k *keeper.Keeper,
				bk types.BankKeeper,
				ctx sdk.Context,
			) {
				_, err := k.PublishAuction(ctx, types.NewDutchAuction(
					a.GetCreator(),
					a.GetDenomId(),
					a.GetTokenId(),
					fund[0],
					amount,
					ctx.BlockTime(),
					ctx.BlockTime().Add(time.Hour*24),
				))
				require.NoError(t, err)
			},
			wantErr: sdkerrors.ErrInvalidRequest,
		},
		{
			desc: "auction expired",
			arrange: func(
				msg *types.MsgAcceptBid,
				a types.Auction,
				k *keeper.Keeper,
				bk types.BankKeeper,
				ctx sdk.Context,
			) {
				_, err := k.PublishAuction(ctx, a)
				require.NoError(t, err)
			},
			addBlockTime: time.Hour * 25,
			wantErr:      types.ErrAuctionExpired,
		},
		{
			desc: "not owner",
			arrange: func(
				msg *types.MsgAcceptBid,
				a types.Auction,
				k *keeper.Keeper,
				bk types.BankKeeper,
				ctx sdk.Context,
			) {
				_, err := k.PublishAuction(ctx, a)
				require.NoError(t, err)

				msg.Sender = acc2.String()
			},
			wantErr: types.ErrNotNftOwner,
		},
		{
			desc: "invalid auction id",
			arrange: func(
				msg *types.MsgAcceptBid,
				a types.Auction,
				k *keeper.Keeper,
				bk types.BankKeeper,
				ctx sdk.Context,
			) {
				_, err := k.PublishAuction(ctx, a)
				require.NoError(t, err)

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

			err = bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, acc2, fund)
			require.NoError(t, err)

			err = nk.IssueDenom(
				ctx, denomId, "asd", "", "", "", acc1.String(), "", "", acc1,
			)
			require.NoError(t, err)

			_, err = nk.MintNFT(ctx, denomId, "asd", "", "", acc1, acc1)
			require.NoError(t, err)

			a := types.NewEnglishAuction(
				acc1.String(),
				denomId,
				tokenId,
				amount,
				ctx.BlockTime(),
				ctx.BlockTime().Add(time.Hour*24),
			)
			msg := types.NewMsgAcceptBid(acc1.String(), 0)
			tc.arrange(msg, a, k, bk, ctx)
			ctx = ctx.WithBlockTime(ctx.BlockTime().Add(tc.addBlockTime))

			_, err = msgServer.AcceptBid(ctx, msg)
			require.ErrorIs(t, err, tc.wantErr)

			_, errGetAuction := k.GetAuction(ctx, 0)
			a.Creator = acc2.String()
			_, errPublishAuction := k.PublishAuction(ctx, a)
			if tc.wantErr == nil {
				require.ErrorIs(t, errGetAuction, types.ErrAuctionNotFound)
				require.NoError(t, errPublishAuction)
			} else {
				require.NoError(t, errGetAuction)
				require.ErrorIs(t, errPublishAuction, types.ErrNftAlreadyPublished)
			}
		})
	}
}
