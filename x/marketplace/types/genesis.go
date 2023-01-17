package types

import (
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

var _ codectypes.UnpackInterfacesMessage = GenesisState{}

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		CollectionList: []Collection{},
		NftList:        []Nft{},
		AuctionList:    []*codectypes.Any{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated ID in collection
	collectionIdMap := make(map[uint64]bool)
	collectionCount := gs.CollectionCount
	for _, elem := range gs.CollectionList {
		if _, ok := collectionIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for collection")
		}
		if elem.Id >= collectionCount {
			return fmt.Errorf("collection id should be lower or equal than the last id")
		}
		collectionIdMap[elem.Id] = true
	}
	// Check for duplicated ID in nft
	nftIdMap := make(map[uint64]bool)
	nftCount := gs.NftCount
	for _, elem := range gs.NftList {
		if _, ok := nftIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for nft")
		}
		if elem.Id >= nftCount {
			return fmt.Errorf("nft id should be lower or equal than the last id")
		}
		nftIdMap[elem.Id] = true
	}
	// Check for duplicated ID in auction
	auctionIdMap := make(map[uint64]bool)
	auctionCount := gs.AuctionCount
	auctionList, err := UnpackAuctions(gs.AuctionList)
	if err != nil {
		return err
	}

	for _, a := range auctionList {
		if _, ok := auctionIdMap[a.GetId()]; ok {
			return fmt.Errorf("duplicated id for auction")
		}
		if a.GetId() >= auctionCount {
			return fmt.Errorf("auction id should be lower or equal than the last id")
		}
		auctionIdMap[a.GetId()] = true
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}

func (gs GenesisState) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	for _, any := range gs.AuctionList {
		var a Auction
		err := unpacker.UnpackAny(any, &a)
		if err != nil {
			return err
		}
	}

	return nil
}
