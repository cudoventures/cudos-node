package simulation

import (
	"encoding/json"
	"fmt"
	"strconv"

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
	collections := types.NewCollections(
		types.NewCollection(
			types.Denom{
				Id:      letters,
				Name:    letters,
				Schema:  "",
				Creator: "",
				Symbol:  simtypes.RandStringOfLength(simState.Rand, 10),
			},
			types.NFTs{},
		),
		types.NewCollection(
			types.Denom{
				Id:      numbers,
				Name:    numbers,
				Schema:  "",
				Creator: "",
				Symbol:  simtypes.RandStringOfLength(simState.Rand, 10),
			},
			types.NFTs{}),
	)
	for i, acc := range simState.Accounts {
		// 10% of accounts own an NFT
		if simState.Rand.Intn(100) < 10 {
			baseNFT := types.NewBaseNFT(
				strconv.FormatInt(int64(i), 10), // id
				simtypes.RandStringOfLength(simState.Rand, 10),
				acc.Address,
				simtypes.RandStringOfLength(simState.Rand, 45), // tokenURI
				simtypes.RandStringOfLength(simState.Rand, 10),
			)

			// 50% letters and 50% numbers
			if simState.Rand.Intn(100) < 50 {
				collections[0].Denom.Creator = baseNFT.Owner
				collections[0] = collections[0].AddNFT(baseNFT)
			} else {
				collections[1].Denom.Creator = baseNFT.Owner
				collections[1] = collections[1].AddNFT(baseNFT)
			}
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
