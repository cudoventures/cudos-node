package simulation

import (
	"math/rand"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	nfttypes "github.com/CudoVentures/cudos-node/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DenomsRandomizer struct {
	usedDenoms map[string]string
}

func NewDenomsRandomizer() *DenomsRandomizer {
	return &DenomsRandomizer{
		usedDenoms: make(map[string]string),
	}
}

func (dr *DenomsRandomizer) GetRandomDenomIdWithOwner(ctx sdk.Context, k types.NftKeeper, r *rand.Rand, unique bool) (address sdk.AccAddress, denomID string) {
	allDenoms := k.GetDenoms(ctx)
	var denoms []nfttypes.Denom

	if unique {
		for _, denomCache := range allDenoms {
			if _, found := dr.usedDenoms[denomCache.Id]; !found {
				denoms = append(denoms, denomCache)
			}
		}
	} else {
		denoms = allDenoms
	}

	denomsLen := len(denoms)
	if denomsLen == 0 {
		return nil, ""
	}

	// get random denom
	i := r.Intn(denomsLen)
	denom := denoms[i]

	if unique {
		dr.usedDenoms[denom.Id] = denom.Id
	}

	ownerAddress, _ := sdk.AccAddressFromBech32(denom.Creator)
	return ownerAddress, denom.Id
}
