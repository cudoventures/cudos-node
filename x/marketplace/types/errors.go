package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/marketplace module sentinel errors
var (
	ErrInvalidRoyaltyPercent          = sdkerrors.Register(ModuleName, 1100, "invalid royalty percent")
	ErrInvalidRoyaltyPercentPrecision = sdkerrors.Register(ModuleName, 1101, "invalid royalty percent precision")
	ErrInvalidDenom                   = sdkerrors.Register(ModuleName, 1102, "invalid denom")
	ErrNotDenomOwner                  = sdkerrors.Register(ModuleName, 1103, "not denom owner")
	ErrCollectionAlreadyPublished     = sdkerrors.Register(ModuleName, 1104, "collection already published")
)
