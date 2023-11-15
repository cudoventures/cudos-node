package types

import (
	"testing"

	"github.com/CudoVentures/cudos-node/testutil/sample"
	"github.com/stretchr/testify/require"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func TestMsgCreateAddress_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreateAddress
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreateAddress{
				Network: "n",
				Label:   "l",
				Value:   "v",
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCreateAddress{
				Network: "n",
				Label:   "l",
				Value:   "v",
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgUpdateAddress_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateAddress
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateAddress{
				Network: "n",
				Label:   "l",
				Value:   "v",
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateAddress{
				Network: "n",
				Label:   "l",
				Value:   "v",
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgDeleteAddress_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgDeleteAddress
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgDeleteAddress{
				Network: "n",
				Label:   "l",
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgDeleteAddress{
				Network: "n",
				Label:   "l",
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
