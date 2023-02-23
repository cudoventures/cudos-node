package keeper_test

import (
	"github.com/CudoVentures/cudos-node/x/nft/keeper"
)

func (suite *IntegrationTestKeeperSuite) TestGetOwners_ReturnsCorrect_Owners() {

	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address2)
	suite.NoError(err)

	err = suite.keeper.IssueDenom(suite.ctx, denomID2, denomNm2, schema, denomSymbol2, denomTraits, denomMinter, denomDescription, denomData, address)
	suite.NoError(err)

	_, err = suite.keeper.MintNFT(suite.ctx, denomID, tokenNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	_, err = suite.keeper.MintNFT(suite.ctx, denomID, tokenNm2, tokenURI, tokenData, address2, address2)
	suite.NoError(err)

	_, err = suite.keeper.MintNFT(suite.ctx, denomID, tokenNm3, tokenURI, tokenData, address2, address3)
	suite.NoError(err)

	owners, err := suite.keeper.GetOwners(suite.ctx)
	suite.NoError(err)
	suite.Equal(3, len(owners))

	_, err = suite.keeper.MintNFT(suite.ctx, denomID2, tokenNm, tokenURI, tokenData, address, address)
	suite.NoError(err)

	_, err = suite.keeper.MintNFT(suite.ctx, denomID2, tokenNm2, tokenURI, tokenData, address, address2)
	suite.NoError(err)

	_, err = suite.keeper.MintNFT(suite.ctx, denomID2, tokenNm3, tokenURI, tokenData, address, address3)
	suite.NoError(err)

	owners, err = suite.keeper.GetOwners(suite.ctx)
	suite.NoError(err)
	suite.Equal(3, len(owners))

	msg, fail := keeper.SupplyInvariant(suite.keeper)(suite.ctx)
	suite.False(fail, msg)
}
