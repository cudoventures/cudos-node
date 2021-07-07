package types

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// NewGenesisState creates a new GenesisState object
func NewGenesisState(minter Minter, params Params) *GenesisState {
	return &GenesisState{
		Minter: minter,
		Params: params,
	}
}

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Minter: DefaultInitialMinter(),
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}
	// this line is used by starport scaffolding # genesis/types/validate
	// this line is used by starport scaffolding # ibc/genesistype/validate
	return ValidateMinter(gs.Minter)
}
