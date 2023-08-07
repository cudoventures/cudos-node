package custom_bindings

import (
	"encoding/json"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	"github.com/CudoVentures/cudos-node/x/addressbook/keeper"
	"github.com/CudoVentures/cudos-node/x/addressbook/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ wasmkeeper.Messenger = (*AddressbookWasmMessengerDecorator)(nil)

type AddressbookWasmMessengerDecorator struct {
	old    wasmkeeper.Messenger
	msgSrv types.MsgServer
}

func NewAddressbookWasmMessageDecorator(addressbookKeeper keeper.Keeper) func(wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &AddressbookWasmMessengerDecorator{
			old:    old,
			msgSrv: keeper.NewMsgServerImpl(addressbookKeeper),
		}
	}
}

func (m *AddressbookWasmMessengerDecorator) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) (events []sdk.Event, data [][]byte, err error) {
	if msg.Custom == nil {
		return m.old.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
	}

	var customMsg addressbookCustomMsg
	if err := json.Unmarshal(msg.Custom, &customMsg); err != nil {
		return nil, nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	switch {
	case customMsg.CreateAddress != nil:
		c := customMsg.CreateAddress
		msg := types.NewMsgCreateAddress(c.Creator, c.Network, c.Label, c.Value)
		if err = msg.ValidateBasic(); err != nil {
			return nil, nil, err
		}
		_, err := m.msgSrv.CreateAddress(ctx, msg)
		return nil, nil, err

	case customMsg.UpdateAddress != nil:
		c := customMsg.UpdateAddress
		msg := types.NewMsgUpdateAddress(c.Creator, c.Network, c.Label, c.Value)
		if err = msg.ValidateBasic(); err != nil {
			return nil, nil, err
		}
		_, err := m.msgSrv.UpdateAddress(ctx, msg)
		return nil, nil, err
	case customMsg.DeleteAddress != nil:
		c := customMsg.DeleteAddress
		msg := types.NewMsgDeleteAddress(c.Creator, c.Network, c.Label)
		if err = msg.ValidateBasic(); err != nil {
			return nil, nil, err
		}
		_, err := m.msgSrv.DeleteAddress(ctx, msg)
		return nil, nil, err
	default:
		return m.old.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
	}
}

type addressbookCustomMsg struct {
	CreateAddress *CreateAddressRequest `json:"create_address_msg,omitempty"`
	UpdateAddress *UpdateAddressRequest `json:"update_address_msg,omitempty"`
	DeleteAddress *DeleteAddressRequest `json:"delete_address_msg,omitempty"`
}

type DeleteAddressRequest struct {
	Creator string `json:"creator"`
	Network string `json:"network"`
	Label   string `json:"label"`
}

type CreateAddressRequest struct {
	DeleteAddressRequest
	Value string `json:"value"`
}

type UpdateAddressRequest struct {
	CreateAddressRequest
}
