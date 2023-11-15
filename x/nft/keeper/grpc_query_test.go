package keeper_test

import (
	gocontext "context"

	"github.com/CudoVentures/cudos-node/x/nft/types"
)

func (suite *IntegrationTestKeeperSuite) TestSupplyReturnsCorrect() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address2)
	suite.NoError(err)

	_, err = suite.keeper.MintNFT(suite.ctx, denomID, tokenNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	response, err := suite.queryClient.Supply(gocontext.Background(), &types.QuerySupplyRequest{
		DenomID: denomID,
		Owner:   address.String(),
	})

	suite.NoError(err)
	suite.Equal(1, int(response.Amount))
}

func (suite *IntegrationTestKeeperSuite) TestOwner_ReturnsCorrect() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address2)
	suite.NoError(err)

	tokenID, err := suite.keeper.MintNFT(suite.ctx, denomID, tokenNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	response, err := suite.queryClient.Owner(gocontext.Background(), &types.QueryOwnerRequest{
		DenomID: denomID,
		Owner:   address.String(),
	})

	suite.NoError(err)
	suite.NotNil(response.Owner)
	suite.Contains(response.Owner.IDCollections[0].TokenIds, tokenID)
}

func (suite *IntegrationTestKeeperSuite) TestCollection_ReturnsCorrect() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address2)
	suite.NoError(err)

	tokenID, err := suite.keeper.MintNFT(suite.ctx, denomID, tokenNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	response, err := suite.queryClient.Collection(gocontext.Background(), &types.QueryCollectionRequest{
		DenomID: denomID,
	})

	suite.NoError(err)
	suite.NotNil(response.Collection)
	suite.Len(response.Collection.NFTs, 1)
	suite.Equal(response.Collection.NFTs[0].Id, tokenID)
}

func (suite *IntegrationTestKeeperSuite) TestDenom_ReturnsCorrect() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address2)
	suite.NoError(err)

	_, err = suite.keeper.MintNFT(suite.ctx, denomID, tokenNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	response, err := suite.queryClient.Denom(gocontext.Background(), &types.QueryDenomRequest{
		DenomID: denomID,
	})

	suite.NoError(err)
	suite.NotNil(response.Denom)
	suite.Equal(response.Denom.Id, denomID)
}

func (suite *IntegrationTestKeeperSuite) TestDenoms_ReturnsCorrectCollection() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address2)
	suite.NoError(err)

	_, err = suite.keeper.MintNFT(suite.ctx, denomID, tokenNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	response, err := suite.queryClient.Denoms(gocontext.Background(), &types.QueryDenomsRequest{})

	suite.NoError(err)
	suite.NotEmpty(response.Denoms)
	suite.Equal(response.Denoms[0].Id, denomID)
}

func (suite *IntegrationTestKeeperSuite) TestNFT_ReturnsCorrect() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address2)
	suite.NoError(err)

	tokenID, err := suite.keeper.MintNFT(suite.ctx, denomID, tokenNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	response, err := suite.queryClient.NFT(gocontext.Background(), &types.QueryNFTRequest{
		DenomID: denomID,
		TokenId: tokenID,
	})

	suite.NoError(err)
	suite.NotEmpty(response.NFT)
	suite.Equal(response.NFT.Id, tokenID)
}

func (suite *IntegrationTestKeeperSuite) TestGetApprovalsNFT_Correctly_ReturnsApprovals() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address2)
	suite.NoError(err)

	tokenID, err := suite.keeper.MintNFT(suite.ctx, denomID, tokenNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	err = suite.keeper.AddApproval(suite.ctx, denomID, tokenID, address, address2)
	suite.NoError(err)

	response, err := suite.queryClient.GetApprovalsNFT(gocontext.Background(), &types.QueryApprovalsNFTRequest{
		DenomID: denomID,
		TokenId: tokenID,
	})

	suite.NoError(err)
	suite.Assert().Contains(response.ApprovedAddresses, address2.String())
}

func (suite *IntegrationTestKeeperSuite) TestIsApprovedOperator_ReturnsTrue_WhenOperatorIsApproved() {
	err := suite.keeper.AddApprovalForAll(suite.ctx, address, address2, true)
	suite.NoError(err)

	response, err := suite.queryClient.QueryApprovalsIsApprovedForAll(gocontext.Background(), &types.QueryApprovalsIsApprovedForAllRequest{
		Owner:    address.String(),
		Operator: address2.String(),
	})

	suite.NoError(err)
	suite.Equal(response.IsApproved, true)
}
