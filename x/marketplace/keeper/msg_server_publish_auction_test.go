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
	nftkeeper "github.com/CudoVentures/cudos-node/x/nft/keeper"
	nfttypes "github.com/CudoVentures/cudos-node/x/nft/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
		wantErr error
	}{
		{
			desc:    "valid english auction",
			arrange: func(msg *types.MsgPublishAuction, msgServer types.MsgServer, nk *nftkeeper.Keeper, ctx sdk.Context) {},
		},
		// todo dutch auction
		{
			desc: "valid approved nft address",
			arrange: func(msg *types.MsgPublishAuction, msgServer types.MsgServer, nk *nftkeeper.Keeper, ctx sdk.Context) {
				msg.Creator = accs[1].Address.String()
				nftMsgServer := nftkeeper.NewMsgServerImpl(*nk)
				_, err := nftMsgServer.ApproveNft(
					ctx,
					&nfttypes.MsgApproveNft{
						Id:              tokenId,
						DenomId:         denomId,
						Sender:          accs[0].Address.String(),
						ApprovedAddress: accs[1].Address.String(),
					},
				)
				require.NoError(t, err)
			},
		},
		{
			desc: "valid approved operator",
			arrange: func(msg *types.MsgPublishAuction, msgServer types.MsgServer, nk *nftkeeper.Keeper, ctx sdk.Context) {
				msg.Creator = accs[1].Address.String()
				nftMsgServer := nftkeeper.NewMsgServerImpl(*nk)
				_, err := nftMsgServer.ApproveAllNft(
					ctx,
					&nfttypes.MsgApproveAllNft{
						Operator: accs[1].Address.String(),
						Sender:   accs[0].Address.String(),
						Approved: true,
					},
				)
				require.NoError(t, err)
			},
		},
		{
			desc: "already locked",
			arrange: func(msg *types.MsgPublishAuction, msgServer types.MsgServer, nk *nftkeeper.Keeper, ctx sdk.Context) {
				err := nk.SoftLockNFT(ctx, accs[0].Address.String(), denomId, tokenId)
				require.NoError(t, err)
			},
			wantErr: nfttypes.ErrAlreadySoftLocked,
		},
		{
			desc: "already published auction",
			arrange: func(msg *types.MsgPublishAuction, msgServer types.MsgServer, nk *nftkeeper.Keeper, ctx sdk.Context) {
				_, err := msgServer.PublishAuction(ctx, msg)
				require.NoError(t, err)
			},
			wantErr: types.ErrNftAlreadyPublished,
		},
		{
			desc: "already published nft",
			arrange: func(msg *types.MsgPublishAuction, msgServer types.MsgServer, nk *nftkeeper.Keeper, ctx sdk.Context) {
				_, err := msgServer.PublishNft(
					ctx,
					&types.MsgPublishNft{
						Creator: accs[0].Address.String(),
						TokenId: tokenId,
						DenomId: denomId,
						Price:   sdk.NewCoin("stake", sdk.OneInt()),
					},
				)
				require.NoError(t, err)
			},
			wantErr: types.ErrNftAlreadyPublished,
		},
		{
			desc: "invalid AuctionType",
			arrange: func(msg *types.MsgPublishAuction, msgServer types.MsgServer, nk *nftkeeper.Keeper, ctx sdk.Context) {
				msg.AuctionType = &codectypes.Any{}
			},
			wantErr: sdkerrors.ErrInvalidType,
		},
		{
			desc: "not owner",
			arrange: func(msg *types.MsgPublishAuction, msgServer types.MsgServer, nk *nftkeeper.Keeper, ctx sdk.Context) {
				msg.Creator = accs[1].Address.String()
			},
			wantErr: types.ErrNotNftOwner,
		},
		{
			desc: "not existing nft",
			arrange: func(msg *types.MsgPublishAuction, msgServer types.MsgServer, nk *nftkeeper.Keeper, ctx sdk.Context) {
				msg.TokenId = tokenId + tokenId
			},
			wantErr: types.ErrNftNotFound,
		},
		{
			desc: "invalid creator",
			arrange: func(msg *types.MsgPublishAuction, msgServer types.MsgServer, nk *nftkeeper.Keeper, ctx sdk.Context) {
				msg.Creator = ""
			},
			wantErr: sdkerrors.ErrInvalidAddress,
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
			require.ErrorIs(t, err, tc.wantErr)
		})
	}
}
