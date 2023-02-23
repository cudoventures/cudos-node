package simulation

import (
	"fmt"
	"math/rand"

	"github.com/CudoVentures/cudos-node/x/marketplace/keeper"
	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	nftsim "github.com/CudoVentures/cudos-node/x/nft/simulation"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgBuyNft(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	nk types.NftKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		// Publish NFT for sale

		sellerAddr, denom, nftID := nftsim.GetRandomNFTFromOwner(ctx, nk, r)
		if sellerAddr.Empty() {
			err := fmt.Errorf("invalid account")
			return simtypes.NoOpMsg(types.ModuleName, types.EventBuyNftType, err.Error()), nil, err
		}

		nft := types.NewNft(nftID, denom, sellerAddr.String(), sdk.NewCoin("acudos", sdk.NewInt(100000000)))
		publishedNftID, err := k.PublishNFT(ctx, nft)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventBuyNftType, err.Error()), nil, err
		}

		// Buy the NFT

		buyerAcc, _ := simtypes.RandomAcc(r, accs)
		account := ak.GetAccount(ctx, buyerAcc.Address)

		msg := types.NewMsgBuyNft(buyerAcc.Address.String(), publishedNftID)

		spendable := bk.SpendableCoins(ctx, buyerAcc.Address)
		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventBuyNftType, err.Error()), nil, err
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
			buyerAcc.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		if _, _, err = app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventBuyNftType, err.Error()), nil, err
		}

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "BuyNft simulation not implemented"), nil, nil
	}
}
