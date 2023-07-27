package simulation

import (
	"fmt"
	"math/rand"
	"strconv"

	simappparams "cosmossdk.io/simapp/params"
	"github.com/CudoVentures/cudos-node/x/addressbook/keeper"
	"github.com/CudoVentures/cudos-node/x/addressbook/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func SimulateMsgCreateAddress(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)

		msg := &types.MsgCreateAddress{
			Creator: simAccount.Address.String(),
			Network: "ETH",
			Value:   "ETH_Value",
			Label:   fmt.Sprintf("%d@mydenom", r.Int()),
		}

		if _, found := k.GetAddress(ctx, msg.Creator, msg.Network, msg.Label); found {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "Address already exist"), nil, nil
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: sdk.NewCoins(),
			AccountKeeper:   ak,
			Bankkeeper:      bk,
		}
		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

func SimulateMsgUpdateAddress(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		var (
			simAccount = simtypes.Account{}
			address    = types.Address{}
			msg        = &types.MsgUpdateAddress{}
			allAddress = k.GetAllAddress(ctx)
			found      = false
		)
		for _, obj := range allAddress {
			simAccount, found = FindAccount(accs, obj.Creator)
			if found {
				address = obj
				break
			}
		}
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "address creator not found"), nil, nil
		}

		msg.Creator = simAccount.Address.String()
		msg.Network = address.Network
		msg.Label = address.Label
		msg.Value = "new value"

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: sdk.NewCoins(),
			AccountKeeper:   ak,
			Bankkeeper:      bk,
		}
		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

func SimulateMsgDeleteAddress(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		var (
			simAccount = simtypes.Account{}
			address    = types.Address{}
			msg        = &types.MsgDeleteAddress{}
			allAddress = k.GetAllAddress(ctx)
			found      = false
		)
		for _, obj := range allAddress {
			simAccount, found = FindAccount(accs, obj.Creator)
			if found {
				address = obj
				break
			}
		}
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "address creator not found"), nil, nil
		}

		msg.Creator = simAccount.Address.String()
		msg.Network = address.Network
		msg.Label = address.Label

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: sdk.NewCoins(),
			AccountKeeper:   ak,
			Bankkeeper:      bk,
		}
		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}
