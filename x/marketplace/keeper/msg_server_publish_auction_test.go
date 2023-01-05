package keeper_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/CudoVentures/cudos-node/testutil/keeper"
	"github.com/CudoVentures/cudos-node/x/marketplace/keeper"
	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	nftkeeper "github.com/CudoVentures/cudos-node/x/nft/keeper"
	nfttypes "github.com/CudoVentures/cudos-node/x/nft/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func TestMsgServerPublishAuction(t *testing.T) {
	r := rand.New(rand.NewSource(rand.Int63()))
	accs := simtypes.RandomAccounts(r, 2)
	const denomId = "asd"
	const tokenId = "1"

	for _, tc := range []struct {
		desc    string
		arrange func(msg *types.MsgPublishAuction, msgServer types.MsgServer, nk *nftkeeper.Keeper, ctx sdk.Context)
		errMsg  string
	}{
		{
			desc:    "valid english auction",
			arrange: func(msg *types.MsgPublishAuction, msgServer types.MsgServer, nk *nftkeeper.Keeper, ctx sdk.Context) {},
		},
		// todo valid dutch auction test case
		{
			desc: "valid approved nft address",
			arrange: func(msg *types.MsgPublishAuction, msgServer types.MsgServer, nk *nftkeeper.Keeper, ctx sdk.Context) {
				msg.Creator = accs[1].Address.String()
				nftMsgServer := nftkeeper.NewMsgServerImpl(*nk)
				_, err := nftMsgServer.ApproveNft(ctx, &nfttypes.MsgApproveNft{tokenId, denomId, accs[0].Address.String(), accs[1].Address.String(), ""})
				require.NoError(t, err)
			},
		},
		{
			desc: "valid approved operator",
			arrange: func(msg *types.MsgPublishAuction, msgServer types.MsgServer, nk *nftkeeper.Keeper, ctx sdk.Context) {
				msg.Creator = accs[1].Address.String()
				nftMsgServer := nftkeeper.NewMsgServerImpl(*nk)
				_, err := nftMsgServer.ApproveAllNft(ctx, &nfttypes.MsgApproveAllNft{accs[1].Address.String(), accs[0].Address.String(), true, ""})
				require.NoError(t, err)
			},
		},
		{
			desc: "already locked",
			arrange: func(msg *types.MsgPublishAuction, msgServer types.MsgServer, nk *nftkeeper.Keeper, ctx sdk.Context) {
				err := nk.SoftLockNFT(ctx, accs[0].Address.String(), denomId, tokenId)
				require.NoError(t, err)
			},
			errMsg: fmt.Sprintf("Failed to acquire soft lock on Denom asd NFT 1 for marketplace because already acquired by %s: already soft locked", accs[0].Address.String()),
		},
		{
			desc: "already published auction",
			arrange: func(msg *types.MsgPublishAuction, msgServer types.MsgServer, nk *nftkeeper.Keeper, ctx sdk.Context) {
				_, err := msgServer.PublishAuction(ctx, msg)
				require.NoError(t, err)
			},
			errMsg: "nft is already published: nft already published",
		},
		{
			desc: "already published nft",
			arrange: func(msg *types.MsgPublishAuction, msgServer types.MsgServer, nk *nftkeeper.Keeper, ctx sdk.Context) {
				_, err := msgServer.PublishNft(ctx, &types.MsgPublishNft{accs[0].Address.String(), tokenId, denomId, sdk.NewCoin("stake", sdk.OneInt())})
				require.NoError(t, err)
			},
			errMsg: "nft is already published: nft already published",
		},
		{
			desc: "invalid AuctionType",
			arrange: func(msg *types.MsgPublishAuction, msgServer types.MsgServer, nk *nftkeeper.Keeper, ctx sdk.Context) {
				msg.AuctionType = &codectypes.Any{}
			},
			errMsg: "expected <nil>, got <nil>: invalid type",
		},
		{
			desc: "not owner",
			arrange: func(msg *types.MsgPublishAuction, msgServer types.MsgServer, nk *nftkeeper.Keeper, ctx sdk.Context) {
				msg.Creator = accs[1].Address.String()
			},
			errMsg: fmt.Sprintf("%s not nft owner or approved operator for token id (1) from denom (asd): not nft owner", accs[1].Address.String()),
		},
		{
			desc: "not existing nft",
			arrange: func(msg *types.MsgPublishAuction, msgServer types.MsgServer, nk *nftkeeper.Keeper, ctx sdk.Context) {
				msg.TokenId = tokenId + tokenId
			},
			errMsg: "not found NFT: asd: unknown nft collection",
		},
		{
			desc: "invalid creator",
			arrange: func(msg *types.MsgPublishAuction, msgServer types.MsgServer, nk *nftkeeper.Keeper, ctx sdk.Context) {
				msg.Creator = ""
			},
			errMsg: "empty address string is not allowed",
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {

			k, nk, _, ctx := keepertest.MarketplaceKeeper(t)
			msgServer := keeper.NewMsgServerImpl(*k)

			err := nk.IssueDenom(ctx, denomId, "asd", "{a:a,b:b}", "asd", "", accs[0].Address.String(), "", "", accs[0].Address)
			require.NoError(t, err)

			_, err = nk.MintNFT(ctx, denomId, "asd", "", "", accs[0].Address, accs[0].Address)
			require.NoError(t, err)

			msg, err := types.NewMsgPublishAuction(accs[0].Address.String(), denomId, tokenId, time.Hour*24, &types.EnglishAuction{MinPrice: sdk.NewCoin("stake", sdk.OneInt())})
			require.NoError(t, err)

			tc.arrange(msg, msgServer, nk, ctx)

			_, err = msgServer.PublishAuction(ctx, msg)

			if tc.errMsg != "" {
				require.EqualError(t, err, tc.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
