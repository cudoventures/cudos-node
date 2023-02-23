package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// AddressKeyPrefix is the prefix to retrieve all Address
	AddressKeyPrefix = "Address/value/"
)

// AddressKey returns the store key to retrieve a Address from the index fields
func AddressKey(owner, network, label string) []byte {
	var key []byte

	ownerBytes := []byte(owner)
	key = append(key, ownerBytes...)
	key = append(key, []byte("/")...)
	networkBytes := []byte(network)
	key = append(key, networkBytes...)
	key = append(key, []byte("/")...)
	labelBytes := []byte(label)
	key = append(key, labelBytes...)
	key = append(key, []byte("/")...)

	return key
}
