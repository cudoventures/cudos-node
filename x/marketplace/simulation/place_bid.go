package simulation

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/CudoVentures/cudos-node/x/marketplace/keeper"
	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	nftsim "github.com/CudoVentures/cudos-node/x/nft/simulation"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgPlaceBidEnglishAuction(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	nk types.NftKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return simulateMsgPlaceBid(ak, bk, nk, k, &types.EnglishAuction{MinPrice: sdk.NewCoin("acudos", sdk.NewInt(100000000000))})
}

func SimulateMsgPlaceBidDutchAuction(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	nk types.NftKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	// todo dutch auction
	return nil
}

func simulateMsgPlaceBid(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	nk types.NftKeeper,
	k keeper.Keeper,
	at types.AuctionType,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		seller, denomId, tokenId := nftsim.GetRandomNFTFromOwner(ctx, nk, r)
		if seller.Empty() {
			err := fmt.Errorf("invalid account")
			return simtypes.NoOpMsg(types.ModuleName, types.EventPlaceBidType, err.Error()), nil, err
		}

		a, err := types.NewAuction(
			seller.String(),
			denomId,
			tokenId,
			time.Now().Add(time.Hour*25),
			at,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventPlaceBidType, err.Error()), nil, err
		}

		auctionId, err := k.PublishAuction(ctx, a)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventPlaceBidType, err.Error()), nil, err
		}

		bidder, _ := simtypes.RandomAcc(r, accs)
		account := ak.GetAccount(ctx, bidder.Address)

		msg := types.NewMsgPlaceBid(bidder.Address.String(), auctionId, sdk.NewCoin("acudos", sdk.NewInt(100000000000)))

		spendable := bk.SpendableCoins(ctx, bidder.Address)
		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventPlaceBidType, err.Error()), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			bidder.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		if _, _, err = app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventPlaceBidType, err.Error()), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}
