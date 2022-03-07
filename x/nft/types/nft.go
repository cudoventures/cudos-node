package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/CudoVentures/cudos-node/x/nft/exported"
)

var _ exported.NFT = BaseNFT{}

// NewBaseNFT creates a new NFT instance
func NewBaseNFT(id, name string, owner sdk.AccAddress, uri, data string) BaseNFT {
	return BaseNFT{
		Id:    id,
		Name:  name,
		Owner: owner.String(),
		URI:   uri,
		Data:  data,
	}
}

// GetID return the id of BaseNFT
func (bnft BaseNFT) GetID() string {
	return bnft.Id
}

// GetName return the name of BaseNFT
func (bnft BaseNFT) GetName() string {
	return bnft.Name
}

// GetOwner return the owner of BaseNFT
func (bnft BaseNFT) GetOwner() sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(bnft.Owner)

	if err != nil {
		panic(fmt.Errorf("couldn't convert %q to account address: %v", bnft.Owner, err))
	}

	return owner
}

// GetURI return the URI of BaseNFT
func (bnft BaseNFT) GetURI() string {
	return bnft.URI
}

// GetData return the Data of BaseNFT
func (bnft BaseNFT) GetData() string {
	return bnft.Data
}

// GetApprovedAddresses return the Data of BaseNFT
func (bnft BaseNFT) GetApprovedAddresses() []string {
	return bnft.ApprovedAddresses
}

// ----------------------------------------------------------------------------
// NFT

// NFTs define a list of NFT
type NFTs []exported.NFT

// NewNFTs creates a new set of NFTs
func NewNFTs(nfts ...exported.NFT) NFTs {
	if len(nfts) == 0 {
		return NFTs{}
	}
	return NFTs(nfts)
}
