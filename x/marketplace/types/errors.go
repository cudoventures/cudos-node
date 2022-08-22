package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/marketplace module sentinel errors
var (
	ErrInvalidRoyaltyPercent          = sdkerrors.Register(ModuleName, 1100, "invalid royalty percent")
	ErrInvalidRoyaltyPercentPrecision = sdkerrors.Register(ModuleName, 1101, "invalid royalty percent precision")
	ErrEmptyDenomID                   = sdkerrors.Register(ModuleName, 1102, "empty denom id")
	ErrNotDenomOwner                  = sdkerrors.Register(ModuleName, 1103, "not denom owner")
	ErrCollectionAlreadyPublished     = sdkerrors.Register(ModuleName, 1104, "collection already published")
	ErrEmptyNftID                     = sdkerrors.Register(ModuleName, 1105, "empty nft id")
	ErrInvalidPrice                   = sdkerrors.Register(ModuleName, 1106, "invalid price")
	ErrNftNotFound                    = sdkerrors.Register(ModuleName, 1107, "nft not found")
	ErrCannotBuyOwnNft                = sdkerrors.Register(ModuleName, 1108, "cannot buy own nft")
)
