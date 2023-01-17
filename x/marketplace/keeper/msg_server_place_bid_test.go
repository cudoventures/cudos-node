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
	nfttypes "github.com/CudoVentures/cudos-node/x/nft/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

var (
	accs       = simtypes.RandomAccounts(rand.New(rand.NewSource(rand.Int63())), 2)
	acc1, acc2 = accs[0].Address, accs[1].Address
	fund       = sdk.NewCoins(sdk.NewCoin("acudos", sdk.NewIntFromUint64(100)))
	amount     = sdk.NewCoin("acudos", sdk.OneInt())
	denomId    = "asd"
	tokenId    = "1"
)

func TestMsgServerPlaceBid_EnglishAuction(t *testing.T) {
	for _, tc := range []struct {
		desc    string
		arrange func(
			msg *types.MsgPlaceBid,
			msgServer types.MsgServer,
			bk types.BankKeeper,
			ctx sdk.Context,
		)
		addBlockTime time.Duration
		wantErr      error
	}{
		{
			desc: "valid",
			arrange: func(
				msg *types.MsgPlaceBid,
				msgServer types.MsgServer,
				bk types.BankKeeper,
				ctx sdk.Context,
			) {
			},
		},
		{
			desc: "cannot refund current bidder",
			arrange: func(
				msg *types.MsgPlaceBid,
				msgServer types.MsgServer,
				bk types.BankKeeper,
				ctx sdk.Context,
			) {
				_, err := msgServer.PlaceBid(ctx, msg)
				require.NoError(t, err)

				err = bk.SendCoinsFromModuleToAccount(
					ctx, types.ModuleName, acc2, fund.Add(msg.Amount),
				)
				require.NoError(t, err)

				msg.Amount = msg.Amount.AddAmount(sdk.OneInt())
			},
			wantErr: sdkerrors.ErrInsufficientFunds,
		},
		{
			desc: "insufficient funds bidder",
			arrange: func(
				msg *types.MsgPlaceBid,
				msgServer types.MsgServer,
				bk types.BankKeeper,
				ctx sdk.Context,
			) {
				msg.Amount = fund[0].AddAmount(sdk.OneInt())
			},
			wantErr: sdkerrors.ErrInsufficientFunds,
		},
		{
			desc: "invalid bidder",
			arrange: func(
				msg *types.MsgPlaceBid,
				msgServer types.MsgServer,
				bk types.BankKeeper,
				ctx sdk.Context,
			) {
				msg.Bidder = "invalid"
			},
			wantErr: sdkerrors.ErrInvalidAddress,
		},
		{
			desc: "bid lower than current bid",
			arrange: func(
				msg *types.MsgPlaceBid,
				msgServer types.MsgServer,
				bk types.BankKeeper,
				ctx sdk.Context,
			) {
				_, err := msgServer.PlaceBid(ctx, msg)
				require.NoError(t, err)
			},
			wantErr: types.ErrInvalidPrice,
		},
		{
			desc: "bid lower than min price",
			arrange: func(
				msg *types.MsgPlaceBid,
				msgServer types.MsgServer,
				bk types.BankKeeper,
				ctx sdk.Context,
			) {
				msg.Amount = msg.Amount.SubAmount(sdk.OneInt())
			},
			wantErr: types.ErrInvalidPrice,
		},
		{
			desc: "auction expired",
			arrange: func(
				msg *types.MsgPlaceBid,
				msgServer types.MsgServer,
				bk types.BankKeeper,
				ctx sdk.Context,
			) {
			},
			addBlockTime: time.Hour * 25,
			wantErr:      types.ErrAuctionExpired,
		},
		{
			desc: "bidder same as creator",
			arrange: func(
				msg *types.MsgPlaceBid,
				msgServer types.MsgServer,
				bk types.BankKeeper,
				ctx sdk.Context,
			) {
				msg.Bidder = acc1.String()
			},
			wantErr: types.ErrCannotBuyOwnNft,
		},
		{
			desc: "invalid auction id",
			arrange: func(
				msg *types.MsgPlaceBid,
				msgServer types.MsgServer,
				bk types.BankKeeper,
				ctx sdk.Context,
			) {
				msg.AuctionId++
			},
			wantErr: types.ErrAuctionNotFound,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			ctx, k, bk, _ := setupTestPlaceBid(t)
			msgServer := keeper.NewMsgServerImpl(*k)

			a := types.NewEnglishAuction(
				acc1.String(),
				denomId,
				tokenId,
				amount,
				ctx.BlockTime(),
				ctx.BlockTime().Add(time.Hour*24),
			)
			auctionId, err := k.PublishAuction(ctx, a)
			require.NoError(t, err)

			msg := types.NewMsgPlaceBid(acc2.String(), auctionId, amount)
			tc.arrange(msg, msgServer, bk, ctx)
			ctx = ctx.WithBlockTime(ctx.BlockTime().Add(tc.addBlockTime))

			_, err = msgServer.PlaceBid(ctx, msg)
			require.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestMsgServerPlaceBid_DutchAuction(t *testing.T) {
	for _, tc := range []struct {
		desc    string
		arrange func(
			msg *types.MsgPlaceBid,
			bk types.BankKeeper,
			nk types.NftKeeper,
			ctx sdk.Context,
		)
		addBlockTime time.Duration
		wantErr      error
	}{
		{
			desc: "valid",
			arrange: func(
				msg *types.MsgPlaceBid,
				bk types.BankKeeper,
				nk types.NftKeeper,
				ctx sdk.Context,
			) {
			},
		},
		{
			desc: "err doTrade",
			arrange: func(
				msg *types.MsgPlaceBid,
				bk types.BankKeeper,
				nk types.NftKeeper,
				ctx sdk.Context,
			) {
				err := nk.SoftUnlockNFT(ctx, types.ModuleName, denomId, tokenId)
				require.NoError(t, err)
			},
			wantErr: nfttypes.ErrNotSoftLocked,
		},
		{
			desc: "insufficient funds bidder",
			arrange: func(
				msg *types.MsgPlaceBid,
				bk types.BankKeeper,
				nk types.NftKeeper,
				ctx sdk.Context,
			) {
				err := bk.SendCoinsFromAccountToModule(
					ctx, acc2, types.ModuleName, bk.SpendableCoins(ctx, acc2),
				)
				require.NoError(t, err)

				msg.Amount = fund[0].AddAmount(sdk.OneInt())
			},
			wantErr: sdkerrors.ErrInsufficientFunds,
		},
		{
			desc: "invalid bidder",
			arrange: func(
				msg *types.MsgPlaceBid,
				bk types.BankKeeper,
				nk types.NftKeeper,
				ctx sdk.Context,
			) {
				msg.Bidder = "invalid"
			},
			wantErr: sdkerrors.ErrInvalidAddress,
		},
		{
			desc: "bid lower than current price",
			arrange: func(
				msg *types.MsgPlaceBid,
				bk types.BankKeeper,
				nk types.NftKeeper,
				ctx sdk.Context,
			) {
				msg.Amount = msg.Amount.SubAmount(sdk.OneInt())
			},
			wantErr: types.ErrInvalidPrice,
		},
		{
			desc: "auction expired",
			arrange: func(
				msg *types.MsgPlaceBid,
				bk types.BankKeeper,
				nk types.NftKeeper,
				ctx sdk.Context,
			) {
			},
			addBlockTime: time.Hour * 25,
			wantErr:      types.ErrAuctionExpired,
		},
		{
			desc: "bidder same as creator",
			arrange: func(
				msg *types.MsgPlaceBid,
				bk types.BankKeeper,
				nk types.NftKeeper,
				ctx sdk.Context,
			) {
				msg.Bidder = acc1.String()
			},
			wantErr: types.ErrCannotBuyOwnNft,
		},
		{
			desc: "invalid auction id",
			arrange: func(
				msg *types.MsgPlaceBid,
				bk types.BankKeeper,
				nk types.NftKeeper,
				ctx sdk.Context,
			) {
				msg.AuctionId++
			},
			wantErr: types.ErrAuctionNotFound,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			ctx, k, bk, nk := setupTestPlaceBid(t)
			msgServer := keeper.NewMsgServerImpl(*k)

			a := types.NewDutchAuction(
				acc1.String(),
				denomId,
				tokenId,
				fund[0],
				amount,
				ctx.BlockTime(),
				ctx.BlockTime().Add(time.Hour*24),
			)
			auctionId, err := k.PublishAuction(ctx, a)
			require.NoError(t, err)

			msg := types.NewMsgPlaceBid(acc2.String(), 0, fund[0])
			tc.arrange(msg, bk, nk, ctx)
			ctx = ctx.WithBlockTime(ctx.BlockTime().Add(tc.addBlockTime))

			_, err = msgServer.PlaceBid(ctx, msg)
			require.ErrorIs(t, err, tc.wantErr)

			_, errGetAuction := k.GetAuction(ctx, auctionId)
			a.Creator = msg.Bidder
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

func setupTestPlaceBid(
	t *testing.T,
) (sdk.Context, *keeper.Keeper, types.BankKeeper, types.NftKeeper) {
	k, nk, bk, ctx := keepertest.MarketplaceKeeper(t)

	err := bk.MintCoins(ctx, types.ModuleName, fund.Add(fund...))
	require.NoError(t, err)

	err = bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, acc2, fund)
	require.NoError(t, err)

	err = nk.IssueDenom(ctx, denomId, "asd", "", "", "", acc1.String(), "", "", acc1)
	require.NoError(t, err)

	_, err = nk.MintNFT(ctx, denomId, "asd", "", "", acc1, acc1)
	require.NoError(t, err)
	return ctx, k, bk, nk
}
