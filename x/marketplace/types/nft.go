package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func NewNft(tokenID, denomID, owner string, price sdk.Coin) Nft {
	return Nft{
		TokenId: tokenID,
		DenomId: denomID,
		Price:   price,
		Owner:   owner,
	}
}
