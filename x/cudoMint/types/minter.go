package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewMinter returns a new Minter object with the given inflation and annual
// provisions values.
func NewMinter(mintRemainder sdk.Dec, normTimePassed sdk.Dec) Minter {
	return Minter{
		MintRemainder:  mintRemainder,
		NormTimePassed: normTimePassed,
	}
}

func DefaultInitialMinter() Minter {
	return Minter{
		MintRemainder:  sdk.ZeroDec(),
		NormTimePassed: sdk.ZeroDec(),
	}
}

// ValidateMinter validate minter
func ValidateMinter(minter Minter) error {
	if minter.MintRemainder.IsNegative() {
		return fmt.Errorf("mint parameter MintRemainder should be positive, is %s",
			minter.MintRemainder.String())
	} else if minter.NormTimePassed.IsNegative() {
		return fmt.Errorf("mint parameter NormTimePassed should be positive, is %s",
			minter.MintRemainder.String())
	}

	return nil
}
