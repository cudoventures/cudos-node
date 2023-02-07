package types

import (
	"encoding/json"
	"time"

	nfttypes "github.com/CudoVentures/cudos-node/x/nft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ Auction = (*EnglishAuction)(nil)
	_ Auction = (*DutchAuction)(nil)
)

type Auction interface {
	codec.ProtoMarshaler

	ValidateBasic() error
	GetId() uint64
	GetDenomId() string
	GetTokenId() string
	GetCreator() string
	GetEndTime() time.Time
	SetId(id uint64)
	SetCreator(creator string)
	SetBaseAuction(a *BaseAuction)
	MarshalJSON() ([]byte, error)
}

func AuctionFromMsgPublishAuction(
	msg *MsgPublishAuction, startTime time.Time,
) (Auction, error) {
	a, err := msg.GetAuction()
	if err != nil {
		return nil, err
	}

	endTime := startTime.Add(msg.Duration)

	switch a := a.(type) {
	case *EnglishAuction:
		return NewEnglishAuction(
			msg.Creator,
			msg.DenomId,
			msg.TokenId,
			a.MinPrice,
			startTime,
			endTime,
		), nil
	case *DutchAuction:
		return NewDutchAuction(
			msg.Creator,
			msg.DenomId,
			msg.TokenId,
			a.StartPrice,
			a.MinPrice,
			startTime,
			endTime,
		), nil
	default:
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidType, "invalid auction type")
	}
}

func PackAuction(a Auction) (*codectypes.Any, error) {
	auctionAny, err := codectypes.NewAnyWithValue(a)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidType, "invalid auction")
	}

	return auctionAny, nil
}

func PackAuctions(auctions []Auction) ([]*codectypes.Any, error) {
	auctionsAny := make([]*codectypes.Any, len(auctions))
	for i, a := range auctions {
		auctionAny, err := codectypes.NewAnyWithValue(a)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidType, "invalid auction")
		}
		auctionsAny[i] = auctionAny
	}
	return auctionsAny, nil
}

func UnpackAuction(auctionAny *codectypes.Any) (Auction, error) {
	if auctionAny == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidType, "invalid auction")
	}

	a, ok := auctionAny.GetCachedValue().(Auction)
	if !ok {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidType, "invalid auction")
	}
	return a, nil
}

func UnpackAuctions(auctionsAny []*codectypes.Any) ([]Auction, error) {
	auctions := make([]Auction, len(auctionsAny))
	for i, a := range auctionsAny {
		a, ok := a.GetCachedValue().(Auction)
		if !ok {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidType, "invalid auction")
		}
		auctions[i] = a
	}
	return auctions, nil
}

func NewBaseAuction(
	creator string,
	denomId string,
	tokenId string,
	startTime time.Time,
	endTime time.Time,
) *BaseAuction {
	return &BaseAuction{
		Creator:   creator,
		TokenId:   tokenId,
		DenomId:   denomId,
		StartTime: startTime,
		EndTime:   endTime,
	}
}

func (a *BaseAuction) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(a.Creator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator %s", err)
	}

	if err := nfttypes.ValidateDenomID(a.DenomId); err != nil {
		return nfttypes.ErrInvalidDenom
	}

	if err := nfttypes.ValidateTokenID(a.TokenId); err != nil {
		return nfttypes.ErrInvalidTokenID
	}

	if a.EndTime.Sub(a.StartTime) < time.Hour*24 {
		return sdkerrors.Wrap(ErrInvalidAuctionDuration, "duration is less than 24 hours")
	}

	return nil
}

func (a *BaseAuction) GetId() uint64 {
	return a.Id
}

func (a *BaseAuction) GetDenomId() string {
	return a.DenomId
}

func (a *BaseAuction) GetTokenId() string {
	return a.TokenId
}

func (a *BaseAuction) GetCreator() string {
	return a.Creator
}

func (a *BaseAuction) GetEndTime() time.Time {
	return a.EndTime
}

func (a *BaseAuction) SetId(id uint64) {
	a.Id = id
}

func (a *BaseAuction) SetCreator(creator string) {
	a.Creator = creator
}

func NewEnglishAuction(
	creator string,
	denomId string,
	tokenId string,
	minPrice sdk.Coin,
	startTime time.Time,
	endTime time.Time,
) *EnglishAuction {
	return &EnglishAuction{
		BaseAuction: NewBaseAuction(creator, denomId, tokenId, startTime, endTime),
		MinPrice:    minPrice,
	}
}

