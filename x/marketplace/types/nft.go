package types

func NewNft(tokenID, denomID, price, owner string) Nft {
	return Nft{
		TokenId: tokenID,
		DenomId: denomID,
		Price:   price,
		Owner:   owner,
	}
}
