package types_test

import (
	"testing"
	"time"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

type InvalidAuction struct {
	types.BaseAuction
}

func (a *InvalidAuction) SetBaseAuction(ba *types.BaseAuction) {}

var (
	startPrice   = sdk.NewCoin("acudos", sdk.NewIntFromUint64(100))
	minPrice     = sdk.NewCoin("acudos", sdk.OneInt())
	duration     = time.Hour * 24
	now          = time.Date(2222, 1, 1, 1, 0, 0, 0, time.UTC)
	endTime      = now.Add(duration)
	discountTime = now.Add(time.Hour * 1)
)

func TestAuctionFromMsgPublishAuction(t *testing.T) {
	for _, tc := range []struct {
		desc        string
		wantAuction types.Auction
		wantErr     error
	}{
		{
			desc: "english auction",
			wantAuction: &types.EnglishAuction{
				BaseAuction: &types.BaseAuction{StartTime: now, EndTime: endTime},
				MinPrice:    minPrice,
			},
		},
		{
			desc: "dutch auction",
			wantAuction: &types.DutchAuction{
				BaseAuction:      &types.BaseAuction{StartTime: now, EndTime: endTime},
				StartPrice:       startPrice,
				MinPrice:         minPrice,
				CurrentPrice:     &startPrice,
				NextDiscountTime: &discountTime,
			},
		},
		{
			desc:        "invalid auction type",
			wantAuction: &InvalidAuction{},
			wantErr:     sdkerrors.ErrInvalidType,
		},
		{
			desc:    "invalid auction",
			wantErr: sdkerrors.ErrInvalidType,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			msg := &types.MsgPublishAuction{Duration: duration}
			if tc.wantAuction != nil {
				a, err := types.PackAuction(tc.wantAuction)
				require.NoError(t, err)
				msg.Auction = a
			} else {
				msg.Auction = &codectypes.Any{}
			}

			haveAuction, err := types.AuctionFromMsgPublishAuction(msg, now)
			require.ErrorIs(t, err, tc.wantErr)
			if tc.wantErr == nil {
				require.Equal(t, tc.wantAuction, haveAuction)
			}
		})
	}
}

func TestPackAuctions(t *testing.T) {
	for _, tc := range []struct {
		desc         string
		wantAuctions []types.Auction
		wantErr      error
	}{
		{
			desc:         "valid",
			wantAuctions: []types.Auction{&types.EnglishAuction{}, &types.DutchAuction{}},
		},
		{
			desc:         "invalid",
			wantAuctions: []types.Auction{nil},
			wantErr:      sdkerrors.ErrInvalidType,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			wantAnys := make([]*codectypes.Any, len(tc.wantAuctions))
			for i, wantAuction := range tc.wantAuctions {
				haveAny, err := types.PackAuction(wantAuction)
				require.ErrorIs(t, err, tc.wantErr)

				wantAny, err := codectypes.NewAnyWithValue(wantAuction)
				require.Equal(t, tc.wantErr == nil, err == nil)

				require.Equal(t, wantAny, haveAny)

				wantAnys[i] = wantAny
			}

			haveAnys, err := types.PackAuctions(tc.wantAuctions)
			require.ErrorIs(t, err, tc.wantErr)
			if tc.wantErr == nil {
				require.Equal(t, wantAnys, haveAnys)
			}
		})
	}
}
func TestUnpackAuctions(t *testing.T) {
	englishAny, err := codectypes.NewAnyWithValue(&types.EnglishAuction{})
	require.NoError(t, err)
	dutchAny, err := codectypes.NewAnyWithValue(&types.DutchAuction{})
	require.NoError(t, err)

	for _, tc := range []struct {
		desc     string
		wantAnys []*codectypes.Any
		wantErr  error
	}{
		{
			desc:     "valid",
			wantAnys: []*codectypes.Any{englishAny, dutchAny},
		},
		{
			desc:     "invalid",
			wantAnys: []*codectypes.Any{&codectypes.Any{}, nil},
			wantErr:  sdkerrors.ErrInvalidType,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			wantAuctions := make([]types.Auction, len(tc.wantAnys))
			for i, wantAny := range tc.wantAnys {
				haveAuction, err := types.UnpackAuction(wantAny)
				require.ErrorIs(t, err, tc.wantErr)

				var wantAuction types.Auction
				if wantAny != nil {
					a, ok := wantAny.GetCachedValue().(types.Auction)
					require.Equal(t, tc.wantErr == nil, ok)
					wantAuction = a
				}

				require.Equal(t, wantAuction, haveAuction)

				wantAuctions[i] = wantAuction
			}

			haveAuctions, err := types.UnpackAuctions(tc.wantAnys)
			require.ErrorIs(t, err, tc.wantErr)
			if tc.wantErr == nil {
				require.Equal(t, wantAuctions, haveAuctions)
			}
		})
	}
}

func TestDutchAuction_IsDiscountTime(t *testing.T) {
	a := &types.DutchAuction{
		BaseAuction:      &types.BaseAuction{StartTime: now, EndTime: endTime},
		NextDiscountTime: &discountTime,
	}

	for _, tc := range []struct {
		desc             string
		arrange          func(a *types.DutchAuction)
		time             time.Time
		wantDiscountTime bool
	}{
		{
			desc:             "discount time",
			arrange:          func(a *types.DutchAuction) {},
			time:             now.Add(time.Hour * 2),
			wantDiscountTime: true,
		},
		{
			desc:    "not discount time",
			arrange: func(a *types.DutchAuction) {},
			time:    now,
		},
		{
			desc: "nil NextDiscountTime",
			arrange: func(a *types.DutchAuction) {
				a.NextDiscountTime = nil
			},
			time: now,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			require.Equal(t, tc.wantDiscountTime, a.IsDiscountTime(tc.time))
		})
	}
}

func TestDutchAuction_ApplyPriceDiscount(t *testing.T) {
	a := types.NewDutchAuction("", "", "", startPrice, minPrice, now, endTime)
	wantNextDiscountTime := now.Add(time.Hour * 1)
	wantCurrentPrice := startPrice
	discountAmount := a.StartPrice.Sub(a.MinPrice).Amount.QuoRaw(int64(duration.Hours()))

	for i := 1; i <= int(duration.Hours()); i++ {
		a.ApplyPriceDiscount()

		wantNextDiscountTime = wantNextDiscountTime.Add(time.Hour * 1)
		wantCurrentPrice = wantCurrentPrice.SubAmount(discountAmount)

		if i < int(duration.Hours())-1 {
			require.Equal(t, wantNextDiscountTime, *a.NextDiscountTime)
			require.Equal(t, wantCurrentPrice, *a.CurrentPrice)
		} else {
			require.Nil(t, a.NextDiscountTime)
			require.Equal(t, minPrice, *a.CurrentPrice)
		}
	}

	a.ApplyPriceDiscount()
	require.Equal(t, minPrice, *a.CurrentPrice)
	require.Nil(t, a.NextDiscountTime)

	a.NextDiscountTime = &discountTime
	a.CurrentPrice = nil
	require.Equal(t, discountTime, *a.NextDiscountTime)
	require.Nil(t, a.CurrentPrice)
}
