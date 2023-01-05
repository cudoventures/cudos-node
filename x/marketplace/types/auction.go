package types

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type AuctionType interface {
	codec.ProtoMarshaler

	ValidateBasic() error
	// HandleBid(ctx sdk.Context, bid Bid, bk BankKeeper) error
}

func NewAuction(creator string, denomId string, tokenId string, endTime time.Time, auctionType AuctionType) (Auction, error) {
	auction := Auction{
		Creator: creator,
		TokenId: tokenId,
		DenomId: denomId,
		EndTime: endTime,
	}

	err := auction.SetAuctionType(auctionType)

	return auction, err
}

func (a *Auction) GetAuctionType() (AuctionType, error) {
	auctionType, ok := a.Type.GetCachedValue().(AuctionType)
	if !ok {
		return nil, sdkerrors.ErrInvalidType.Wrapf("expected %T, got %T", (AuctionType)(nil), a.Type.GetCachedValue())
	}
	return auctionType, nil
}

func (a *Auction) SetAuctionType(auctionType AuctionType) error {
	any, err := codectypes.NewAnyWithValue(auctionType)
	if err != nil {
		return err
	}
	a.Type = any
	return nil
}

func (a Auction) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var auctionType AuctionType
	return unpacker.UnpackAny(a.Type, &auctionType)
}

var _ AuctionType = (*EnglishAuction)(nil)

func (a *EnglishAuction) ValidateBasic() error {
	if err := a.MinPrice.Validate(); err != nil {
		return sdkerrors.Wrapf(ErrInvalidPrice, "invalid minimum price (%s)", err)
	}

	if a.MinPrice.IsZero() {
		return sdkerrors.Wrapf(ErrInvalidPrice, "minimum price must be positive")
	}

	return nil
}

var _ AuctionType = (*DutchAuction)(nil)

func (a *DutchAuction) ValidateBasic() error {
	// todo
	return nil
}
