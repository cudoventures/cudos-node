package types

import (
	"testing"

	"github.com/CudoVentures/cudos-node/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgAcceptBid_ValidateBasic(t *testing.T) {
	for _, tc := range []struct {
		desc string
		msg  MsgAcceptBid
		err  error
	}{
		{
			desc: "valid",
			msg: MsgAcceptBid{
				Sender: sample.AccAddress(),
			},
		},
		{
			desc: "invalid address",
			msg: MsgAcceptBid{
				Sender: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
