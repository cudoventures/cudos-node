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
	ErrCollectionNotFound             = sdkerrors.Register(ModuleName, 1109, "collection not published for sale")
	ErrNotNftOwner                    = sdkerrors.Register(ModuleName, 1110, "not nft owner")
	ErrNftAlreadyPublished            = sdkerrors.Register(ModuleName, 1111, "nft already published")
	ErrAlreadyAdmin                   = sdkerrors.Register(ModuleName, 1112, "already admin")
	ErrCollectionAlreadyVerified      = sdkerrors.Register(ModuleName, 1113, "collection is already verified")
	ErrCollectionAlreadyUnverified    = sdkerrors.Register(ModuleName, 1114, "collection is already unverified")
	ErrCollectionIsUnverified         = sdkerrors.Register(ModuleName, 1115, "collection is unverified")
	ErrNotCollectionOwner             = sdkerrors.Register(ModuleName, 1116, "not collection owner")
	ErrNotAdmin                       = sdkerrors.Register(ModuleName, 1117, "not admin")
	ErrAuctionExpired                 = sdkerrors.Register(ModuleName, 1118, "auction expired")
	ErrAuctionNotFound                = sdkerrors.Register(ModuleName, 1119, "auction not found")
	ErrInvalidAuctionDuration         = sdkerrors.Register(ModuleName, 1120, "invalid auction duration")
	ErrInvalidAuctionId               = sdkerrors.Register(ModuleName, 1121, "invalid auction id")
)
