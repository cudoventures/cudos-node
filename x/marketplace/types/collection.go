package types

func NewCollection(denomId, mintRoyalties, resaleRoyalties, owner string, verified bool) Collection {
	return Collection{
		DenomId:         denomId,
		MintRoyalties:   mintRoyalties,
		ResaleRoyalties: resaleRoyalties,
		Verified:        verified,
		Owner:           owner,
	}
}
