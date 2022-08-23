package types

import "encoding/binary"

const (
	// ModuleName defines the module name
	ModuleName = "token"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_token"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

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
