package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	customModuleErrorPrefix uint32 = 100000
	ErrUnknownCollection           = sdkerrors.Register(ModuleName, customModuleErrorPrefix+3, "unknown nft collection")
	ErrInvalidNFT                  = sdkerrors.Register(ModuleName, customModuleErrorPrefix+4, "invalid nft")
	ErrNFTAlreadyExists            = sdkerrors.Register(ModuleName, customModuleErrorPrefix+5, "nft already exists")
	ErrUnknownNFT                  = sdkerrors.Register(ModuleName, customModuleErrorPrefix+6, "unknown nft")
	ErrUnauthorized                = sdkerrors.Register(ModuleName, customModuleErrorPrefix+8, "unauthorized address")
	ErrInvalidDenom                = sdkerrors.Register(ModuleName, customModuleErrorPrefix+9, "invalid denom")
	ErrInvalidTokenID              = sdkerrors.Register(ModuleName, customModuleErrorPrefix+10, "invalid nft id")
	ErrInvalidTokenURI             = sdkerrors.Register(ModuleName, customModuleErrorPrefix+11, "invalid nft uri")
	ErrInvalidDenomName            = sdkerrors.Register(ModuleName, customModuleErrorPrefix+12, "invalid denom name")
	ErrNoApprovedAddresses         = sdkerrors.Register(ModuleName, customModuleErrorPrefix+13, "no approved addresses!")
	ErrNotFoundNFT                 = sdkerrors.Register(ModuleName, customModuleErrorPrefix+14, "nft not found")
	ErrInvalidDenomSymbol          = sdkerrors.Register(ModuleName, customModuleErrorPrefix+15, "invalid denom symbol")
	ErrInvalidNftName              = sdkerrors.Register(ModuleName, customModuleErrorPrefix+16, "invalid nft name")
	ErrInvalidTraits               = sdkerrors.Register(ModuleName, customModuleErrorPrefix+17, "invalid traits")
	ErrAlreadySoftLocked           = sdkerrors.Register(ModuleName, customModuleErrorPrefix+18, "already soft locked")
	ErrNotSoftLocked               = sdkerrors.Register(ModuleName, customModuleErrorPrefix+19, "not soft locked")
	ErrNotOwnerOfSoftLock          = sdkerrors.Register(ModuleName, customModuleErrorPrefix+20, "not owner of soft lock")
	ErrSoftLocked                  = sdkerrors.Register(ModuleName, customModuleErrorPrefix+21, "soft locked")
	ErrNotEditable                 = sdkerrors.Register(ModuleName, customModuleErrorPrefix+22, "not editable")
	ErrInvalidDescription          = sdkerrors.Register(ModuleName, customModuleErrorPrefix+23, "invalid description")
	ErrInvalidTokenData            = sdkerrors.Register(ModuleName, customModuleErrorPrefix+24, "invalid token data")
	ErrInvalidDenomData            = sdkerrors.Register(ModuleName, customModuleErrorPrefix+25, "invalid denom data")
	ErrInvalidSchema               = sdkerrors.Register(ModuleName, customModuleErrorPrefix+26, "invalid schema")
)
