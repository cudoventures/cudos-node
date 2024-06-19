package messaging

import (
	"fmt"

	// "cosmossdk.io/api/tendermint/abci"
	"github.com/CudoVentures/cudos-node/x/messaging/keeper"
	"github.com/CudoVentures/cudos-node/x/messaging/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewHandler ...
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgSendMessage:
			return handleMsgSendMessage(ctx, k, msg)
		// this line is used by starport scaffolding # 1
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgSendMessage(ctx sdk.Context, k keeper.Keeper, msg *types.MsgSendMessage) (*sdk.Result, error) {
	err := k.SendMessage(ctx, *msg)
	if err != nil {
		return nil, err
	}

	// Convert SDK events to Tendermint ABCI events using the helper function
	events := ConvertSDKEventsToABCIEvents(ctx.EventManager().Events())

	return &sdk.Result{
		Events: events,
	}, nil
}

// ConvertSDKEventsToABCIEvents converts SDK events to Tendermint ABCI events
func ConvertSDKEventsToABCIEvents(events sdk.Events) []abci.Event {
	abciEvents := make([]abci.Event, len(events))
	for i, event := range events {
		abciAttributes := make([]abci.EventAttribute, len(event.Attributes))
		for j, attr := range event.Attributes {
			abciAttributes[j] = abci.EventAttribute{
				Key:   []byte(attr.Key),
				Value: []byte(attr.Value),
			}
		}
		abciEvents[i] = abci.Event{
			Type:       event.Type,
			Attributes: abciAttributes,
		}
	}
	return abciEvents
}
