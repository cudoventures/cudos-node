package types

import (
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func ValidateRoyalties(royalties string) error {
	if royalties == "" {
		return nil
	}

	royaltiesList := strings.Split(royalties, ",")

	for _, royalty := range royaltiesList {
		royaltyParts := strings.Split(royalty, ":")

		_, err := sdk.AccAddressFromBech32(royaltyParts[0])
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid royalty address (%s)", err)
		}

		if _, err := strconv.ParseFloat(royaltyParts[1], 32); err != nil {
			return sdkerrors.Wrapf(ErrInvalidRoyaltyPercent, "invalid royalty percent (%s)", royaltyParts[1])
		}

		percentParts := strings.Split(royaltyParts[1], ".")

		if len(percentParts) == 2 && len(percentParts[1]) > 2 {
			return sdkerrors.Wrapf(ErrInvalidRoyaltyPercent, "invalid royalty percent precision (%s)", royaltyParts[1])
		}
	}

	return nil
}
