package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/CudoVentures/cudos-node/testutil/sample"
)

func TestMsgCreateCollection_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreateCollection
		err  error
	}{
		{
			name: "valid",
			msg: MsgCreateCollection{
				Id:      "testdenom",
				Name:    "testname",
				Symbol:  "testsymbol",
				Schema:  "testschema",
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
