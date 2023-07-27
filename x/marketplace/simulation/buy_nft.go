package simulation

import (
	"fmt"
	"math/rand"

	simappparams "cosmossdk.io/simapp/params"
	"github.com/CudoVentures/cudos-node/x/marketplace/keeper"
	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgBuyNft(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	nk types.NftKeeper,
	k keeper.Keeper,
	tr *TokensRandomizer,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		buyerAcc, _ := simtypes.RandomAcc(r, accs)
		account := ak.GetAccount(ctx, buyerAcc.Address)

		// Publish NFT for sale
		sellerAddr, denom, nftID := tr.GetRandomTokenIdWithOwner(ctx, nk, r, true, buyerAcc.Address.String())
		if sellerAddr.Empty() {
			err := fmt.Errorf("invalid account")
			return simtypes.NoOpMsg(types.ModuleName, types.EventBuyNftType, err.Error()), nil, err
		}

		nftPrice := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(r.Int63n(100)+1))
		nft := types.NewNft(nftID, denom, sellerAddr.String(), nftPrice)
		publishedNftID, err := k.PublishNFT(ctx, nft)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventBuyNftType, err.Error()), nil, err
		}

		// Buy the NFT
		// buyerAcc, _ := simtypes.RandomAcc(r, accs)
		// account := ak.GetAccount(ctx, buyerAcc.Address)

		fundAcc(ctx, bk, account.GetAddress(), nftPrice)
		msg := types.NewMsgBuyNft(buyerAcc.Address.String(), publishedNftID)

		spendable := bk.SpendableCoins(ctx, buyerAcc.Address)
		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventBuyNftType, err.Error()), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := simtestutil.GenSignedMockTx(
			r,
			txGen,
			[]sdk.Msg{msg},
			fees,
			simtestutil.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			buyerAcc.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		if _, _, err = app.SimDeliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventBuyNftType, err.Error()), nil, err
		}

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "BuyNft simulation not implemented"), nil, nil
	}
}
