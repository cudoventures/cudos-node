package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func ValidateMintRoyalties(royalties []Royalty) error {
	requiredPercent := sdk.NewDec(100)
	return ValidateRoyalties(royalties, &requiredPercent)
}

func ValidateResaleRoyalties(royalties []Royalty) error {
	return ValidateRoyalties(royalties, nil)
}

func ValidateRoyalties(royalties []Royalty, requiredPercent *sdk.Dec) error {
	if len(royalties) == 0 {
		return nil
	}

	totalPercent := sdk.NewDecFromInt(sdk.NewInt(0))

	for _, royalty := range royalties {

		_, err := sdk.AccAddressFromBech32(royalty.Address)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid royalty address (%s): %s", royalty.Address, err)
		}

		totalPercent = totalPercent.Add(royalty.Percent)

		percentParts := strings.Split(royalty.Percent.String(), ".")

		if len(percentParts) == 2 {

			trailingZeroesCount := 0
			i := len(percentParts[1]) - 1
			for i >= 0 && percentParts[1][i] == '0' {
				trailingZeroesCount++
				i--
			}

			if len(percentParts[1])-trailingZeroesCount > 2 {
				return sdkerrors.Wrapf(ErrInvalidRoyaltyPercent, "invalid royalty percent precision (%s)", percentParts[1])
			}
		}
	}

	if totalPercent.LTE(sdk.NewDec(0)) {
		return sdkerrors.Wrapf(ErrInvalidRoyaltyPercent, "total royalty percent (%s) cannot be less than or equal to zero", totalPercent.String())
	}

	if totalPercent.GT(sdk.NewDec(100)) {
		return sdkerrors.Wrapf(ErrInvalidRoyaltyPercent, "total royalty percent (%s) cannot be greater than 100", totalPercent.String())
	}

	if requiredPercent != nil && !totalPercent.Equal(*requiredPercent) {
		return sdkerrors.Wrapf(ErrInvalidRoyaltyPercent, "total royalty percent (%s) must be equal to required royalrty percent (%s)", totalPercent.String(), requiredPercent.String())
	}

	return nil
}
