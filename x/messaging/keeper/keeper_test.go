package keeper_test

import (
	"testing"

	keepertest "github.com/CudoVentures/cudos-node/testutil/keeper"
	"github.com/CudoVentures/cudos-node/x/messaging/keeper"
	"github.com/CudoVentures/cudos-node/x/messaging/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func setupKeeper(t *testing.T) (keeper.Keeper, sdk.Context) {
	k, ctx := keepertest.MessagingKeeper(t)
	return *k, ctx
}

func TestSendMessage(t *testing.T) {
	k, ctx := setupKeeper(t)
	sender := sdk.AccAddress([]byte("sender"))
	receiver := sdk.AccAddress([]byte("receiver"))
	msg := types.MsgSendMessage{
		FromAddress: sender,
		ToAddress:   receiver,
		Subject:     "Test Message",
		Body:        "Hello, Cosmos!",
	}

	// Try sending the message
	err := k.SendMessage(ctx, msg)
	require.NoError(t, err)

	// Check the new balance to confirm fee deduction
	balance := k.BankKeeper.GetBalance(ctx, sender, "token")
	require.Equal(t, sdk.NewInt(9000), balance.Amount) // Assuming a fee of 1000
}

func TestPermissionDenied(t *testing.T) {
	k, ctx := setupKeeper(t)
	sender := sdk.AccAddress([]byte("unauthorizedUser"))
	msg := types.MsgSendMessage{
		FromAddress: sender,
		ToAddress:   sdk.AccAddress([]byte("receiver")),
		Subject:     "Unauthorized Test",
		Body:        "Should fail",
	}

	err := k.SendMessage(ctx, msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "sender is not authorized")
}

func TestInsufficientFunds(t *testing.T) {
	k, ctx := setupKeeper(t)
	sender := sdk.AccAddress([]byte("poorUser"))
	msg := types.MsgSendMessage{
		FromAddress: sender,
		ToAddress:   sdk.AccAddress([]byte("receiver")),
		Subject:     "No Money Test",
		Body:        "I have no coins",
	}

	err := k.SendMessage(ctx, msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "insufficient funds")
}
