package types

import (
	"testing"

	"github.com/CudoVentures/cudos-node/testutil/sample"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

var (
	amount = sdk.NewCoin("cosmos", sdk.OneInt())
)

func TestMsgPlaceBid_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgPlaceBid
		err  error
	}{
		{
			// todo cover all ValidateBasic() errors
			name: "invalid address",
			msg: MsgPlaceBid{
				Bidder: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgPlaceBid{
				Bidder:    sample.AccAddress(),
				AuctionId: 1,
				Amount:    amount,
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
