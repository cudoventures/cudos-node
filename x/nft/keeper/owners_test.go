package keeper_test

import (
	"cudos.org/cudos-node/x/nft/keeper"
)

func (suite *IntegrationTestKeeperSuite) TestGetOwners_ReturnsCorrect_Owners() {

	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, address2)
	suite.NoError(err)

	err = suite.keeper.IssueDenom(suite.ctx, denomID2, denomNm2, schema, address)
	suite.NoError(err)

	err = suite.keeper.MintNFT(suite.ctx, denomID, tokenID, tokenNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	err = suite.keeper.MintNFT(suite.ctx, denomID, tokenID2, tokenNm2, tokenURI, tokenData, address2, address2)
	suite.NoError(err)

	err = suite.keeper.MintNFT(suite.ctx, denomID, tokenID3, tokenNm3, tokenURI, tokenData, address2, address3)
	suite.NoError(err)

	owners, err := suite.keeper.GetOwners(suite.ctx)
	suite.NoError(err)
	suite.Equal(3, len(owners))

	err = suite.keeper.MintNFT(suite.ctx, denomID2, tokenID, tokenNm, tokenURI, tokenData, address, address)
	suite.NoError(err)

	err = suite.keeper.MintNFT(suite.ctx, denomID2, tokenID2, tokenNm2, tokenURI, tokenData, address, address2)
	suite.NoError(err)

	err = suite.keeper.MintNFT(suite.ctx, denomID2, tokenID3, tokenNm3, tokenURI, tokenData, address, address3)
	suite.NoError(err)

	owners, err = suite.keeper.GetOwners(suite.ctx)
	suite.NoError(err)
	suite.Equal(3, len(owners))

	msg, fail := keeper.SupplyInvariant(suite.keeper)(suite.ctx)
	suite.False(fail, msg)
}
