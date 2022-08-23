package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/token module sentinel errors
var (
	// todo add error and throw it instead of errors.New
	ErrSample = sdkerrors.Register(ModuleName, 1100, "sample error")
)
