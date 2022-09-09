package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrUnknownCollection   = sdkerrors.Register(ModuleName, 3, "unknown nft collection")
	ErrInvalidNFT          = sdkerrors.Register(ModuleName, 4, "invalid nft")
	ErrNFTAlreadyExists    = sdkerrors.Register(ModuleName, 5, "nft already exists")
	ErrUnknownNFT          = sdkerrors.Register(ModuleName, 6, "unknown nft")
	ErrUnauthorized        = sdkerrors.Register(ModuleName, 8, "unauthorized address")
	ErrInvalidDenom        = sdkerrors.Register(ModuleName, 9, "invalid denom")
	ErrInvalidTokenID      = sdkerrors.Register(ModuleName, 10, "invalid nft id")
	ErrInvalidTokenURI     = sdkerrors.Register(ModuleName, 11, "invalid nft uri")
	ErrInvalidDenomName    = sdkerrors.Register(ModuleName, 12, "invalid denom name")
	ErrNoApprovedAddresses = sdkerrors.Register(ModuleName, 13, "no approved addresses!")
	ErrNotFoundNFT         = sdkerrors.Register(ModuleName, 14, "nft not found")
	ErrInvalidDenomSymbol  = sdkerrors.Register(ModuleName, 15, "invalid denom symbol")
	ErrInvalidNftName      = sdkerrors.Register(ModuleName, 16, "invalid nft name")
	ErrInvalidTraits       = sdkerrors.Register(ModuleName, 17, "invalid traits")
	ErrAlreadySoftLocked   = sdkerrors.Register(ModuleName, 18, "already soft locked")
	ErrNotSoftLocked       = sdkerrors.Register(ModuleName, 19, "not soft locked")
	ErrNotOwnerOfSoftLock  = sdkerrors.Register(ModuleName, 20, "not owner of soft lock")
	ErrSoftLocked          = sdkerrors.Register(ModuleName, 21, "soft locked")
	ErrNotEditable         = sdkerrors.Register(ModuleName, 22, "not editable")
	ErrInvalidDescription  = sdkerrors.Register(ModuleName, 23, "invalid description")
	ErrInvalidTokenData    = sdkerrors.Register(ModuleName, 24, "invalid token data")
	ErrInvalidDenomData    = sdkerrors.Register(ModuleName, 25, "invalid denom data")
)
