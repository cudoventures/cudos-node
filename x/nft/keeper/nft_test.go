package keeper_test

import "cudos.org/cudos-node/x/nft/types"

func (suite *IntegrationTestKeeperSuite) TestGetNFT_ShouldCorrectly_ReturnNFT() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, address)
	suite.NoError(err)

	err = suite.keeper.MintNFT(suite.ctx, denomID, tokenID, denomNm, tokenURI, tokenData, address, address2)
	suite.NoError(err)

	nft, err := suite.keeper.GetBaseNFT(suite.ctx, denomID, tokenID)
	suite.NoError(err)
	suite.Equal(nft.GetID(), tokenID)
	suite.True(nft.GetOwner().Equals(address2))
	suite.Equal(nft.GetURI(), tokenURI)
}

func (suite *IntegrationTestKeeperSuite) TestGetNFT_ShouldErr_WhenNFTNotFound() {
	_, err := suite.keeper.GetBaseNFT(suite.ctx, denomID, tokenID)
	suite.ErrorIs(err, types.ErrNotFoundNFT)
}

func (suite *IntegrationTestKeeperSuite) TestGetNFTs_CorrectlyReturns_CollectionOfNFTs() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID2, denomNm, schema, address)
	suite.NoError(err)

	err = suite.keeper.MintNFT(suite.ctx, denomID2, tokenID, tokenNm, tokenURI, tokenData, address, address2)
	suite.NoError(err)

	err = suite.keeper.MintNFT(suite.ctx, denomID2, tokenID2, tokenNm2, tokenURI, tokenData, address, address2)
	suite.NoError(err)

	err = suite.keeper.MintNFT(suite.ctx, denomID2, tokenID3, tokenNm3, tokenURI, tokenData, address, address2)
	suite.NoError(err)

	nfts := suite.keeper.GetNFTs(suite.ctx, denomID2)
	suite.Len(nfts, 3)
}

func (suite *IntegrationTestKeeperSuite) TestHasNFT_ReturnsFalse_WhenNFTDoesNotExist() {
	isNFT := suite.keeper.HasNFT(suite.ctx, denomID, tokenID)
	suite.False(isNFT)
}

func (suite *IntegrationTestKeeperSuite) TestHasNFT_ReturnsCorrect_WhenNFTDoesExist() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID2, denomNm, schema, address)
	suite.NoError(err)

	err = suite.keeper.MintNFT(suite.ctx, denomID2, tokenID, tokenNm, tokenURI, tokenData, address, address2)
	suite.NoError(err)

	isNFT := suite.keeper.HasNFT(suite.ctx, denomID2, tokenID)
	suite.True(isNFT)
}

//TODO: Test approve and revoke
