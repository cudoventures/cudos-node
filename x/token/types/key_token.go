package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// TokenKeyPrefix is the prefix to retrieve all Token
	TokenKeyPrefix = "Token/value/"
)

// TokenKey returns the store key to retrieve a Token from the index fields
func TokenKey(
	denom string,
) []byte {
	var key []byte

	denomBytes := []byte(denom)
	key = append(key, denomBytes...)
	key = append(key, []byte("/")...)

	return key
}
