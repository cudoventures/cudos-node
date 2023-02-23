package keeper_test

import (
	"github.com/CudoVentures/cudos-node/x/nft/keeper"
	"github.com/CudoVentures/cudos-node/x/nft/types"
)

func (suite *IntegrationTestKeeperSuite) TestSetCollection_Correctly_MintsNFTsFromCollection() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address)
	suite.NoError(err)
	nft := types.NewBaseNFT("", tokenNm, address, tokenURI, tokenData)
	nft2 := types.NewBaseNFT("", tokenNm, address, tokenURI, tokenData)

	denomE := types.Denom{
		Id:      denomID,
		Name:    denomNm,
		Schema:  schema,
		Creator: address.String(),
	}

	collection2 := types.Collection{
		Denom: denomE,
		NFTs:  []types.BaseNFT{nft2, nft},
	}

	err = suite.keeper.SetCollection(suite.ctx, collection2, address)
	suite.NoError(err)

	collection2, err = suite.keeper.GetCollection(suite.ctx, denomID)
	suite.NoError(err)
	suite.Len(collection2.NFTs, 2)

	msg, fail := keeper.SupplyInvariant(suite.keeper)(suite.ctx)
	suite.False(fail, msg)
}

func (suite *IntegrationTestKeeperSuite) TestGetCollection_Returns_ValidCollection() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address2)
	suite.NoError(err)

	_, err = suite.keeper.MintNFT(suite.ctx, denomID, tokenNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	// collection should exist
	collection, err := suite.keeper.GetCollection(suite.ctx, denomID)
	suite.NoError(err)
	suite.NotEmpty(collection)

	msg, fail := keeper.SupplyInvariant(suite.keeper)(suite.ctx)
	suite.False(fail, msg)
}

func (suite *IntegrationTestKeeperSuite) TestGetSupply() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address2)
	suite.NoError(err)
	err = suite.keeper.IssueDenom(suite.ctx, denomID2, denomNm2, schema, denomSymbol2, denomTraits, denomMinter, denomDescription, denomData, address)
	suite.NoError(err)

	tokenId, err := suite.keeper.MintNFT(suite.ctx, denomID, tokenNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	tokenId2, err := suite.keeper.MintNFT(suite.ctx, denomID, tokenNm2, tokenURI, tokenData, address2, address2)
	suite.NoError(err)

	tokenId3, err := suite.keeper.MintNFT(suite.ctx, denomID2, tokenNm3, tokenURI, tokenData, address, address3)
	suite.NoError(err)

	supply := suite.keeper.GetTotalSupply(suite.ctx, denomID)
	suite.Equal(uint64(2), supply)

	supply = suite.keeper.GetTotalSupply(suite.ctx, denomID2)
	suite.Equal(uint64(1), supply)

	supply = suite.keeper.GetTotalSupplyOfOwner(suite.ctx, denomID, address)
	suite.Equal(uint64(1), supply)

	supply = suite.keeper.GetTotalSupplyOfOwner(suite.ctx, denomID, address2)
	suite.Equal(uint64(1), supply)

	supply = suite.keeper.GetTotalSupply(suite.ctx, denomID)
	suite.Equal(uint64(2), supply)

	supply = suite.keeper.GetTotalSupply(suite.ctx, denomID2)
	suite.Equal(uint64(1), supply)

	err = suite.keeper.BurnNFT(suite.ctx, denomID, tokenId, address)
	suite.NoError(err)
	supply = suite.keeper.GetTotalSupply(suite.ctx, denomID)
	suite.Equal(uint64(1), supply)

	err = suite.keeper.BurnNFT(suite.ctx, denomID, tokenId2, address2)
	suite.NoError(err)
	supply = suite.keeper.GetTotalSupply(suite.ctx, denomID)
	suite.Equal(uint64(0), supply)

	err = suite.keeper.BurnNFT(suite.ctx, denomID2, tokenId3, address3)
	suite.NoError(err)
	supply = suite.keeper.GetTotalSupply(suite.ctx, denomID2)
	suite.Equal(uint64(0), supply)
}
