package types

import (
	"testing"

	"github.com/CudoVentures/cudos-node/testutil/sample"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgPlaceBid_ValidateBasic(t *testing.T) {
	for _, tc := range []struct {
		desc    string
		arrange func(msg *MsgPlaceBid)
		wantErr error
	}{
		{
			desc:    "valid",
			arrange: func(msg *MsgPlaceBid) {},
		},
		{
			desc: "zero amount",
			arrange: func(msg *MsgPlaceBid) {
				msg.Amount = sdk.NewCoin("acudos", sdk.ZeroInt())
			},
			wantErr: ErrInvalidPrice,
		},
		{
			desc: "invalid amount denom",
			arrange: func(msg *MsgPlaceBid) {
				msg.Amount = sdk.Coin{Denom: "", Amount: sdk.ZeroInt()}

			},
			wantErr: ErrInvalidPrice,
		},
		{
			desc: "invalid address",
			arrange: func(msg *MsgPlaceBid) {
				msg.Bidder = "invalid"
			},
			wantErr: sdkerrors.ErrInvalidAddress,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			msg := NewMsgPlaceBid(
				sample.AccAddress(),
				0,
				sdk.NewCoin("acudos", sdk.OneInt()),
			)

			tc.arrange(msg)

			err := msg.ValidateBasic()
			require.ErrorIs(t, err, tc.wantErr)
		})
	}
}
