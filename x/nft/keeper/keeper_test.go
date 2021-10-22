package keeper_test

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"cudos.org/cudos-node/simapp"
	"cudos.org/cudos-node/x/nft/keeper"
	"cudos.org/cudos-node/x/nft/types"
)

var (
	denomID     = "denomid"
	denomNm     = "denomnm"
	denomSymbol = "denomSymbol"
	schema      = "{a:a,b:b}"

	denomID2     = "denomid2"
	denomNm2     = "denom2nm"
	denomSymbol2 = "denomSymbol2"

	tokenID  = "tokenid"
	tokenID2 = "tokenid2"
	tokenID3 = "tokenid3"

	tokenNm  = "tokennm"
	tokenNm2 = "tokennm2"
	tokenNm3 = "tokennm3"

	denomID3     = "denomid3"
	denomNm3     = "denom3nm"
	denomSymbol3 = "denomSymbol3"

	address    = CreateTestAddrs(1)[0]
	address2   = CreateTestAddrs(2)[1]
	address3   = CreateTestAddrs(3)[2]
	tokenURI   = "https://google.com/token-1.json"
	tokenURI2  = "https://google.com/token-2.json"
	tokenData  = "{a:a,b:b}"
	tokenData2 = "{a:a,b:b,c:c}"

	isCheckTx = false
)

type KeeperSuite struct {
	suite.Suite

	legacyAmino *codec.LegacyAmino
	ctx         sdk.Context
	keeper      keeper.Keeper
	app         *simapp.SimApp

	queryClient types.QueryClient
}

type MockedKeeper struct {
	mock.Mock
}

func (suite *KeeperSuite) SetupTest() {

	app := simapp.Setup(isCheckTx)

	suite.app = app
	suite.legacyAmino = app.LegacyAmino()
	suite.ctx = app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	suite.keeper = app.NftKeeper

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.NftKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)

}

func (suite *KeeperSuite) AfterTest(_, _ string) {
}

//TODO: Refactor with Should syntaxis
func (suite *KeeperSuite) TestIssueDenom() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, address)
	suite.NoError(err)

	// MintNFT shouldn't fail when collection does not exist
	err = suite.keeper.IssueDenom(suite.ctx, denomID2, denomNm2, schema, address)
	suite.NoError(err)

	err = suite.keeper.IssueDenom(suite.ctx, denomID3, denomNm3, schema, address3)
	suite.NoError(err)

	// collections should equal 3
	collections := suite.keeper.GetCollections(suite.ctx)
	suite.NotEmpty(collections)
	suite.Equal(len(collections), 3)
}

func (suite *KeeperSuite) TestMintNFT_ShouldError_WhenSenderIsNotDenomCreator() {

	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, address)
	suite.NoError(err)

	err = suite.keeper.MintNFT(suite.ctx, denomID, tokenID, denomNm, tokenURI, tokenData, address2, address)
	suite.ErrorIs(err, types.ErrUnauthorized)
}

func (suite *KeeperSuite) TestMintNFT_ShouldError_WhenDenomDoesNotExist() {
	err := suite.keeper.MintNFT(suite.ctx, denomID, tokenID, denomNm, tokenURI, tokenData, address2, address)
	suite.ErrorIs(err, types.ErrInvalidDenom)
}

func (suite *KeeperSuite) TestMintNFT_ShouldError_WhenNFTAlreadyExists() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, address)
	suite.NoError(err)

	err = suite.keeper.MintNFT(suite.ctx, denomID, tokenID, denomNm, tokenURI, tokenData, address, address2)
	suite.NoError(err)

	err = suite.keeper.MintNFT(suite.ctx, denomID, tokenID, denomNm, tokenURI, tokenData, address, address2)
	suite.ErrorIs(err, types.ErrNFTAlreadyExists)

}

func (suite *KeeperSuite) TestMintNFT_ShouldCorrectly_MintNewNFT() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, address)
	suite.NoError(err)

	err = suite.keeper.MintNFT(suite.ctx, denomID, tokenID, denomNm, tokenURI, tokenData, address, address2)
	suite.NoError(err)

	nftSuccessfullyMinted := suite.keeper.HasNFT(suite.ctx, denomID, tokenID)
	assert.Equal(suite.T(), true, nftSuccessfullyMinted)

}

func (suite *KeeperSuite) TestMintNFT_ShouldCorrectly_SetOwner() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, address)
	suite.NoError(err)

	err = suite.keeper.MintNFT(suite.ctx, denomID, tokenID, denomNm, tokenURI, tokenData, address, address2)
	suite.NoError(err)

	nft, err := suite.keeper.GetBaseNFT(suite.ctx, denomID, tokenID)
	suite.NoError(err)

	assert.Equal(suite.T(), nft.Owner, address2.String())

}

