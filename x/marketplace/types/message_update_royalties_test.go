package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/CudoVentures/cudos-node/testutil/sample"
)

func TestMsgUpdateRoyalties_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateRoyalties
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateRoyalties{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateRoyalties{
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
