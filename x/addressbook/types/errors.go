package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/addressbook module sentinel errors
var (
	ErrInvalidNetwork = sdkerrors.Register(ModuleName, 1100, "invalid network")
	ErrInvalidLabel   = sdkerrors.Register(ModuleName, 1101, "invalid label")
)
