package keeper_test

import (
	gocontext "context"

	"github.com/CudoVentures/cudos-node/x/nft/types"
)

func (suite *IntegrationTestKeeperSuite) TestSupplyReturnsCorrect() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address2)
	suite.NoError(err)

	_, err = suite.keeper.MintNFT(suite.ctx, denomID, tokenNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	response, err := suite.queryClient.Supply(gocontext.Background(), &types.QuerySupplyRequest{
		DenomId: denomID,
		Owner:   address.String(),
	})

	suite.NoError(err)
	suite.Equal(1, int(response.Amount))
}

func (suite *IntegrationTestKeeperSuite) TestOwner_ReturnsCorrect() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address2)
	suite.NoError(err)

	tokenID, err := suite.keeper.MintNFT(suite.ctx, denomID, tokenNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	response, err := suite.queryClient.Owner(gocontext.Background(), &types.QueryOwnerRequest{
		DenomId: denomID,
		Owner:   address.String(),
	})

	suite.NoError(err)
	suite.NotNil(response.Owner)
	suite.Contains(response.Owner.IDCollections[0].TokenIds, tokenID)
}

func (suite *IntegrationTestKeeperSuite) TestCollection_ReturnsCorrect() {

	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address2)
	suite.NoError(err)

	tokenId, err := suite.keeper.MintNFT(suite.ctx, denomID, tokenNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	response, err := suite.queryClient.Collection(gocontext.Background(), &types.QueryCollectionRequest{
		DenomId: denomID,
	})

	suite.NoError(err)
	suite.NotNil(response.Collection)
	suite.Len(response.Collection.NFTs, 1)
	suite.Equal(response.Collection.NFTs[0].Id, tokenId)
}

func (suite *IntegrationTestKeeperSuite) TestDenom_ReturnsCorrect() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address2)
	suite.NoError(err)

	_, err = suite.keeper.MintNFT(suite.ctx, denomID, tokenNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	response, err := suite.queryClient.Denom(gocontext.Background(), &types.QueryDenomRequest{
		DenomId: denomID,
	})

	suite.NoError(err)
	suite.NotNil(response.Denom)
	suite.Equal(response.Denom.Id, denomID)
}

func (suite *IntegrationTestKeeperSuite) TestDenoms_ReturnsCorrectCollection() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address2)
	suite.NoError(err)

	_, err = suite.keeper.MintNFT(suite.ctx, denomID, tokenNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	response, err := suite.queryClient.Denoms(gocontext.Background(), &types.QueryDenomsRequest{})

	suite.NoError(err)
	suite.NotEmpty(response.Denoms)
	suite.Equal(response.Denoms[0].Id, denomID)
}

func (suite *IntegrationTestKeeperSuite) TestNFT_ReturnsCorrect() {

	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address2)
	suite.NoError(err)

	tokenId, err := suite.keeper.MintNFT(suite.ctx, denomID, tokenNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	response, err := suite.queryClient.NFT(gocontext.Background(), &types.QueryNFTRequest{
		DenomId: denomID,
		TokenId: tokenId,
	})

	suite.NoError(err)
	suite.NotEmpty(response.NFT)
	suite.Equal(response.NFT.Id, tokenId)
}

func (suite *IntegrationTestKeeperSuite) TestGetApprovalsNFT_Correctly_ReturnsApprovals() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address2)
	suite.NoError(err)

	tokenId, err := suite.keeper.MintNFT(suite.ctx, denomID, tokenNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	err = suite.keeper.AddApproval(suite.ctx, denomID, tokenId, address, address2)
	suite.NoError(err)

	response, err := suite.queryClient.GetApprovalsNFT(gocontext.Background(), &types.QueryApprovalsNFTRequest{
		DenomId: denomID,
		TokenId: tokenId,
	})

	suite.NoError(err)
	suite.Equal(response.ApprovedAddresses[address2.String()], true)
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
