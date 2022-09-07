package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/addressbook module sentinel errors
var (
	ErrInvalidNetwork   = sdkerrors.Register(ModuleName, 1100, "invalid network")
	ErrInvalidLabel     = sdkerrors.Register(ModuleName, 1101, "invalid label")
	ErrKeyAlreadyExists = sdkerrors.Register(ModuleName, 1102, "key already exist")
	ErrKeyNotFound      = sdkerrors.Register(ModuleName, 1103, "key not found")
	ErrInvalidValue     = sdkerrors.Register(ModuleName, 1104, "invalid value")
)
