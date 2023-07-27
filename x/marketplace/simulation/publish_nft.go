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

func SimulateMsgPublishNft(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	nk types.NftKeeper,
	k keeper.Keeper,
	tr *TokensRandomizer,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		ownerAddr, denom, nftID := tr.GetRandomTokenIdWithOwner(ctx, nk, r, true, "")
		if ownerAddr.Empty() {
			err := fmt.Errorf("invalid account")
			return simtypes.NoOpMsg(types.ModuleName, types.EventPublishNftType, err.Error()), nil, err
		}

		nftPrice := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(r.Int63n(100)+1))
		fundAcc(ctx, bk, ownerAddr, nftPrice)
		msg := types.NewMsgPublishNft(
			ownerAddr.String(),
			nftID,
			denom,
			nftPrice,
		)

		account := ak.GetAccount(ctx, ownerAddr)

		ownerAccount, found := simtypes.FindAccount(accs, ownerAddr)
		if !found {
			err := fmt.Errorf("account %s not found", msg.Creator)
			return simtypes.NoOpMsg(types.ModuleName, types.EventPublishNftType, err.Error()), nil, err
		}

		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventPublishNftType, err.Error()), nil, err
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
			ownerAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		if _, _, err = app.SimDeliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventPublishNftType, err.Error()), nil, err
		}

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "PublishNft simulation not implemented"), nil, nil
	}
}
