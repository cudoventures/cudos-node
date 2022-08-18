package types

import (
	"fmt"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		CollectionList: []Collection{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated ID in collection
	collectionIdMap := make(map[uint64]bool)
	collectionCount := gs.GetCollectionCount()
	for _, elem := range gs.CollectionList {
		if _, ok := collectionIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for collection")
		}
		if elem.Id >= collectionCount {
			return fmt.Errorf("collection id should be lower or equal than the last id")
		}
		collectionIdMap[elem.Id] = true
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
