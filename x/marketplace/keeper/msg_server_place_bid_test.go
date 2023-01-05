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
		errMsg       string
	}{
		{
			desc:    "valid",
			arrange: func(msg *types.MsgPlaceBid, msgServer types.MsgServer, ctx sdk.Context) {},
		},
		// todo dutch auction tests
		{
			desc: "insufficient balance",
			arrange: func(msg *types.MsgPlaceBid, msgServer types.MsgServer, ctx sdk.Context) {
				msg.Amount = fund[0].AddAmount(sdk.OneInt())
			},
			errMsg: "10000acudos is smaller than 10001acudos: insufficient funds",
		},
		{
			desc: "invalid bidder",
			arrange: func(msg *types.MsgPlaceBid, msgServer types.MsgServer, ctx sdk.Context) {
				msg.Bidder = "invalid"
			},
			errMsg: "decoding bech32 failed: invalid bech32 string length 7",
		},
		{
			desc: "bid lower than current bid",
			arrange: func(msg *types.MsgPlaceBid, msgServer types.MsgServer, ctx sdk.Context) {
				_, err := msgServer.PlaceBid(ctx, msg)
				require.NoError(t, err)
			},
			errMsg: "bid is lower than current bid: invalid price",
		},
		{
			desc: "bid lower than min price",
			arrange: func(msg *types.MsgPlaceBid, msgServer types.MsgServer, ctx sdk.Context) {
				msg.Amount = msg.Amount.SubAmount(sdk.OneInt())
			},
			errMsg: "bid is lower than auction minimum price: invalid price",
		},
		{
			desc:         "auction expired",
			arrange:      func(msg *types.MsgPlaceBid, msgServer types.MsgServer, ctx sdk.Context) {},
			addBlockTime: time.Hour * 25,
			errMsg:       "cannot place a bid for inactive auction 0: auction expired",
		},
		{
			desc: "bidder same as creator",
			arrange: func(msg *types.MsgPlaceBid, msgServer types.MsgServer, ctx sdk.Context) {
				msg.Bidder = accs[0].Address.String()
			},
			errMsg: "cannot bid own auctions: cannot buy own nft",
		},
		{
			desc: "invalid auction id",
			arrange: func(msg *types.MsgPlaceBid, msgServer types.MsgServer, ctx sdk.Context) {
				msg.AuctionId++
			},
			errMsg: "auction with id (1) does not exist: auction not found",
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, nk, bk, ctx := keepertest.MarketplaceKeeper(t)
			msgServer := keeper.NewMsgServerImpl(*k)

			if err := bk.MintCoins(ctx, types.ModuleName, fund); err != nil {
				panic(err)
			}
			if err := bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, accs[1].Address, fund); err != nil {
				panic(err)
			}

			err := nk.IssueDenom(ctx, "asd", "asd", "{a:a,b:b}", "asd", "", accs[0].Address.String(), "", "", accs[0].Address)
			require.NoError(t, err)

			_, err = nk.MintNFT(ctx, "asd", "asd", "", "", accs[0].Address, accs[0].Address)
			require.NoError(t, err)

			msgPublishAuction, err := types.NewMsgPublishAuction(accs[0].Address.String(), "asd", "1", time.Hour*24, &types.EnglishAuction{MinPrice: amount})
			require.NoError(t, err)

			_, err = msgServer.PublishAuction(ctx, msgPublishAuction)
			require.NoError(t, err)

			msg := &types.MsgPlaceBid{auctionId, amount, accs[1].Address.String()}

			tc.arrange(msg, msgServer, ctx)

			if tc.addBlockTime > 0 {
				ctx = ctx.WithBlockTime(time.Now().Add(tc.addBlockTime))
			}
			_, err = msgServer.PlaceBid(ctx, msg)

			if tc.errMsg != "" {
				require.EqualError(t, err, tc.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
