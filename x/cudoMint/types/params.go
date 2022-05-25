package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// IncrementModifier Parameter store keys
var (
	IncrementModifier = []byte("IncrementModifier")
)

// ParamKeyTable ParamTable for minting module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(
	incrementModifier sdk.Int,
) Params {

	return Params{
		IncrementModifier: incrementModifier,
	}
}

// DefaultParams default minting module parameters
func DefaultParams() Params {
	return Params{
		IncrementModifier: sdk.NewInt(17280), // assuming 5 second block times
	}
}

// Validate validate params
func (p Params) Validate() error {
	if err := validateIncrementModifier(p.IncrementModifier); err != nil {
		return err
	}

	return nil

}

// ParamSetPairs Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(IncrementModifier, &p.IncrementModifier, validateIncrementModifier),
	}
}

func validateIncrementModifier(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() || v.IsZero() {
		return fmt.Errorf("blocks per day must be positive: %s", v)
	}
	return nil
}
