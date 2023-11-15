package types

import (
	"testing"

	"github.com/CudoVentures/cudos-node/testutil/sample"
	"github.com/stretchr/testify/require"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func TestMsgRemoveAdmin_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgRemoveAdmin
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgRemoveAdmin{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgRemoveAdmin{
				Creator: sample.AccAddress(),
				Address: sample.AccAddress(),
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
