package types

import (
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// TODO: There is duplication of this logic in keeper.distributeRoyalties
func ValidateRoyalties(royalties string) error {
	if royalties == "" {
		return nil
	}

	var totalPercent float64

	splitFn := func(c rune) bool {
		return c == ','
	}

	royaltiesList := strings.FieldsFunc(royalties, splitFn)

	for _, royalty := range royaltiesList {
		royaltyParts := strings.Split(royalty, ":")

		_, err := sdk.AccAddressFromBech32(royaltyParts[0])
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid royalty address (%s): %s", royaltyParts[0], err)
		}

		value, err := strconv.ParseFloat(royaltyParts[1], 32)
		if err != nil {
			return sdkerrors.Wrapf(ErrInvalidRoyaltyPercent, "invalid royalty percent (%s): %s", royaltyParts[1], err)
		}

		totalPercent += value

		percentParts := strings.Split(royaltyParts[1], ".")

		if len(percentParts) == 2 && len(percentParts[1]) > 2 {
			return sdkerrors.Wrapf(ErrInvalidRoyaltyPercent, "invalid royalty percent precision (%s)", royaltyParts[1])
		}
	}

	if totalPercent <= 0 {
		return sdkerrors.Wrapf(ErrInvalidRoyaltyPercent, "total royalty percent (%f) cannot be less than or equal to zero", totalPercent)
	}

	if totalPercent > 100 {
		return sdkerrors.Wrapf(ErrInvalidRoyaltyPercent, "total royalty percent (%f) cannot be greater than 100", totalPercent)
	}

	return nil
}
