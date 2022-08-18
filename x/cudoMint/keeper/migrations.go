package keeper

import (
	"fmt"
	"strings"

	paramsV1 "github.com/CudoVentures/cudos-node/x/cudoMint/legacy/v1"
	paramsV2 "github.com/CudoVentures/cudos-node/x/cudoMint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{keeper: keeper}
}

// Migrate1to2 migrates from version 1 to 2.
// change blocksPerDay to IncrementModifier
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	bpdByte := m.keeper.paramSpace.GetRaw(ctx, paramsV1.BlocksPerDay)

	bpdStr := strings.ReplaceAll(string(bpdByte), "\"", "")

	bpd, ok := sdk.NewIntFromString(bpdStr)

	if ok {
		newParams := paramsV2.Params{IncrementModifier: bpd}

		m.keeper.paramSpace.SetParamSet(ctx, &newParams)

		return nil
	}

	return fmt.Errorf("failed to parse. blocksperday = %v", bpdStr)
}
