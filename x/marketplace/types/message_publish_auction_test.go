package types

import (
	"testing"
	"time"

	"github.com/CudoVentures/cudos-node/testutil/sample"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgPublishAuction_ValidateBasic(t *testing.T) {
	tests := []struct {
		desc    string
		arrange func(msg *MsgPublishAuction)
		errMsg  string
	}{
		{
			desc:    "valid",
			arrange: func(msg *MsgPublishAuction) {},
		},
		{
			desc: "invalid address",
			arrange: func(msg *MsgPublishAuction) {
				msg.Creator = "invalid"
			},
			errMsg: "",
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			msg, err := NewMsgPublishAuction(
				sample.AccAddress(),
				"test",
				"123",
				time.Hour*24,
				&EnglishAuction{MinPrice: sdk.NewCoin("acudos", sdk.OneInt())},
			)
			require.NoError(t, err)

			tc.arrange(msg)

			if tc.errMsg != "" {
				require.EqualError(t, err, tc.errMsg)
				return
			}
			require.NoError(t, err)
		})
	}
}
