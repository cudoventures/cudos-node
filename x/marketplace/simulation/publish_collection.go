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

func SimulateMsgPublishCollection(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	nk types.NftKeeper,
	k keeper.Keeper,
	dr *DenomsRandomizer,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		ownerAddr, denom := dr.GetRandomDenomIdWithOwner(ctx, nk, r, true)
		if ownerAddr.Empty() {
			err := fmt.Errorf("invalid account")
			return simtypes.NoOpMsg(types.ModuleName, types.EventPublishCollectionType, err.Error()), nil, err
		}

		msg := types.NewMsgPublishCollection(
			ownerAddr.String(),
			denom,
			[]types.Royalty{},
			[]types.Royalty{},
		)

		account := ak.GetAccount(ctx, ownerAddr)

		ownerAccount, found := simtypes.FindAccount(accs, ownerAddr)
		if !found {
			err := fmt.Errorf("account %s not found", msg.Creator)
			return simtypes.NoOpMsg(types.ModuleName, types.EventPublishCollectionType, err.Error()), nil, err
		}

		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventPublishCollectionType, err.Error()), nil, err
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
			return simtypes.NoOpMsg(types.ModuleName, types.EventPublishCollectionType, err.Error()), nil, err
		}

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "PublishCollection simulation not implemented"), nil, nil
	}
}
