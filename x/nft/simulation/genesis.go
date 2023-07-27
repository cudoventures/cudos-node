package simulation

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/CudoVentures/cudos-node/x/nft/types"
)

const (
	numbers = "numbers"
	letters = "letters"
)

// RandomizedGenState generates a random GenesisState for nft
func RandomizedGenState(simState *module.SimulationState) {
	var collections []types.Collection
	var denom string

	for i, acc := range simState.Accounts {
		// 70% of accounts own an NFT
		if simState.Rand.Intn(100) < 70 {
			baseNFT := types.NewBaseNFT(
				strconv.FormatInt(int64(i), 10), // id
				simtypes.RandStringOfLength(simState.Rand, 10),
				acc.Address,
				simtypes.RandStringOfLength(simState.Rand, 45), // tokenURI
				simtypes.RandStringOfLength(simState.Rand, 10),
			)

			// 50% letters and 50% numbers
			if simState.Rand.Intn(100) < 50 {
				denom = fmt.Sprintf("%s%s", letters, strings.ToLower(simtypes.RandStringOfLength(simState.Rand, 10)))
			} else {
				denom = fmt.Sprintf("%s%s", numbers, strings.ToLower(simtypes.RandStringOfLength(simState.Rand, 10)))
			}
			collections = append(collections, types.NewCollection(
				types.Denom{
					Id:      denom,
					Name:    denom,
					Schema:  "",
					Creator: baseNFT.Owner,
					Symbol:  simtypes.RandStringOfLength(simState.Rand, 10),
				},
				types.NewNFTs(baseNFT),
			))
		}
	}

	nftGenesis := types.NewGenesisState(collections)

	bz, err := json.MarshalIndent(nftGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated %s parameters:\n%s\n", types.ModuleName, bz)

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(nftGenesis)
}
