package keeper_test

import (
	"cudos.org/cudos-node/x/nft/exported"
	"cudos.org/cudos-node/x/nft/types"
	"encoding/binary"

	keep "cudos.org/cudos-node/x/nft/keeper"
	abci "github.com/tendermint/tendermint/abci/types"
)

func (suite *IntegrationTestKeeperSuite) TestNewQuerier_ReturnsCorrect() {
	querier := keep.NewQuerier(suite.keeper, suite.legacyAmino)
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	_, err := querier(suite.ctx, []string{"foo", "bar"}, query)
	suite.Error(err)
}

func (suite *IntegrationTestKeeperSuite) TestQuerySupply_ReturnsCorrectSupply() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, address2)
	suite.NoError(err)

	_, err = suite.keeper.MintNFT(suite.ctx, denomID, tokenNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	querier := keep.NewQuerier(suite.keeper, suite.legacyAmino)

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	query.Path = "/custom/nft/supply"
	query.Data = []byte("?")

	res, err := querier(suite.ctx, []string{"supply"}, query)
	suite.Error(err)
	suite.Nil(res)

	queryCollectionParams := types.NewQuerySupplyParams(denomID2, nil)
	bz, errRes := suite.legacyAmino.MarshalJSON(queryCollectionParams)
	suite.Nil(errRes)
	query.Data = bz
	res, err = querier(suite.ctx, []string{"supply"}, query)
	suite.NoError(err)
	supplyResp := binary.LittleEndian.Uint64(res)
	suite.Equal(0, int(supplyResp))

	queryCollectionParams = types.NewQuerySupplyParams(denomID, nil)
	bz, errRes = suite.legacyAmino.MarshalJSON(queryCollectionParams)
	suite.Nil(errRes)
	query.Data = bz

	res, err = querier(suite.ctx, []string{"supply"}, query)
	suite.NoError(err)
	suite.NotNil(res)

	supplyResp = binary.LittleEndian.Uint64(res)
	suite.Equal(1, int(supplyResp))
}

func (suite *IntegrationTestKeeperSuite) TestQueryCollection_ReturnsCorrectCollection() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, address2)
	suite.NoError(err)

	err = suite.keeper.IssueDenom(suite.ctx, denomID2, denomNm2, schema, address2)
	suite.NoError(err)

	_, err = suite.keeper.MintNFT(suite.ctx, denomID, tokenNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	querier := keep.NewQuerier(suite.keeper, suite.legacyAmino)

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	query.Path = "/custom/nft/collection"

	query.Data = []byte("?")
	res, err := querier(suite.ctx, []string{"collection"}, query)
	suite.Error(err)
	suite.Nil(res)

	queryCollectionParams := types.NewQuerySupplyParams(denomID2, nil)
	bz, errRes := suite.legacyAmino.MarshalJSON(queryCollectionParams)
	suite.Nil(errRes)

	query.Data = bz
	_, err = querier(suite.ctx, []string{"collection"}, query)
	suite.NoError(err)

	queryCollectionParams = types.NewQuerySupplyParams(denomID, nil)
	bz, errRes = suite.legacyAmino.MarshalJSON(queryCollectionParams)
	suite.Nil(errRes)

	query.Data = bz
	res, err = querier(suite.ctx, []string{"collection"}, query)
	suite.NoError(err)
	suite.NotNil(res)

	var collection types.Collection
	types.ModuleCdc.MustUnmarshalJSON(res, &collection)
	suite.Len(collection.NFTs, 1)
}

func (suite *IntegrationTestKeeperSuite) TestQueryOwner_ReturnsCorrectOwner() {

	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, address2)
	suite.NoError(err)
	err = suite.keeper.IssueDenom(suite.ctx, denomID2, denomNm2, schema, address2)
	suite.NoError(err)

	tokenId, err := suite.keeper.MintNFT(suite.ctx, denomID, tokenNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)
	tokenId, err = suite.keeper.MintNFT(suite.ctx, denomID2, tokenNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	querier := keep.NewQuerier(suite.keeper, suite.legacyAmino)
	query := abci.RequestQuery{
		Path: "/custom/nft/owner",
		Data: []byte{},
	}

	query.Data = []byte("?")
	_, err = querier(suite.ctx, []string{"owner"}, query)
	suite.Error(err)

	// query the balance using no denomID so that all denoms will be returns
	params := types.NewQuerySupplyParams("", address)
	bz, err2 := suite.legacyAmino.MarshalJSON(params)
	suite.Nil(err2)
	query.Data = bz

	var out types.Owner
	res, err := querier(suite.ctx, []string{"owner"}, query)
	suite.NoError(err)
	suite.NotNil(res)

	suite.legacyAmino.MustUnmarshalJSON(res, &out)

	// build the owner using both denoms
	idCollection1 := types.NewIDCollection(denomID, []string{tokenId})
	idCollection2 := types.NewIDCollection(denomID2, []string{tokenId})
	owner := types.NewOwner(address, idCollection1, idCollection2)

	suite.EqualValues(out.String(), owner.String())
}

func (suite *IntegrationTestKeeperSuite) TestQueryNFT_ReturnsCorrectNFT() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, address2)
	suite.NoError(err)

	tokenId, err := suite.keeper.MintNFT(suite.ctx, denomID, tokenNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	querier := keep.NewQuerier(suite.keeper, suite.legacyAmino)

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	query.Path = "/custom/nft/nft"
	var res []byte

	query.Data = []byte("?")
	res, err = querier(suite.ctx, []string{"nft"}, query)
	suite.Error(err)
	suite.Nil(res)

	params := types.NewQueryNFTParams(denomID2, "1234")
	bz, err2 := suite.legacyAmino.MarshalJSON(params)
	suite.Nil(err2)

	query.Data = bz
	res, err = querier(suite.ctx, []string{"nft"}, query)
	suite.Error(err)
	suite.Nil(res)

	params = types.NewQueryNFTParams(denomID, tokenId)
	bz, err2 = suite.legacyAmino.MarshalJSON(params)
	suite.Nil(err2)

	query.Data = bz
	res, err = querier(suite.ctx, []string{"nft"}, query)
	suite.NoError(err)
	suite.NotNil(res)

	var out exported.NFT
	suite.legacyAmino.MustUnmarshalJSON(res, &out)

	suite.Equal(out.GetID(), tokenId)
	suite.Equal(out.GetURI(), tokenURI)
	suite.Equal(out.GetOwner(), address)
}

func (suite *IntegrationTestKeeperSuite) TestQueryDenoms_ReturnsCorrectDenoms() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, address2)
	suite.NoError(err)

	_, err = suite.keeper.MintNFT(suite.ctx, denomID, tokenNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	err = suite.keeper.IssueDenom(suite.ctx, denomID2, denomNm2, schema, address2)
	suite.NoError(err)

	_, err = suite.keeper.MintNFT(suite.ctx, denomID2, tokenNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	querier := keep.NewQuerier(suite.keeper, suite.legacyAmino)

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	var res []byte
	query.Path = "/custom/nft/denoms"

	res, err = querier(suite.ctx, []string{"denoms"}, query)
	suite.NoError(err)
	suite.NotNil(res)

	denoms := []string{denomID, denomID2, denomID3}

	var out []types.Denom
	suite.legacyAmino.MustUnmarshalJSON(res, &out)

	for key, denomInQuestion := range out {
		suite.Equal(denomInQuestion.Id, denoms[key])
	}
}
