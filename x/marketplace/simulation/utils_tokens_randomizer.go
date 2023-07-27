package simulation

import (
	"fmt"
	"math/rand"

	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	nfttypes "github.com/CudoVentures/cudos-node/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TokensRandomizer struct {
	usedTokens map[string]string
}

func NewTokensRandomizer() *TokensRandomizer {
	return &TokensRandomizer{
		usedTokens: make(map[string]string),
	}
}

func (tr *TokensRandomizer) GetRandomTokenIdWithOwner(ctx sdk.Context, k types.NftKeeper, r *rand.Rand, unique bool, excludedOwnerAddress string) (address sdk.AccAddress, denomID, tokenID string) {
	allOwners, err := k.GetOwners(ctx)
	if err != nil {
		return nil, "", ""
	}

	var owners nfttypes.Owners
	if unique {
		for _, owner := range allOwners {
			if excludedOwnerAddress == owner.Address {
				continue
			}

			var idcs []nfttypes.IDCollection
			for _, col := range owner.IDCollections {
				var tokenIds []string
				for _, token := range col.TokenIds {
					mapKey := fmt.Sprintf("%s_%s", col.DenomId, token)
					if _, found := tr.usedTokens[mapKey]; !found {
						tokenIds = append(tokenIds, token)
					}
				}
				if len(tokenIds) > 0 {
					idcs = append(idcs, nfttypes.IDCollection{
						DenomId:  col.DenomId,
						TokenIds: tokenIds,
					})
				}
			}
			if len(idcs) > 0 {
				owners = append(
					owners,
					nfttypes.Owner{
						Address:       owner.Address,
						IDCollections: idcs,
					},
				)
			}
		}
	} else {
		owners = allOwners
	}

	ownersLen := len(owners)
	if ownersLen == 0 {
		return nil, "", ""
	}

	// get random owner
	i := r.Intn(ownersLen)
	owner := owners[i]

	idCollectionsLen := len(owner.IDCollections)
	if idCollectionsLen == 0 {
		return nil, "", ""
	}

	// get random collection from owner's balance
	i = r.Intn(idCollectionsLen)
	idCollection := owner.IDCollections[i] // nfts IDs
	denomID = idCollection.DenomId

	idsLen := len(idCollection.TokenIds)
	if idsLen == 0 {
		return nil, "", ""
	}

	// get random nft from collection
	i = r.Intn(idsLen)
	tokenID = idCollection.TokenIds[i]

	if unique {
		mapKey := fmt.Sprintf("%s_%s", denomID, tokenID)
		tr.usedTokens[mapKey] = tokenID
	}

	ownerAddress, _ := sdk.AccAddressFromBech32(owner.Address)
	return ownerAddress, denomID, tokenID
}
