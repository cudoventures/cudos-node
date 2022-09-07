package simulation

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/CudoVentures/cudos-node/x/nft/keeper"
	"github.com/CudoVentures/cudos-node/x/nft/types"
)

// Simulation operation weights constants
const (
	OpWeightMsgIssueDenom    = "op_weight_msg_issue_denom"
	OpWeightMsgMintNFT       = "op_weight_msg_mint_nft"
	OpWeightMsgEditNFT       = "op_weight_msg_edit_nft_tokenData"
	OpWeightMsgTransferNFT   = "op_weight_msg_transfer_nft"
	OpWeightMsgBurnNFT       = "op_weight_msg_transfer_burn_nft"
	OpWeightMsgTransferDenom = "op_weight_msg_transfer_denom"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams,
	cdc codec.JSONCodec,
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
) simulation.WeightedOperations {
	var weightIssueDenom, weightMint, weightEdit, weightBurn, weightTransfer, weightTransferDenom int

	appParams.GetOrGenerate(
		cdc, OpWeightMsgIssueDenom, &weightIssueDenom, nil,
		func(_ *rand.Rand) {
			weightIssueDenom = 50
		},
	)

	appParams.GetOrGenerate(
		cdc, OpWeightMsgMintNFT, &weightMint, nil,
		func(_ *rand.Rand) {
			weightMint = 100
		},
	)

	appParams.GetOrGenerate(
		cdc, OpWeightMsgEditNFT, &weightEdit, nil,
		func(_ *rand.Rand) {
			weightEdit = 50
		},
	)

	appParams.GetOrGenerate(
		cdc, OpWeightMsgTransferNFT, &weightTransfer, nil,
		func(_ *rand.Rand) {
			weightTransfer = 50
		},
	)

	appParams.GetOrGenerate(
		cdc, OpWeightMsgBurnNFT, &weightBurn, nil,
		func(_ *rand.Rand) {
			weightBurn = 10
		},
	)
	appParams.GetOrGenerate(
		cdc, OpWeightMsgTransferDenom, &weightTransferDenom, nil,
		func(_ *rand.Rand) {
			weightTransferDenom = 10
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightIssueDenom,
			SimulateMsgIssueDenom(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightMint,
			SimulateMsgMintNFT(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightEdit,
			SimulateMsgEditNFT(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightTransfer,
			SimulateMsgTransferNFT(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightBurn,
			SimulateMsgBurnNFT(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightTransferDenom,
			SimulateMsgTransferDenom(k, ak, bk),
		),
	}
}

// SimulateMsgTransferNFT simulates the transfer of an NFT
func SimulateMsgTransferNFT(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		opMsg simtypes.OperationMsg, fOps []simtypes.FutureOperation, err error,
	) {
		ownerAddr, denom, nftID := GetRandomNFTFromOwner(ctx, k, r)
		if ownerAddr.Empty() {
			err = fmt.Errorf("invalid account")
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeTransferNft, err.Error()), nil, err
		}

		recipientAccount, _ := simtypes.RandomAcc(r, accs)
		msg := types.NewMsgTransferNft(
			nftID,
			denom,
			ownerAddr.String(),
			recipientAccount.Address.String(),
			ownerAddr.String(),
			"",
		)

		account := ak.GetAccount(ctx, ownerAddr)

		ownerAccount, found := simtypes.FindAccount(accs, ownerAddr)
		if !found {
			err = fmt.Errorf("account %s not found", msg.Sender)
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeTransferNft, err.Error()), nil, err
		}

		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeTransferNft, err.Error()), nil, err
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
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeTransferNft, err.Error()), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// SimulateMsgEditNFT simulates an edit tokenData transaction
func SimulateMsgEditNFT(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		opMsg simtypes.OperationMsg, fOps []simtypes.FutureOperation, err error,
	) {
		ownerAddr, denom, nftID := GetRandomNFTFromOwner(ctx, k, r)
		if ownerAddr.Empty() {
			err = fmt.Errorf("account invalid")
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeEditNFT, err.Error()), nil, err
		}

		msg := types.NewMsgEditNFT(
			nftID,
			denom,
			"",
			simtypes.RandStringOfLength(r, 45), // tokenURI
			simtypes.RandStringOfLength(r, 10), // tokenData
			ownerAddr.String(),
			"",
		)

		account := ak.GetAccount(ctx, ownerAddr)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeEditNFT, err.Error()), nil, err
		}

		ownerAccount, found := simtypes.FindAccount(accs, ownerAddr)
		if !found {
			err = fmt.Errorf("account %s not found", ownerAddr)
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeEditNFT, err.Error()), nil, err
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
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeEditNFT, err.Error()), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// SimulateMsgMintNFT simulates a mint of an NFT
func SimulateMsgMintNFT(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		opMsg simtypes.OperationMsg, fOps []simtypes.FutureOperation, err error,
	) {
		randomSender, _ := simtypes.RandomAcc(r, accs)
		randomRecipient, _ := simtypes.RandomAcc(r, accs)

		msg := types.NewMsgMintNFT(
			getRandomDenom(ctx, k, r), // denom
			"",
			simtypes.RandStringOfLength(r, 45), // tokenURI
			simtypes.RandStringOfLength(r, 10), // tokenData
			randomSender.Address.String(),      // sender
			randomRecipient.Address.String(),   // recipient
			"",
		)

		account := ak.GetAccount(ctx, randomSender.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeMintNFT, err.Error()), nil, err
		}

		simAccount, found := simtypes.FindAccount(accs, randomSender.Address)
		if !found {
			err = fmt.Errorf("account %s not found", msg.Sender)
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeMintNFT, err.Error()), nil, err
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
			simAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		if _, _, err = app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeMintNFT, err.Error()), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// SimulateMsgBurnNFT simulates a burn of an existing NFT
