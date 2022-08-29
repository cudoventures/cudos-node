package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func ValidateRoyalties(royalties []Royalty) error {
	if len(royalties) == 0 {
		return nil
	}

	var totalPercent sdk.Dec

	for _, royalty := range royalties {

		_, err := sdk.AccAddressFromBech32(royalty.Address)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid royalty address (%s): %s", royalty.Address, err)
		}

		totalPercent = totalPercent.Add(royalty.Percent)

		percentParts := strings.Split(royalty.Percent.String(), ".")

		if len(percentParts) == 2 && len(percentParts[1]) > 2 {
			return sdkerrors.Wrapf(ErrInvalidRoyaltyPercent, "invalid royalty percent precision (%s)", percentParts[1])
		}
	}

	if totalPercent.LTE(sdk.NewDec(0)) {
		return sdkerrors.Wrapf(ErrInvalidRoyaltyPercent, "total royalty percent (%s) cannot be less than or equal to zero", totalPercent.String())
	}

	if totalPercent.GT(sdk.NewDec(100)) {
		return sdkerrors.Wrapf(ErrInvalidRoyaltyPercent, "total royalty percent (%s) cannot be greater than 100", totalPercent.String())
	}

	return nil
}
