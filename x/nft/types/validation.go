package types

import (
	"regexp"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	DoNotModify = "[do-not-modify]"

	MinDenomIdLen          = 3
	MaxDenomIdLen          = 64
	MinDenomNameLen        = 3
	MaxDenomNameLen        = 64
	MinDenomSymbolLen      = 3
	MaxDenomSymbolLen      = 64
	MaxDenomTraitsLen      = 256
	MaxDenomDescriptionLen = 256

	MaxTokenURILen  = 256
	MaxTokenDataLen = 512
)

var (
	// IsAlphaNumeric only accepts [a-z0-9]
	IsAlphaNumeric = regexp.MustCompile(`^[a-z0-9]+$`).MatchString
	// IsBeginWithAlpha only begin with [a-z]
	IsBeginWithAlpha = regexp.MustCompile(`^[a-z].*`).MatchString
)

// ValidateDenomID verifies whether the  parameters are legal
func ValidateDenomID(denomID string) error {
	if len(denomID) < MinDenomIdLen || len(denomID) > MaxDenomIdLen {
		return sdkerrors.Wrapf(ErrInvalidDenom, "the length of denom id(%s) only accepts value [%d, %d]", denomID, MinDenomIdLen, MaxDenomIdLen)
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
	if len(denomName) < MinDenomNameLen || len(denomName) > MaxDenomNameLen {
		return sdkerrors.Wrapf(ErrInvalidDenomName, "the length of denom name(%s) only accepts value [%d, %d]", denomName, MinDenomNameLen, MaxDenomNameLen)
	}
	return nil
}

// ValidateDenomSymbol verifies whether the  parameters are legal
func ValidateDenomSymbol(symbol string) error {
	symbol = strings.TrimSpace(symbol)
	if len(symbol) == 0 {
		return sdkerrors.Wrapf(ErrInvalidDenomSymbol, "denom symbol(%s) can not be space", symbol)
	}
	if len(symbol) < MinDenomSymbolLen || len(symbol) > MaxDenomSymbolLen {
		return sdkerrors.Wrapf(ErrInvalidDenomSymbol, "the length of denom symbol(%s) only accepts value [%d, %d]", symbol, MinDenomNameLen, MaxDenomNameLen)
	}
	return nil
}

func ValidateDenomTraits(traits string) error {
	if traits == "" {
		return nil
	}

	if len(traits) > MaxDenomTraitsLen {
		return sdkerrors.Wrapf(ErrInvalidTraits, "the length of traits %d is exceeding max accepted length %d", len(traits), MaxDenomTraitsLen)
	}

	traitsList := strings.Split(traits, ",")
	for _, trait := range traitsList {
		if _, ok := DenomTraitsMapStrToType[trait]; !ok {
			return sdkerrors.Wrapf(ErrInvalidTraits, "denom trait(%s) is not supported.", trait)
		}
	}
	return nil
}

func ValidateMinter(minter string) error {
	if minter == "" {
		return nil
	}

	if _, err := sdk.AccAddressFromBech32(minter); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid minter address (%s)", err)
	}

	return nil
}

func ValidateDescription(description string) error {
	if len(description) > MaxDenomTraitsLen {
		return sdkerrors.Wrapf(ErrInvalidDescription, "the length of denom description %d is exceeding max accepted length %d", len(description), MaxDenomTraitsLen)
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

func ValidateTokenData(tokenData string) error {
	if len(tokenData) > MaxTokenDataLen {
		return sdkerrors.Wrapf(ErrInvalidTokenURI, "the length of nft data(%s) only accepts value [0, %d]", tokenData, MaxTokenDataLen)
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