func SimulateMsgBurnNFT(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		opMsg simtypes.OperationMsg, fOps []simtypes.FutureOperation, err error,
	) {
		ownerAddr, denom, nftID := GetRandomNFTFromOwner(ctx, k, r)
		if ownerAddr.Empty() {
			err = fmt.Errorf("invalid account")
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeBurnNFT, err.Error()), nil, err
		}

		msg := types.NewMsgBurnNFT(ownerAddr.String(), nftID, denom, "")

		account := ak.GetAccount(ctx, ownerAddr)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeBurnNFT, err.Error()), nil, err
		}

		simAccount, found := simtypes.FindAccount(accs, ownerAddr)
		if !found {
			err = fmt.Errorf("account %s not found", msg.Sender)
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeBurnNFT, err.Error()), nil, err
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
			simAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		if _, _, err = app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeEditNFT, err.Error()), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// SimulateMsgTransferDenom simulates the transfer of an denom
func SimulateMsgTransferDenom(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		opMsg simtypes.OperationMsg, fOps []simtypes.FutureOperation, err error,
	) {

		denomId := getRandomDenom(ctx, k, r)
		denom, err := k.GetDenom(ctx, denomId)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgTransferDenom, err.Error()), nil, err
		}

		creator, err := sdk.AccAddressFromBech32(denom.Creator)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgTransferDenom, err.Error()), nil, err
		}
		account := ak.GetAccount(ctx, creator)
		owner, found := simtypes.FindAccount(accs, account.GetAddress())
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgTransferDenom, "creator not found"), nil, nil
		}

		recipient, _ := simtypes.RandomAcc(r, accs)
		msg := types.NewMsgTransferDenom(
			denomId,
			denom.Creator,
			recipient.Address.String(),
			"",
		)

		spendable := bk.SpendableCoins(ctx, owner.Address)
		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgTransferDenom, err.Error()), nil, err
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
			owner.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		if _, _, err = app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgTransferDenom, err.Error()), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// SimulateMsgIssueDenom simulates issue an denom
func SimulateMsgIssueDenom(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		opMsg simtypes.OperationMsg, fOps []simtypes.FutureOperation, err error,
	) {

		denomId := strings.ToLower(simtypes.RandStringOfLength(r, 10))
		denomName := strings.ToLower(simtypes.RandStringOfLength(r, 10))
		denomSymbol := strings.ToLower(simtypes.RandStringOfLength(r, 10))
		sender, _ := simtypes.RandomAcc(r, accs)

		if err := types.ValidateDenomID(denomId); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeIssueDenom, "invalid denom"), nil, nil
		}

		denom, _ := k.GetDenom(ctx, denomId)
		if denom.Size() != 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeIssueDenom, "denom exist"), nil, nil
		}

		msg := types.NewMsgIssueDenom(
			sender.Address.String(),
			denomId,
			denomName,
			"schema",
			"",
			denomSymbol,
			"",
			"",
			"",
		)
		account := ak.GetAccount(ctx, sender.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeIssueDenom, err.Error()), nil, err
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
			sender.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		if _, _, err = app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeIssueDenom, err.Error()), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

func GetRandomNFTFromOwner(ctx sdk.Context, k NftKeeper, r *rand.Rand) (address sdk.AccAddress, denomID, tokenID string) {
	owners, err := k.GetOwners(ctx)
	if err != nil {
		return nil, "", ""
	}

	ownersLen := len(owners)
	if ownersLen == 0 {
		return nil, "", ""
	}

	// get random owner
	i := r.Intn(ownersLen)
	owner := owners[i]

	idCollectionsLen := len(owner.IDCollections)
	if idCollectionsLen == 0 {
		return nil, "", ""
	}

	// get random collection from owner's balance
	i = r.Intn(idCollectionsLen)
	idCollection := owner.IDCollections[i] // nfts IDs
	denomID = idCollection.DenomId

	idsLen := len(idCollection.TokenIds)
	if idsLen == 0 {
		return nil, "", ""
	}

	// get random nft from collection
	i = r.Intn(idsLen)
	tokenID = idCollection.TokenIds[i]

	ownerAddress, _ := sdk.AccAddressFromBech32(owner.Address)
	return ownerAddress, denomID, tokenID
}

func getRandomDenom(ctx sdk.Context, k keeper.Keeper, r *rand.Rand) string {
	var denoms = []string{numbers, letters}
	i := r.Intn(len(denoms))
	return denoms[i]
}

func genRandomBool(r *rand.Rand) bool {
	return r.Int()%2 == 0
}

type NftKeeper interface {
	GetOwners(ctx sdk.Context) (owners types.Owners, err error)
}
