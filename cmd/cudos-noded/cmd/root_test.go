package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/CudoVentures/cudos-node/app"
	"github.com/CudoVentures/cudos-node/cmd/cudos-noded/cmd"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func TestRootCmdConfig(t *testing.T) {
	rootCmd, _ := cmd.NewRootCmd()
	rootCmd.SetArgs([]string{
		"config",          // Test the config cmd
		"keyring-backend", // key
		"test",            // value
	})
	err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome)
	require.NoError(t, err)
}