func (a *EnglishAuction) ValidateBasic() error {
	if a.BaseAuction != nil {
		if err := a.BaseAuction.ValidateBasic(); err != nil {
			return err
		}
	}

	if err := a.MinPrice.Validate(); err != nil {
		return sdkerrors.Wrapf(ErrInvalidPrice, "invalid min price: %s", err)
	}

	if a.MinPrice.IsZero() {
		return sdkerrors.Wrapf(ErrInvalidPrice, "min price must be positive")
	}

	return nil
}

func (a *EnglishAuction) SetBaseAuction(ba *BaseAuction) {
	a.BaseAuction = ba
}

type englishAuctionJson struct {
	Type     string   `json:"@type"`
	MinPrice sdk.Coin `json:"minPrice"`
}

func (a *EnglishAuction) MarshalJSON() ([]byte, error) {
	return json.Marshal(englishAuctionJson{
		Type:     "/cudoventures.cudosnode.marketplace.EnglishAuction",
		MinPrice: a.MinPrice,
	})
}

func NewDutchAuction(
	creator string,
	denomId string,
	tokenId string,
	startPrice sdk.Coin,
	minPrice sdk.Coin,
	startTime time.Time,
	endTime time.Time,
) *DutchAuction {
	nextDiscountTime := startTime.Add(time.Hour * 1)
	return &DutchAuction{
		BaseAuction:      NewBaseAuction(creator, denomId, tokenId, startTime, endTime),
		StartPrice:       startPrice,
		MinPrice:         minPrice,
		CurrentPrice:     &sdk.Coin{Denom: startPrice.Denom, Amount: startPrice.Amount},
		NextDiscountTime: &nextDiscountTime,
	}
}

func (a *DutchAuction) ValidateBasic() error {
	if a.BaseAuction != nil {
		if err := a.BaseAuction.ValidateBasic(); err != nil {
			return err
		}
	}

	if err := a.StartPrice.Validate(); err != nil {
		return sdkerrors.Wrapf(ErrInvalidPrice, "invalid start price: %s", err)
	}

	if a.StartPrice.IsZero() {
		return sdkerrors.Wrapf(ErrInvalidPrice, "start price must be positive")
	}

	if err := a.MinPrice.Validate(); err != nil {
		return sdkerrors.Wrapf(ErrInvalidPrice, "invalid min price: %s", err)
	}

	if a.MinPrice.IsZero() {
		return sdkerrors.Wrapf(ErrInvalidPrice, "min price must be positive")
	}

	if a.StartPrice.IsLT(a.MinPrice) {
		return sdkerrors.Wrapf(ErrInvalidPrice, "start price is lower than min price")
	}

	return nil
}

func (a *DutchAuction) SetBaseAuction(ba *BaseAuction) {
	a.BaseAuction = ba
}

type dutchAuctionJson struct {
	Type             string     `json:"@type"`
	StartPrice       sdk.Coin   `json:"startPrice"`
	MinPrice         sdk.Coin   `json:"minPrice"`
	CurrentPrice     *sdk.Coin  `json:"currentPrice"`
	NextDiscountTime *time.Time `json:"nextDiscountTime"`
}

func (a *DutchAuction) MarshalJSON() ([]byte, error) {
	return json.Marshal(dutchAuctionJson{
		Type:             "/cudoventures.cudosnode.marketplace.DutchAuction",
		StartPrice:       a.StartPrice,
		MinPrice:         a.MinPrice,
		CurrentPrice:     a.CurrentPrice,
		NextDiscountTime: a.NextDiscountTime,
	})
}

func (a *DutchAuction) IsDiscountTime(time time.Time) bool {
	return a.NextDiscountTime != nil && time.After(*a.NextDiscountTime)
}

func (a *DutchAuction) ApplyPriceDiscount() {
	if a.NextDiscountTime == nil || a.CurrentPrice == nil {
		return
	}

	durationHours := int64(a.EndTime.Sub(a.StartTime).Hours()) - 1
	discount := a.StartPrice.Sub(a.MinPrice).Amount.QuoRaw(durationHours)
	priceAfterDiscount := a.CurrentPrice.SubAmount(discount)
	a.CurrentPrice = &priceAfterDiscount

	if a.EndTime.Sub(*a.NextDiscountTime) >= time.Hour*2 {
		nextDiscountTime := a.NextDiscountTime.Add(time.Hour * 1)
		a.NextDiscountTime = &nextDiscountTime
	} else {
		a.NextDiscountTime = nil
		a.CurrentPrice = &sdk.Coin{Denom: a.MinPrice.Denom, Amount: a.MinPrice.Amount}
	}
}
