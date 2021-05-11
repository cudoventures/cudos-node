package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// NewMinter returns a new Minter object with the given inflation and annual
// provisions values.
func NewMinter(mintRemainder sdk.Dec) Minter {
	return Minter{
		MintRemainder: mintRemainder,
	}
}
