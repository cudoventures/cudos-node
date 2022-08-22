package types

import (
	"regexp"
	"strconv"
	"strings"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	DoNotModify = "[do-not-modify]"
	MinDenomLen = 3
	MaxDenomLen = 64

	MaxTokenURILen = 256
)

var (
	// IsAlphaNumeric only accepts [a-z0-9]
	IsAlphaNumeric = regexp.MustCompile(`^[a-z0-9]+$`).MatchString
	// IsBeginWithAlpha only begin with [a-z]
	IsBeginWithAlpha = regexp.MustCompile(`^[a-z].*`).MatchString
)

// ValidateDenomID verifies whether the  parameters are legal
func ValidateDenomID(denomID string) error {
	if len(denomID) < MinDenomLen || len(denomID) > MaxDenomLen {
		return sdkerrors.Wrapf(ErrInvalidDenom, "the length of denom(%s) only accepts value [%d, %d]", denomID, MinDenomLen, MaxDenomLen)
	}
	if !IsBeginWithAlpha(denomID) || !IsAlphaNumeric(denomID) {
		return sdkerrors.Wrapf(ErrInvalidDenom, "the denom(%s) only accepts lowercase alphanumeric characters, and begin with an english letter", denomID)
	}
	return nil
}

// ValidateDenomName verifies whether the  parameters are legal
func ValidateDenomName(denomName string) error {
	denomName = strings.TrimSpace(denomName)
	if len(denomName) == 0 {
		return sdkerrors.Wrapf(ErrInvalidDenomName, "denom name(%s) can not be space", denomName)
	}
	return nil
}

// ValidateDenomSymbol verifies whether the  parameters are legal
func ValidateDenomSymbol(symbol string) error {
	symbol = strings.TrimSpace(symbol)
	if len(symbol) == 0 {
		return sdkerrors.Wrapf(ErrInvalidDenomSymbol, "denom symbol(%s) can not be space", symbol)
	}
	return nil
}

func ValidateDenomTraits(traits string) error {
	traits = strings.TrimSpace(traits)
	if traits == "" {
		return nil
	}

	traitsList := strings.Split(traits, ",")
	for _, trait := range traitsList {
		if _, ok := DenomTraitsMapStrToType[trait]; !ok {
			return sdkerrors.Wrapf(ErrInvalidTraits, "denom trait(%s) is not supported.", trait)
		}
	}
	return nil
}

// ValidateTokenID verify that the tokenID is legal
func ValidateTokenID(tokenID string) error {

	_, err := isUint64(tokenID)
	if err != nil {
		return err
	}
	return nil
}

// ValidateTokenURI verify that the tokenURI is legal
func ValidateTokenURI(tokenURI string) error {
	if len(tokenURI) > MaxTokenURILen {
		return sdkerrors.Wrapf(ErrInvalidTokenURI, "the length of nft uri(%s) only accepts value [0, %d]", tokenURI, MaxTokenURILen)
	}
	return nil
}

// Modified returns whether the field is modified
func Modified(target string) bool {
	return target != DoNotModify
}

func isUint64(v string) (bool, error) {
	if val, err := strconv.ParseInt(v, 10, 64); err == nil {
		if val > 0 {
			return true, nil
		} else {
			return false, sdkerrors.Wrapf(ErrInvalidTokenID, "The tokenId must be a positive integer, you passed [%s]", v)
		}
	}

	return false, sdkerrors.Wrapf(ErrInvalidTokenID, "The tokenId must be a positive integer, you passed [%s]", v)

}
