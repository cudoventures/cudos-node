package simulation

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/CudoVentures/cudos-node/x/marketplace/keeper"
	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	nftsim "github.com/CudoVentures/cudos-node/x/nft/simulation"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgMintNft(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	nk types.NftKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		// Publish collection

		ownerAddr, denom, _ := nftsim.GetRandomNFTFromOwner(ctx, nk, r)
		if ownerAddr.Empty() {
			err := fmt.Errorf("invalid account")
			return simtypes.NoOpMsg(types.ModuleName, types.EventMintNftType, err.Error()), nil, err
		}

		collection := types.NewCollection(denom, "", "", ownerAddr.String(), false)
		_, err := k.PublishCollection(ctx, collection)
		if err != nil {
			err := fmt.Errorf("failed to publish collection")
			return simtypes.NoOpMsg(types.ModuleName, types.EventMintNftType, err.Error()), nil, err
		}

		// Mint NFT

		name := strings.ToLower(simtypes.RandStringOfLength(r, 10))
		uri := strings.ToLower(simtypes.RandStringOfLength(r, 10))
		data := strings.ToLower(simtypes.RandStringOfLength(r, 10))

		recipientAcc, _ := simtypes.RandomAcc(r, accs)

		msg := types.NewMsgMintNft(
			ownerAddr.String(),
			denom,
			recipientAcc.Address.String(),
			"100000000000acudos",
			name,
			uri,
			data,
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
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		if _, _, err = app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventMintNftType, err.Error()), nil, err
		}

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "MintNft simulation not implemented"), nil, nil
	}
}
