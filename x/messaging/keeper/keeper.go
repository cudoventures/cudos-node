// Package keeper implements the keeper functions for the messaging module
// in a Cosmos SDK application. The keeper is responsible for managing the
// state and providing the core functionality for the module.
package keeper

import (
	"encoding/json"
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/CudoVentures/cudos-node/x/messaging/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

// DefaultFee represents the default fee for sending a message.
// TODO: Make this value configurable.
const DefaultFee int64 = 1000

// Keeper defines the messaging module's keeper. It holds references to the
// necessary codec, store keys, parameter subspace, and other keepers.
type Keeper struct {
	cdc           codec.BinaryCodec
	storeKey      sdk.StoreKey
	memKey        sdk.StoreKey
	paramstore    paramtypes.Subspace
	accountKeeper types.AccountKeeper
	BankKeeper    types.BankKeeper
}

// BankKeeperAdapter is an adapter for the bank keeper to add additional methods if needed.
type BankKeeperAdapter struct {
	keeper.Keeper
}

// GetModuleAddress returns the address of a module account given its name.
func (bka *BankKeeperAdapter) GetModuleAddress(name string) sdk.AccAddress {
	// Implement this method if it's not already in the Keeper
	// This might involve accessing another component that can resolve module addresses
	return sdk.AccAddress{} // Return the correct address
}

// NewKeeper creates a new Keeper instance.
func NewKeeper(
	cdc codec.BinaryCodec,
	amino *codec.LegacyAmino,
	storeKey, memKey sdk.StoreKey,
	ps paramtypes.Subspace,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
) *Keeper {
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		memKey:        memKey,
		paramstore:    ps,
		accountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
	}
}

// Logger returns a logger for the messaging module.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// AppendMessage appends a message to the store.
func (k Keeper) AppendMessage(ctx sdk.Context, message *types.Message) error {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(message)
	store.Set([]byte(message.Subject), bz)
	return nil
}

// NewQuerier creates a new querier for the messaging module.
func (k Keeper) NewQuerier() sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case "getMessages":
			return getMessages(ctx, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown query path")
		}
	}
}

// GetAllMessages retrieves all messages from the store.
func (k Keeper) GetAllMessages(ctx sdk.Context) ([]*types.Message, error) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte("MessagePrefix"))
	defer iterator.Close()

	var messages []*types.Message
	for ; iterator.Valid(); iterator.Next() {
		var msg types.Message
		k.cdc.MustUnmarshal(iterator.Value(), &msg)
		messages = append(messages, &msg)
	}

	return messages, nil
}

// getMessages retrieves and marshals all messages for querying.
func getMessages(ctx sdk.Context, k Keeper) ([]byte, error) {
	messages, err := k.GetAllMessages(ctx)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := json.MarshalIndent(messages, "", "  ")
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

// SendMessage handles sending a message, including fee deduction and permission checks.
func (k Keeper) SendMessage(ctx sdk.Context, msg types.MsgSendMessage) error {
	if msg.FromAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "sender address cannot be empty")
	}

	if !k.IsSenderAuthorized(ctx, msg.FromAddress) {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "sender is not authorized to send messages")
	}

	fee := k.CalculateFeeForMessage(msg)
	if !k.HasSufficientCoins(ctx, msg.FromAddress, sdk.NewCoins(fee)) {
		return sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "insufficient funds")
	}

	err := k.BankKeeper.SendCoinsFromAccountToModule(ctx, msg.FromAddress, "msgFees", sdk.NewCoins(fee))
	if err != nil {
		return sdkerrors.Wrap(err, "failed to deduct fees for message sending")
	}

	return k.StoreMessage(ctx, msg)
}

// IsSenderAuthorized checks if the sender is authorized to send messages.
func (k Keeper) IsSenderAuthorized(ctx sdk.Context, sender sdk.AccAddress) bool {
	account := k.accountKeeper.GetAccount(ctx, sender)
	if account == nil {
		ctx.Logger().Error("Sender account does not exist", "address", sender.String())
		return false
	}

	if acc, ok := account.(authtypes.AccountI); ok {
		permissions := k.GetAccountPermissions(ctx, acc)
		if !permissions["canSendMessages"] {
			ctx.Logger().Error("Sender does not have permission to send messages", "address", sender.String())
			return false
		}
	} else {
		ctx.Logger().Error("Failed to cast account to AccountI", "address", sender.String())
		return false
	}

	return true
}

// GetAccountPermissions retrieves permissions for the given account.
func (k Keeper) GetAccountPermissions(ctx sdk.Context, account authtypes.AccountI) map[string]bool {
	return map[string]bool{
		"canSendMessages": true, // TODO: Replace with actual logic
	}
}

// CalculateFeeForMessage calculates the fee for sending a message.
func (k Keeper) CalculateFeeForMessage(msg types.MsgSendMessage) sdk.Coin {
	feeAmount := DefaultFee
	return sdk.NewCoin("yourToken", sdk.NewInt(int64(feeAmount)))
}

// StoreMessage stores the message in the KVStore.
func (k Keeper) StoreMessage(ctx sdk.Context, msg types.MsgSendMessage) error {
	message := types.Message{
		FromAddress: msg.FromAddress.String(),
		ToAddress:   msg.ToAddress.String(),
		Subject:     msg.Subject,
		Body:        msg.Body,
	}

	bz, err := k.cdc.Marshal(&message)
	if err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	messageStore := prefix.NewStore(store, []byte("Message"))

	messageKey := []byte(msg.Subject)
	messageStore.Set(messageKey, bz)

	return nil
}

// GetModuleAddress returns the address of the module account.
func (k Keeper) GetModuleAddress(moduleName string) sdk.AccAddress {
	return k.accountKeeper.GetModuleAddress(moduleName)
}

// HasSufficientCoins checks if the given address has sufficient coins.
func (k Keeper) HasSufficientCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool {
	spendableCoins := k.BankKeeper.SpendableCoins(ctx, addr)
	return spendableCoins.IsAllGTE(amt)
}

// SendCoins sends coins from one address to another.
func (k Keeper) SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
	if !k.HasSufficientCoins(ctx, fromAddr, amt) {
		return sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "insufficient funds")
	}

	return k.BankKeeper.SendCoins(ctx, fromAddr, toAddr, amt)
}
