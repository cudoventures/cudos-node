package types

func NewCollection(denomId, firstSaleRoyalties, resaleRoyalties, owner string, verified bool) Collection {
	return Collection{
		DenomId:            denomId,
		FirstSaleRoyalties: firstSaleRoyalties,
		ResaleRoyalties:    resaleRoyalties,
		Verified:           verified,
		Owner:              owner,
	}
}
