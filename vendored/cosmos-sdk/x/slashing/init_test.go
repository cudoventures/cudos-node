package slashing_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	// The default power validators are initialized to have within tests
	InitTokens = sdk.TokensFromConsensusPower(4000000, sdk.DefaultPowerReduction)
)
