package types

import (
	"testing"

	"github.com/CudoVentures/cudos-node/testutil/sample"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func TestMsgMintNft_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgMintNft
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgMintNft{
				DenomId: "abc",
				Name:    "abc",
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgMintNft{
				DenomId:   "abc",
				Name:      "abc",
				Price:     sdk.NewCoin("acudos", sdk.NewInt(1000)),
				Creator:   sample.AccAddress(),
				Recipient: sample.AccAddress(),
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