func (suite *KeeperSuite) TestMintNFT_ShouldCorrectly_IncreasesTotalSupply() {
	supplyBeforeMinting := suite.keeper.GetTotalSupply(suite.ctx, denomID)
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, address)
	suite.NoError(err)

	err = suite.keeper.MintNFT(suite.ctx, denomID, tokenID, denomNm, tokenURI, tokenData, address, address2)
	suite.NoError(err)
	supplyAfterMinting := suite.keeper.GetTotalSupply(suite.ctx, denomID)

	assert.Greater(suite.T(), supplyAfterMinting, supplyBeforeMinting)

}

func (suite *KeeperSuite) TestEditNFT_ShouldError_WhenDenomDoesNotExist() {
	err := suite.keeper.EditNFT(suite.ctx, denomID, tokenID, tokenNm, tokenURI, tokenData, address)
	suite.ErrorIs(err, types.ErrInvalidDenom)
}

func (suite *KeeperSuite) TestEditNFT_ShouldError_WhenSenderIsNotDenomCreator() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, address)
	suite.NoError(err)

	err = suite.keeper.MintNFT(suite.ctx, denomID, tokenID, denomNm, tokenURI, tokenData, address, address2)
	suite.NoError(err)

	err = suite.keeper.EditNFT(suite.ctx, denomID, tokenID, tokenNm, tokenURI, tokenData, address2)
	suite.ErrorIs(err, types.ErrUnauthorized)
}

func (suite *KeeperSuite) TestEditNFT_ShouldError_WhenNFTDoesNotExit() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, address)
	suite.NoError(err)

	err = suite.keeper.EditNFT(suite.ctx, denomID, tokenID, tokenNm, tokenURI, tokenData, address)
	suite.ErrorIs(err, types.ErrNotFoundNFT)
}

func (suite *KeeperSuite) TestEditNFT_ShouldError_WhenSenderIsNotOwner() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, address)
	suite.NoError(err)

	err = suite.keeper.MintNFT(suite.ctx, denomID, tokenID, denomNm, tokenURI, tokenData, address, address2)
	suite.NoError(err)

	err = suite.keeper.EditNFT(suite.ctx, denomID, tokenID, tokenNm, tokenURI, tokenData, address)
	suite.ErrorIs(err, types.ErrUnauthorized)
}

func (suite *KeeperSuite) TestEditNFT_ShouldCorrectly_UpdateNFTProperties() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, address2)
	suite.NoError(err)

	err = suite.keeper.MintNFT(suite.ctx, denomID, tokenID, denomNm, tokenURI, tokenData, address2, address2)
	suite.NoError(err)

	originalNFT, _ := suite.keeper.GetNFT(suite.ctx, denomID, tokenID)
	err = suite.keeper.EditNFT(suite.ctx, denomID, tokenID, tokenNm2, tokenURI2, tokenData2, address2)
	suite.NoError(err)

	editedNFT, _ := suite.keeper.GetNFT(suite.ctx, denomID, tokenID)

	assert.Equal(suite.T(), editedNFT.GetName(), tokenNm2)
	assert.Equal(suite.T(), editedNFT.GetData(), tokenData2)
	assert.Equal(suite.T(), editedNFT.GetURI(), tokenURI2)

	assert.NotEqual(suite.T(), originalNFT.GetName(), editedNFT.GetName())
	assert.NotEqual(suite.T(), originalNFT.GetData(), editedNFT.GetData())
	assert.NotEqual(suite.T(), originalNFT.GetURI(), editedNFT.GetURI())

}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(KeeperSuite))
}

// CreateTestAddrs creates test addresses
func CreateTestAddrs(numAddrs int) []sdk.AccAddress {
	var addresses []sdk.AccAddress
	var buffer bytes.Buffer

	// start at 100 so we can make up to 999 test addresses with valid test addresses
	for i := 100; i < (numAddrs + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString("A58856F0FD53BF058B4909A21AEC019107BA6") //base address string

		buffer.WriteString(numString) //adding on final two digits to make addresses unique
		res, _ := sdk.AccAddressFromHex(buffer.String())
		bech := res.String()
		addresses = append(addresses, testAddr(buffer.String(), bech))
		buffer.Reset()
	}

	return addresses
}

// for incode address generation
func testAddr(addr string, bech string) sdk.AccAddress {
	res, err := sdk.AccAddressFromHex(addr)
	if err != nil {
		panic(err)
	}
	bechexpected := res.String()
	if bech != bechexpected {
		panic("Bech encoding doesn't match reference")
	}

	bechres, err := sdk.AccAddressFromBech32(bech)
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(bechres, res) {
		panic("Bech decode and hex decode don't match")
	}

	return res
}
