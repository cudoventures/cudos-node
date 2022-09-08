package keeper_test

import (
	"github.com/CudoVentures/cudos-node/x/nft/types"
	"github.com/stretchr/testify/assert"
)

func (suite *IntegrationTestKeeperSuite) TestGetNFT_ShouldCorrectly_ReturnNFT() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address)
	suite.NoError(err)

	tokenId, err := suite.keeper.MintNFT(suite.ctx, denomID, denomNm, tokenURI, tokenData, address, address2)
	suite.NoError(err)

	nft, err := suite.keeper.GetBaseNFT(suite.ctx, denomID, tokenId)
	suite.NoError(err)
	suite.Equal(nft.GetID(), tokenId)
	suite.True(nft.GetOwner().Equals(address2))
	suite.Equal(nft.GetURI(), tokenURI)
}

func (suite *IntegrationTestKeeperSuite) TestGetNFT_ShouldErr_WhenNFTNotFound() {
	_, err := suite.keeper.GetBaseNFT(suite.ctx, denomID, "1234")
	suite.ErrorIs(err, types.ErrNotFoundNFT)
}

func (suite *IntegrationTestKeeperSuite) TestGetNFTs_CorrectlyReturns_CollectionOfNFTs() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID2, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address)
	suite.NoError(err)

	_, err = suite.keeper.MintNFT(suite.ctx, denomID2, tokenNm, tokenURI, tokenData, address, address2)
	suite.NoError(err)

	_, err = suite.keeper.MintNFT(suite.ctx, denomID2, tokenNm2, tokenURI, tokenData, address, address2)
	suite.NoError(err)

	_, err = suite.keeper.MintNFT(suite.ctx, denomID2, tokenNm3, tokenURI, tokenData, address, address2)
	suite.NoError(err)

	nfts := suite.keeper.GetNFTs(suite.ctx, denomID2)
	suite.Len(nfts, 3)
}

func (suite *IntegrationTestKeeperSuite) TestHasNFT_ReturnsFalse_WhenNFTDoesNotExist() {
	isNFT := suite.keeper.HasNFT(suite.ctx, denomID, "1234")
	suite.False(isNFT)
}

func (suite *IntegrationTestKeeperSuite) TestHasNFT_ReturnsCorrect_WhenNFTDoesExist() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID2, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address)
	suite.NoError(err)

	tokenId, err := suite.keeper.MintNFT(suite.ctx, denomID2, tokenNm, tokenURI, tokenData, address, address2)
	suite.NoError(err)

	isNFT := suite.keeper.HasNFT(suite.ctx, denomID2, tokenId)
	suite.True(isNFT)
}

func (suite *IntegrationTestKeeperSuite) TestApproveNFT_ReturnsCorrect() {

	err := suite.keeper.IssueDenom(suite.ctx, denomID2, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address)
	suite.NoError(err)

	tokenId, err := suite.keeper.MintNFT(suite.ctx, denomID2, tokenNm, tokenURI, tokenData, address, address2)
	suite.NoError(err)

	nft, err := suite.keeper.GetBaseNFT(suite.ctx, denomID2, tokenId)
	suite.NoError(err)

	suite.keeper.ApproveNFT(suite.ctx, nft, address3, denomID2)

	nft, err = suite.keeper.GetBaseNFT(suite.ctx, denomID2, tokenId)
	suite.NoError(err)

	assert.Equal(suite.T(), true, suite.keeper.IsApprovedAddress(&nft, address3.String()))
}

func (suite *IntegrationTestKeeperSuite) TestRevokeApprovalNFT_ReturnsCorrect() {

	err := suite.keeper.IssueDenom(suite.ctx, denomID2, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address)
	suite.NoError(err)

	tokenId, err := suite.keeper.MintNFT(suite.ctx, denomID2, tokenNm, tokenURI, tokenData, address, address2)
	suite.NoError(err)

	nft, err := suite.keeper.GetBaseNFT(suite.ctx, denomID2, tokenId)
	suite.NoError(err)

	suite.keeper.ApproveNFT(suite.ctx, nft, address2, denomID2)
	suite.NoError(err)

	nft, err = suite.keeper.GetBaseNFT(suite.ctx, denomID2, tokenId)
	suite.NoError(err)
	assert.Equal(suite.T(), suite.keeper.IsApprovedAddress(&nft, address2.String()), true)

	err = suite.keeper.RevokeApprovalNFT(suite.ctx, nft, address2, denomID2)

	nft, err = suite.keeper.GetBaseNFT(suite.ctx, denomID2, tokenId)

	suite.NoError(err)
	assert.Equal(suite.T(), []string([]string(nil)), nft.ApprovedAddresses)

}

func (suite *IntegrationTestKeeperSuite) TestRevokeApprovalNFT_ReturnsError_WhenUserHasNoApprovedAddresses() {

	err := suite.keeper.IssueDenom(suite.ctx, denomID2, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address)
	suite.NoError(err)

	tokenId, err := suite.keeper.MintNFT(suite.ctx, denomID2, tokenNm, tokenURI, tokenData, address, address2)
	suite.NoError(err)

	nft, err := suite.keeper.GetBaseNFT(suite.ctx, denomID2, tokenId)
	suite.NoError(err)

	err = suite.keeper.RevokeApprovalNFT(suite.ctx, nft, address2, denomID2)
	suite.ErrorIs(err, types.ErrNoApprovedAddresses)

}

func (suite *IntegrationTestKeeperSuite) TestRevokeApprovalNFT_ReturnsError_WhenApprovedAddressIsNotFound() {

	err := suite.keeper.IssueDenom(suite.ctx, denomID2, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address)
	suite.NoError(err)

	tokenId, err := suite.keeper.MintNFT(suite.ctx, denomID2, tokenNm, tokenURI, tokenData, address, address2)
	suite.NoError(err)

	nft, err := suite.keeper.GetBaseNFT(suite.ctx, denomID2, tokenId)
	suite.NoError(err)

	suite.keeper.ApproveNFT(suite.ctx, nft, address2, denomID2)
	suite.NoError(err)

	nft, err = suite.keeper.GetBaseNFT(suite.ctx, denomID2, tokenId)
	suite.NoError(err)

	err = suite.keeper.RevokeApprovalNFT(suite.ctx, nft, address3, denomID2)
	suite.ErrorIs(err, types.ErrNoApprovedAddresses)

}
