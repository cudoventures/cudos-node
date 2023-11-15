package types

import (
	"testing"

	"github.com/CudoVentures/cudos-node/testutil/sample"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func TestMsgUpdatePrice_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdatePrice
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdatePrice{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdatePrice{
				Creator: sample.AccAddress(),
				Price:   sdk.NewCoin("acudos", sdk.NewInt(1000)),
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
