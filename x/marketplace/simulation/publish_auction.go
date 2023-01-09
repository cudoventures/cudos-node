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

func SimulateMsgPublishEnglishAuction(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	nk types.NftKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return simulateMsgPublishAuction(ak, bk, nk, k, &types.EnglishAuction{
		MinPrice: sdk.NewCoin("acudos", sdk.NewInt(10)),
	})
}

func SimulateMsgPublishDutchAuction(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	nk types.NftKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return simulateMsgPublishAuction(ak, bk, nk, k, &types.DutchAuction{
		StartPrice: sdk.NewCoin("acudos", sdk.NewInt(10)),
		MinPrice:   sdk.NewCoin("acudos", sdk.NewInt(1)),
	})
}

func simulateMsgPublishAuction(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	nk types.NftKeeper,
	k keeper.Keeper,
	a types.Auction,
) simtypes.Operation {
	return func(
		r *rand.Rand,
		app *baseapp.BaseApp,
		ctx sdk.Context,
		accs []simtypes.Account,
		chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		module := types.ModuleName
		owner, denomId, tokenId := nftsim.GetRandomNFTFromOwner(ctx, nk, r)
		if owner.Empty() {
			err := fmt.Errorf("invalid account")
			op := simtypes.NoOpMsg(module, types.EventPublishAuctionType, err.Error())
			return op, nil, err
		}

		msg, err := types.NewMsgPublishAuction(
			owner.String(),
			denomId,
			tokenId,
			time.Hour*24,
			a,
		)
		if err != nil {
			op := simtypes.NoOpMsg(module, types.EventPublishAuctionType, err.Error())
			return op, nil, err
		}

		account := ak.GetAccount(ctx, owner)

		ownerAccount, found := simtypes.FindAccount(accs, owner)
		if !found {
			err := fmt.Errorf("account %s not found", msg.Creator)
			op := simtypes.NoOpMsg(module, types.EventPublishAuctionType, err.Error())
			return op, nil, err
		}

		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			op := simtypes.NoOpMsg(module, types.EventPublishAuctionType, err.Error())
			return op, nil, err
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
			ownerAccount.PrivKey,
		)
		if err != nil {
			op := simtypes.NoOpMsg(module, msg.Type(), err.Error())
			return op, nil, err
		}

		if _, _, err = app.Deliver(txGen.TxEncoder(), tx); err != nil {
			op := simtypes.NoOpMsg(module, types.EventPublishAuctionType, err.Error())
			return op, nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}
