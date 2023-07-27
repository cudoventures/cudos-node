package simulation

import (
	"fmt"
	"math/rand"
	"strings"

	simappparams "cosmossdk.io/simapp/params"
	"github.com/CudoVentures/cudos-node/x/marketplace/keeper"
	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgMintNft(
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

		collection := types.NewCollection(denom, []types.Royalty{}, []types.Royalty{}, ownerAddr.String(), true)
		_, err := k.PublishCollection(ctx, collection)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventMintNftType, err.Error()), nil, err
		}

		// Mint NFT
		nftPrice := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(r.Int63n(100)+1))
		fundAcc(ctx, bk, ownerAddr, nftPrice)
		name := strings.ToLower(simtypes.RandStringOfLength(r, 10))
		uri := strings.ToLower(simtypes.RandStringOfLength(r, 10))
		data := strings.ToLower(simtypes.RandStringOfLength(r, 10))

		recipientAcc, _ := simtypes.RandomAcc(r, accs)

		msg := types.NewMsgMintNft(
			ownerAddr.String(),
			denom,
			recipientAcc.Address.String(),
			name,
			uri,
			data,
			"",
			nftPrice,
		)

		account := ak.GetAccount(ctx, ownerAddr)

		ownerAccount, found := simtypes.FindAccount(accs, ownerAddr)
		if !found {
			err := fmt.Errorf("account %s not found", msg.Creator)
			return simtypes.NoOpMsg(types.ModuleName, types.EventMintNftType, err.Error()), nil, err
		}

		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventMintNftType, err.Error()), nil, err
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
			return simtypes.NoOpMsg(types.ModuleName, types.EventMintNftType, err.Error()), nil, err
		}

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "MintNft simulation not implemented"), nil, nil
	}
}
