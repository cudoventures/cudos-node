package types

import "encoding/binary"

const (
	// ModuleName defines the module name
	ModuleName = "marketplace"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_marketplace"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func Uint64ToBytes(value uint64) []byte {
	bValue := make([]byte, 8)
	binary.LittleEndian.PutUint64(bValue, value)
	return bValue
}

func BytesToUint64(value []byte) uint64 {
	return binary.LittleEndian.Uint64(value)
}

func KeyCollectionDenomID(denomID string) []byte {
	key := append(KeyPrefix(CollectionDenomIDKey), delimiter...)
	return append(key, KeyPrefix(denomID)...)
}

func KeyNftDenomTokenID(denomID, tokenID string) []byte {
	key := append(KeyPrefix(NftDenomTokenIDKey), delimiter...)
	key = append(key, KeyPrefix(denomID)...)
	key = append(key, delimiter...)
	return append(key, KeyPrefix(tokenID)...)
}

func KeyAdmins() []byte {
	return append(KeyPrefix(AdminsKey), delimiter...)
}

var delimiter = []byte("/")

const (
	CollectionKey        = "Collection-value-"
	CollectionCountKey   = "Collection-count-"
	CollectionDenomIDKey = "Collection-denom-id"
)

const (
	NftKey             = "Nft-value-"
	NftCountKey        = "Nft-count-"
	NextUniqueIdKey    = "Next-unique-id-"
	NftDenomTokenIDKey = "Nft-denom-token-id-"
)

const (
	AdminsKey = "Admins-"
)
