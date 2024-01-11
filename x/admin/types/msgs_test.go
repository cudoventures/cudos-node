package types_test

import (
	"testing"

	"github.com/CudoVentures/cudos-node/x/admin/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestMsgAdminSpendCommunityPoolRoute(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from"))
	addr2 := sdk.AccAddress([]byte("to"))
	coins := sdk.NewCoins(sdk.NewInt64Coin("acudos", 10))
	msg := types.NewMsgAdminSpendCommunityPool(addr1, addr2, coins)

	require.Equal(t, msg.Route(), types.RouterKey)
	require.Equal(t, msg.Type(), "adminSpendCommunityPool")
}

func TestMsgSend_ValidateBasic(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from________________"))
	addr2 := sdk.AccAddress([]byte("to__________________"))
	addrEmpty := sdk.AccAddress([]byte(""))
	addrLong := sdk.AccAddress([]byte("Purposefully long address"))

	acudos123 := sdk.NewCoins(sdk.NewInt64Coin("acudos", 123))
	acudos0 := sdk.NewCoins(sdk.NewInt64Coin("acudos", 0))
	acudos123eth123 := sdk.NewCoins(sdk.NewInt64Coin("acudos", 123), sdk.NewInt64Coin("eth", 123))
	acudos123eth0 := sdk.Coins{sdk.NewInt64Coin("acudos", 123), sdk.NewInt64Coin("eth", 0)}

	cases := []struct {
		expectedErr string // empty means no error expected
		msg         *types.MsgAdminSpendCommunityPool
	}{
		{"", types.NewMsgAdminSpendCommunityPool(addr1, addr2, acudos123)},                                  // valid send
		{"", types.NewMsgAdminSpendCommunityPool(addr1, addr2, acudos123eth123)},                            // valid send with multiple coins
		{"", types.NewMsgAdminSpendCommunityPool(addrLong, addr2, acudos123)},                               // valid send with long addr sender
		{"", types.NewMsgAdminSpendCommunityPool(addr1, addrLong, acudos123)},                               // valid send with long addr recipient
		{": invalid coins", types.NewMsgAdminSpendCommunityPool(addr1, addr2, acudos0)},                     // non positive coin
		{"123acudos,0eth: invalid coins", types.NewMsgAdminSpendCommunityPool(addr1, addr2, acudos123eth0)}, // non positive coin in multicoins
		{"Invalid sender address (empty address string is not allowed): invalid address", types.NewMsgAdminSpendCommunityPool(addrEmpty, addr2, acudos123)},
		{"Invalid recipient address (empty address string is not allowed): invalid address", types.NewMsgAdminSpendCommunityPool(addr1, addrEmpty, acudos123)},
	}

	for _, tc := range cases {
		err := tc.msg.ValidateBasic()
		if tc.expectedErr == "" {
			require.Nil(t, err)
		} else {
			require.EqualError(t, err, tc.expectedErr)
		}
	}
}

func TestMsgAdminSpendCommunityPoolGetSignBytes(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("input"))
	addr2 := sdk.AccAddress([]byte("output"))
	coins := sdk.NewCoins(sdk.NewInt64Coin("acudos", 10))
	msg := types.NewMsgAdminSpendCommunityPool(addr1, addr2, coins)
	res := msg.GetSignBytes()

	expected := `{"coins":[{"amount":"10","denom":"acudos"}],"initiator":"cosmos1d9h8qat57ljhcm","to_address":"cosmos1da6hgur4wsmpnjyg"}`
	require.Equal(t, expected, string(res))
}

func TestMsgAdminSpendCommunityPoolGetSigners(t *testing.T) {
	from := sdk.AccAddress([]byte("input111111111111111"))
	msg := types.NewMsgAdminSpendCommunityPool(from, sdk.AccAddress{}, sdk.NewCoins())
	res := msg.GetSigners()
	require.Equal(t, 1, len(res))
	require.True(t, from.Equals(res[0]))
}
