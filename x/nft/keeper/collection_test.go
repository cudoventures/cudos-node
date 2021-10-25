package keeper_test

import (
	"cudos.org/cudos-node/x/nft/keeper"
	"cudos.org/cudos-node/x/nft/types"
)

func (suite *IntegrationTestKeeperSuite) TestSetCollection_Correctly_MintsNFTsFromCollection() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, address)
	suite.NoError(err)
	nft := types.NewBaseNFT(tokenID, tokenNm, address, tokenURI, tokenData)
	nft2 := types.NewBaseNFT(tokenID2, tokenNm, address, tokenURI, tokenData)

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
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, address2)
	suite.NoError(err)

	err = suite.keeper.MintNFT(suite.ctx, denomID, tokenID, tokenNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	// collection should exist
	collection, err := suite.keeper.GetCollection(suite.ctx, denomID)
	suite.NoError(err)
	suite.NotEmpty(collection)

	msg, fail := keeper.SupplyInvariant(suite.keeper)(suite.ctx)
	suite.False(fail, msg)
}

func (suite *IntegrationTestKeeperSuite) TestGetSupply() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, address2)
	suite.NoError(err)
	err = suite.keeper.IssueDenom(suite.ctx, denomID2, denomNm2, schema, address)
	suite.NoError(err)

	err = suite.keeper.MintNFT(suite.ctx, denomID, tokenID, tokenNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	err = suite.keeper.MintNFT(suite.ctx, denomID, tokenID2, tokenNm2, tokenURI, tokenData, address2, address2)
	suite.NoError(err)

	err = suite.keeper.MintNFT(suite.ctx, denomID2, tokenID3, tokenNm3, tokenURI, tokenData, address, address3)
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

	//burn nft
	err = suite.keeper.BurnNFT(suite.ctx, denomID, tokenID, address)
	suite.NoError(err)

	supply = suite.keeper.GetTotalSupply(suite.ctx, denomID)
	suite.Equal(uint64(1), supply)

	supply = suite.keeper.GetTotalSupply(suite.ctx, denomID)
	suite.Equal(uint64(1), supply)

	//burn nft
	err = suite.keeper.BurnNFT(suite.ctx, denomID, tokenID2, address2)
	suite.NoError(err)

	supply = suite.keeper.GetTotalSupply(suite.ctx, denomID)
	suite.Equal(uint64(0), supply)

	supply = suite.keeper.GetTotalSupply(suite.ctx, denomID)
	suite.Equal(uint64(0), supply)
}
