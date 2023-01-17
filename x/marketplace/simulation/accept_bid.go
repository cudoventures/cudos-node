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

func SimulateMsgAcceptBid(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	nk types.NftKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand,
		app *baseapp.BaseApp,
		ctx sdk.Context,
		accs []simtypes.Account,
		chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		owner, denomId, tokenId := nftsim.GetRandomNFTFromOwner(ctx, nk, r)
		if owner.Empty() {
			err := fmt.Errorf("invalid account")
			op := simtypes.NoOpMsg(types.ModuleName, types.EventAcceptBidType, err.Error())
			return op, nil, err
		}

		a := types.NewEnglishAuction(
			owner.String(),
			denomId,
			tokenId,
			sdk.NewCoin("acudos", sdk.NewInt(10)),
			ctx.BlockTime(),
			ctx.BlockTime().Add(time.Hour*24),
		)

		auctionId, err := k.PublishAuction(ctx, a)
		if err != nil {
			op := simtypes.NoOpMsg(types.ModuleName, types.EventAcceptBidType, err.Error())
			return op, nil, err
		}

		bidder, _ := simtypes.RandomAcc(r, accs)
		account := ak.GetAccount(ctx, bidder.Address)
		balance := bk.SpendableCoins(ctx, bidder.Address)
		fees, err := simtypes.RandomFees(r, ctx, balance)
		if err != nil {
			op := simtypes.NoOpMsg(types.ModuleName, types.EventAcceptBidType, err.Error())
			return op, nil, err
		}

		err = k.PlaceBid(ctx, auctionId, types.Bid{
			Amount: sdk.NewCoin("acudos", sdk.NewInt(10)),
			Bidder: bidder.Address.String(),
		})
		if err != nil {
			op := simtypes.NoOpMsg(types.ModuleName, types.EventAcceptBidType, err.Error())
			return op, nil, err
		}

		msg := types.NewMsgAcceptBid(owner.String(), auctionId)

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
			op := simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error())
			return op, nil, err
		}

		if _, _, err = app.Deliver(txGen.TxEncoder(), tx); err != nil {
			op := simtypes.NoOpMsg(types.ModuleName, types.EventAcceptBidType, err.Error())
			return op, nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}
