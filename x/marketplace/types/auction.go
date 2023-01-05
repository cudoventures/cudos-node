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

func NewAuction(creator string, denomId string, tokenId string, endTime time.Time, at AuctionType) (Auction, error) {
	a := Auction{
		Creator: creator,
		TokenId: tokenId,
		DenomId: denomId,
		EndTime: endTime,
	}

	err := a.SetAuctionType(at)

	return a, err
}

func (a *Auction) GetAuctionType() (AuctionType, error) {
	at, ok := a.Type.GetCachedValue().(AuctionType)
	if !ok {
		return nil, sdkerrors.ErrInvalidType.Wrapf("expected %T, got %T", (AuctionType)(nil), a.Type.GetCachedValue())
	}
	return at, nil
}

func (a *Auction) SetAuctionType(at AuctionType) error {
	any, err := codectypes.NewAnyWithValue(at)
	if err != nil {
		return err
	}
	a.Type = any
	return nil
}

func (a Auction) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var at AuctionType
	return unpacker.UnpackAny(a.Type, &at)
}

var _ AuctionType = (*EnglishAuction)(nil)

func (a *EnglishAuction) ValidateBasic() error {
	if a.MinPrice.Validate() != nil || a.MinPrice.IsZero() {
		return ErrInvalidPrice
	}

	return nil
}

var _ AuctionType = (*DutchAuction)(nil)

func (a *DutchAuction) ValidateBasic() error {
	// todo
	return nil
}
