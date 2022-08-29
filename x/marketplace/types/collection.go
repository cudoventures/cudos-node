package types

func NewCollection(denomId string, mintRoyalties, resaleRoyalties []Royalty, owner string, verified bool) Collection {
	return Collection{
		DenomId:         denomId,
		MintRoyalties:   mintRoyalties,
		ResaleRoyalties: resaleRoyalties,
		Verified:        verified,
		Owner:           owner,
	}
}
