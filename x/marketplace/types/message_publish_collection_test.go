package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/CudoVentures/cudos-node/testutil/sample"
)

func TestMsgPublishCollection_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgPublishCollection
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgPublishCollection{
				DenomId: "abc",
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgPublishCollection{
				DenomId: "abc",
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
