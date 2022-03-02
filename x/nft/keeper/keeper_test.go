package keeper_test

import (
	"bytes"
	"strconv"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/CudoVentures/cudos-node/simapp"
	"github.com/CudoVentures/cudos-node/x/nft/keeper"
	"github.com/CudoVentures/cudos-node/x/nft/types"
)

var (
	denomID     = "denomid"
	denomNm     = "denomnm"
	denomSymbol = "denomSymbol"
	schema      = "{a:a,b:b}"

	denomID2     = "denomid2"
	denomNm2     = "denom2nm"
	denomSymbol2 = "denom2Symbol"

	tokenNm  = "tokennm"
	tokenNm2 = "tokennm2"
	tokenNm3 = "tokennm3"

	denomID3     = "denomid3"
	denomNm3     = "denom3nm"
	denomSymbol3 = "denom3Symbol"

	address    = CreateTestAddrs(1)[0]
	address2   = CreateTestAddrs(2)[1]
	address3   = CreateTestAddrs(3)[2]
	tokenURI   = "https://google.com/token-1.json"
	tokenURI2  = "https://google.com/token-2.json"
	tokenData  = "{a:a,b:b}"
	tokenData2 = "{a:a,b:b,c:c}"

	isCheckTx = false
)

type IntegrationTestKeeperSuite struct {
	suite.Suite

	legacyAmino *codec.LegacyAmino
	ctx         sdk.Context
	keeper      keeper.Keeper
	app         *simapp.SimApp

	queryClient types.QueryClient
}

func (suite *IntegrationTestKeeperSuite) SetupTest() {

	app := simapp.Setup(isCheckTx)

	suite.app = app
	suite.legacyAmino = app.LegacyAmino()
	suite.ctx = app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	suite.keeper = app.NftKeeper

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.NftKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)

}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestKeeperSuite))
}

func (suite *IntegrationTestKeeperSuite) isApprovedAddress(nft *types.BaseNFT, sender *sdk.AccAddress) bool {
	for _, address := range nft.ApprovedAddresses {
		if address == sender.String() {
			return true
		}
	}
	return false
}

func (suite *IntegrationTestKeeperSuite) TestIssueDenom_ShouldError_WhenDenomIdAlreadyExists() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address)
	suite.NoError(err)

	err = suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address)
	suite.ErrorIs(err, types.ErrInvalidDenom)
}

func (suite *IntegrationTestKeeperSuite) TestIssueDenom_ShouldError_WhenDenomNameAlreadyExists() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID2, denomNm, denomSymbol, schema, address)
	suite.NoError(err)

	err = suite.keeper.IssueDenom(suite.ctx, denomID2, denomNm, denomSymbol, schema, address)
	suite.ErrorIs(err, types.ErrInvalidDenom)
}

func (suite *IntegrationTestKeeperSuite) TestIssueDenom_ShouldError_WhenDenomSymbolAlreadyExists() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID2, denomNm, denomSymbol, schema, address)
	suite.NoError(err)

	err = suite.keeper.IssueDenom(suite.ctx, denomID2, denomNm, denomSymbol, schema, address)
	suite.ErrorIs(err, types.ErrInvalidDenom)
}

func (suite *IntegrationTestKeeperSuite) TestIssueDenom_ShouldCorrectly_SetDenomIdAndNameAndSymbol() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, denomSymbol, schema, address)
	suite.NoError(err)
}

func (suite *IntegrationTestKeeperSuite) TestMintNFT_ShouldError_WhenSenderIsNotDenomCreator() {

	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address)
	suite.NoError(err)

	_, err = suite.keeper.MintNFT(suite.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	suite.ErrorIs(err, types.ErrUnauthorized)
}

func (suite *IntegrationTestKeeperSuite) TestMintNFT_ShouldError_WhenDenomDoesNotExist() {
	_, err := suite.keeper.MintNFT(suite.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	suite.ErrorIs(err, types.ErrInvalidDenom)
}

func (suite *IntegrationTestKeeperSuite) TestGetCollection_ShouldCorrectly_ReturnDenomCollections() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address)
	suite.NoError(err)

	// MintNFT shouldn't fail when collection does not exist
	err = suite.keeper.IssueDenom(suite.ctx, denomID2, denomNm2, schema, denomSymbol2, address)
	suite.NoError(err)

	err = suite.keeper.IssueDenom(suite.ctx, denomID3, denomNm3, schema, denomSymbol3, address3)
	suite.NoError(err)

	// collections should equal 3
	collections := suite.keeper.GetCollections(suite.ctx)
	suite.NotEmpty(collections)
	suite.Equal(len(collections), 3)
}

func (suite *IntegrationTestKeeperSuite) TestMintNFT_ShouldCorrectly_MintNewNFT() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address)
	suite.NoError(err)

	tokenId, err := suite.keeper.MintNFT(suite.ctx, denomID, denomNm, tokenURI, tokenData, address, address2)
	suite.NoError(err)

	nftSuccessfullyMinted := suite.keeper.HasNFT(suite.ctx, denomID, tokenId)
	assert.Equal(suite.T(), true, nftSuccessfullyMinted)

}

func (suite *IntegrationTestKeeperSuite) TestMintNFT_ShouldCorrectly_SetOwner() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address)
	suite.NoError(err)

	tokenId, err := suite.keeper.MintNFT(suite.ctx, denomID, denomNm, tokenURI, tokenData, address, address2)
	suite.NoError(err)

	nft, err := suite.keeper.GetBaseNFT(suite.ctx, denomID, tokenId)
	suite.NoError(err)

	assert.Equal(suite.T(), nft.Owner, address2.String())

	owner, err := suite.keeper.GetOwner(suite.ctx, address2, denomID)
	isOwnerCorrectlySavedInDb := false

	for _, collection := range owner.IDCollections {
		if collection.DenomId == denomID {
			for _, ownedTokenId := range collection.TokenIds {
				if ownedTokenId == tokenId {
					isOwnerCorrectlySavedInDb = true
				}
			}
		}
	}

	assert.Equal(suite.T(), true, isOwnerCorrectlySavedInDb)

}

func (suite *IntegrationTestKeeperSuite) TestMintNFT_ShouldCorrectly_IncreasesTotalSupply() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address)
	suite.NoError(err)
	supplyBeforeMinting := suite.keeper.GetTotalSupply(suite.ctx, denomID)

	_, err = suite.keeper.MintNFT(suite.ctx, denomID, denomNm, tokenURI, tokenData, address, address2)
	suite.NoError(err)
	supplyAfterMinting := suite.keeper.GetTotalSupply(suite.ctx, denomID)

	assert.Greater(suite.T(), supplyAfterMinting, supplyBeforeMinting)

}

// TODO:
// test total count function

func (suite *IntegrationTestKeeperSuite) TestEditNFT_ShouldError_WhenDenomDoesNotExist() {
	err := suite.keeper.EditNFT(suite.ctx, denomID, "1234", tokenNm, tokenURI, tokenData, address)
	suite.ErrorIs(err, types.ErrInvalidDenom)
}

func (suite *IntegrationTestKeeperSuite) TestEditNFT_ShouldError_WhenNFTDoesNotExit() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address)
	suite.NoError(err)

	err = suite.keeper.EditNFT(suite.ctx, denomID, "1234", tokenNm, tokenURI, tokenData, address)
	suite.ErrorIs(err, types.ErrNotFoundNFT)
}

func (suite *IntegrationTestKeeperSuite) TestEditNFT_ShouldError_WhenSenderIsNotOwner() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address)
	suite.NoError(err)

	tokenId, err := suite.keeper.MintNFT(suite.ctx, denomID, denomNm, tokenURI, tokenData, address, address2)
	suite.NoError(err)

	err = suite.keeper.EditNFT(suite.ctx, denomID, tokenId, tokenNm, tokenURI, tokenData, address)
	suite.ErrorIs(err, types.ErrUnauthorized)
}

func (suite *IntegrationTestKeeperSuite) TestEditNFT_ShouldCorrectly_UpdateNFTProperties() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address2)
	suite.NoError(err)

	tokenId, err := suite.keeper.MintNFT(suite.ctx, denomID, denomNm, tokenURI, tokenData, address2, address2)
	suite.NoError(err)

	originalNFT, _ := suite.keeper.GetNFT(suite.ctx, denomID, tokenId)
	err = suite.keeper.EditNFT(suite.ctx, denomID, tokenId, tokenNm2, tokenURI2, tokenData2, address2)
	suite.NoError(err)

	editedNFT, _ := suite.keeper.GetNFT(suite.ctx, denomID, tokenId)

	assert.Equal(suite.T(), editedNFT.GetName(), tokenNm2)
	assert.Equal(suite.T(), editedNFT.GetData(), tokenData2)
	assert.Equal(suite.T(), editedNFT.GetURI(), tokenURI2)

	assert.NotEqual(suite.T(), originalNFT.GetName(), editedNFT.GetName())
	assert.NotEqual(suite.T(), originalNFT.GetData(), editedNFT.GetData())
	assert.NotEqual(suite.T(), originalNFT.GetURI(), editedNFT.GetURI())

}

func (suite *IntegrationTestKeeperSuite) TestTransferOwner_ShouldError_WhenDenomDoesNotExist() {
	err := suite.keeper.TransferOwner(suite.ctx, denomID, "1234", address, address2, address3)
	suite.ErrorIs(err, types.ErrInvalidDenom)
}

func (suite *IntegrationTestKeeperSuite) TestTransferOwner_ShouldError_WhenNFTDoesNotBelongToFromAddress() {

	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address2)
	suite.NoError(err)

	tokenId, err := suite.keeper.MintNFT(suite.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	err = suite.keeper.TransferOwner(suite.ctx, denomID, tokenId, address3, address2, address2)
	suite.ErrorIs(err, types.ErrUnauthorized)
}

func (suite *IntegrationTestKeeperSuite) TestTransferOwner_ShouldError_WhenSenderDoesNotHavePermissionForTransfer() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address2)
	suite.NoError(err)

	tokenId, err := suite.keeper.MintNFT(suite.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	err = suite.keeper.TransferOwner(suite.ctx, denomID, tokenId, address, address2, address2)
	suite.ErrorIs(err, types.ErrUnauthorized)
}

func (suite *IntegrationTestKeeperSuite) TestTransferOwner_ShouldCorrectly_TransferWhenSenderIsOwner() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address2)
	suite.NoError(err)

	tokenId, err := suite.keeper.MintNFT(suite.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	err = suite.keeper.TransferOwner(suite.ctx, denomID, tokenId, address, address2, address)
	suite.NoError(err)
}

func (suite *IntegrationTestKeeperSuite) TestTransferOwner_ShouldCorrectly_TransferWhenSenderIsApprovedOnNFT() {

	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address2)
	suite.NoError(err)

	tokenId, err := suite.keeper.MintNFT(suite.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	err = suite.keeper.AddApproval(suite.ctx, denomID, tokenId, address, address3)
	suite.NoError(err)

	err = suite.keeper.TransferOwner(suite.ctx, denomID, tokenId, address, address2, address3)
	suite.NoError(err)

	nft, err := suite.keeper.GetBaseNFT(suite.ctx, denomID, tokenId)
	assert.Equal(suite.T(), nft.Owner, address2.String())

}

func (suite *IntegrationTestKeeperSuite) TestTransferOwner_ShouldCorrectly_TransferWhenSenderIsApprovedOperatorAllForNFTOwner() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address2)
	suite.NoError(err)

	tokenId, err := suite.keeper.MintNFT(suite.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	err = suite.keeper.AddApprovalForAll(suite.ctx, address, address3, true)
	suite.NoError(err)

	err = suite.keeper.TransferOwner(suite.ctx, denomID, tokenId, address, address2, address3)
	suite.NoError(err)
}

func (suite *IntegrationTestKeeperSuite) TestTransferOwner_ShouldCorrectly_SwapOwner() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address2)
	suite.NoError(err)

	tokenId, err := suite.keeper.MintNFT(suite.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	err = suite.keeper.TransferOwner(suite.ctx, denomID, tokenId, address, address2, address)
	suite.NoError(err)

	nft, err := suite.keeper.GetNFT(suite.ctx, denomID, tokenId)
	suite.NoError(err)

	assert.Equal(suite.T(), nft.GetOwner().String(), address2.String())

	owner, err := suite.keeper.GetOwner(suite.ctx, address2, denomID)
	isOwnerCorrectlySwappedInDb := false

	for _, collection := range owner.IDCollections {
		if collection.DenomId == denomID {
			for _, ownedTokenId := range collection.TokenIds {
				if ownedTokenId == tokenId {
					isOwnerCorrectlySwappedInDb = true
				}
			}
		}
	}

	assert.Equal(suite.T(), true, isOwnerCorrectlySwappedInDb)

}

func (suite *IntegrationTestKeeperSuite) TestAddApproval_ShouldError_WhenSenderIsNotOwnerOfNftOrIsNotApproved() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address2)
	suite.NoError(err)

	tokenId, err := suite.keeper.MintNFT(suite.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	err = suite.keeper.AddApproval(suite.ctx, denomID, tokenId, address2, address2)
	suite.ErrorIs(err, types.ErrUnauthorized)
}

func (suite *IntegrationTestKeeperSuite) TestAddApproval_ShouldCorrectly_AddAddressToNFTApprovedList() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address2)
	suite.NoError(err)

	tokenId, err := suite.keeper.MintNFT(suite.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	err = suite.keeper.AddApproval(suite.ctx, denomID, tokenId, address, address2)
	suite.NoError(err)

	nft, _ := suite.keeper.GetBaseNFT(suite.ctx, denomID, tokenId)
	isApproved := suite.isApprovedAddress(&nft, &address2)
	assert.Equal(suite.T(), isApproved, true)
}

func (suite *IntegrationTestKeeperSuite) TestAddApprovalAll_ShouldError_WhenSenderAddressIsTheSameAsApproved() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address2)
	suite.NoError(err)

	_, err = suite.keeper.MintNFT(suite.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	err = suite.keeper.AddApprovalForAll(suite.ctx, address2, address2, true)
	suite.ErrorIs(err, sdkerrors.ErrInvalidRequest)
}

func (suite *IntegrationTestKeeperSuite) TestAddApprovalAll_ShouldCorrectly_AddAddressToNFTApprovedList() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address2)
	suite.NoError(err)

	_, err = suite.keeper.MintNFT(suite.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	err = suite.keeper.AddApprovalForAll(suite.ctx, address, address2, true)
	suite.NoError(err)

	approvedAddresses, _ := suite.keeper.GetApprovedAddresses(suite.ctx, address)
	isApproved := approvedAddresses.ApprovedAddressesData[address2.String()]
	assert.Equal(suite.T(), isApproved, true)
}

func (suite *IntegrationTestKeeperSuite) TestRevokeApproval_ShouldError_WhenSenderIsNotOwnerOrApprovedOperator() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address2)
	suite.NoError(err)

	tokenId, err := suite.keeper.MintNFT(suite.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	err = suite.keeper.RevokeApproval(suite.ctx, denomID, tokenId, address2, address2)
	suite.ErrorIs(err, types.ErrUnauthorized)
}

func (suite *IntegrationTestKeeperSuite) TestRevokeApproval_ShouldCorrectly_RevokeNFTApproval() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address2)
	suite.NoError(err)

	tokenId, err := suite.keeper.MintNFT(suite.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	err = suite.keeper.AddApproval(suite.ctx, denomID, tokenId, address, address2)
	suite.NoError(err)

	nft, _ := suite.keeper.GetBaseNFT(suite.ctx, denomID, tokenId)
	isApproved := suite.isApprovedAddress(&nft, &address2)
	assert.Equal(suite.T(), isApproved, true)

	err = suite.keeper.RevokeApproval(suite.ctx, denomID, tokenId, address, address2)

	nft, _ = suite.keeper.GetBaseNFT(suite.ctx, denomID, tokenId)
	isApproved = suite.isApprovedAddress(&nft, &address2)
	assert.Equal(suite.T(), isApproved, false)
}

func (suite *IntegrationTestKeeperSuite) TestBurnNFT_ShouldError_WhenDenomIdDoesNotExist() {
	err := suite.keeper.BurnNFT(suite.ctx, denomID, "1234", address)
	suite.ErrorIs(err, types.ErrInvalidDenom)
}

func (suite *IntegrationTestKeeperSuite) TestBurnNFT_ShouldError_WhenSenderIsNotOwner() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address2)
	suite.NoError(err)

	tokenId, err := suite.keeper.MintNFT(suite.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	err = suite.keeper.BurnNFT(suite.ctx, denomID, tokenId, address2)
	suite.ErrorIs(err, types.ErrUnauthorized)
}

func (suite *IntegrationTestKeeperSuite) TestBurnNFT_ShouldCorrectly_DeleteNFT() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address2)
	suite.NoError(err)

	tokenId, err := suite.keeper.MintNFT(suite.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	nft, err := suite.keeper.GetBaseNFT(suite.ctx, denomID, tokenId)
	suite.NoError(err)
	assert.NotNil(suite.T(), nft)

	err = suite.keeper.BurnNFT(suite.ctx, denomID, tokenId, address)
	suite.NoError(err)

	_, err = suite.keeper.GetBaseNFT(suite.ctx, denomID, tokenId)
	suite.ErrorIs(err, types.ErrNotFoundNFT)
}

func (suite *IntegrationTestKeeperSuite) TestBurnNFT_ShouldCorrectly_DeleteNFTOwner() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address2)
	suite.NoError(err)

	tokenId, err := suite.keeper.MintNFT(suite.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	nft, err := suite.keeper.GetBaseNFT(suite.ctx, denomID, tokenId)
	suite.NoError(err)
	assert.NotNil(suite.T(), nft)

	err = suite.keeper.BurnNFT(suite.ctx, denomID, tokenId, address)
	suite.NoError(err)

	_, err = suite.keeper.GetBaseNFT(suite.ctx, denomID, tokenId)
	suite.ErrorIs(err, types.ErrNotFoundNFT)

	owner, err := suite.keeper.GetOwner(suite.ctx, address, denomID)
	isOwnerCorrectlySwappedInDb := false

	for _, collection := range owner.IDCollections {
		if collection.DenomId == denomID {
			for _, ownedTokenId := range collection.TokenIds {
				if ownedTokenId == tokenId {
					isOwnerCorrectlySwappedInDb = true
				}
			}
		}
	}

	assert.Equal(suite.T(), false, isOwnerCorrectlySwappedInDb)

}

func (suite *IntegrationTestKeeperSuite) TestBurnNFT_ShouldCorrectly_DecreaseSupply() {
	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address2)
	suite.NoError(err)

	tokenId, err := suite.keeper.MintNFT(suite.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	suite.NoError(err)

	nft, err := suite.keeper.GetBaseNFT(suite.ctx, denomID, tokenId)
	suite.NoError(err, types.ErrNotFoundNFT)
	assert.NotNil(suite.T(), nft)

	err = suite.keeper.BurnNFT(suite.ctx, denomID, tokenId, address)
	suite.NoError(err)

	supplyAfterBurn := suite.keeper.GetTotalSupply(suite.ctx, denomID)

	assert.Equal(suite.T(), uint64(0), supplyAfterBurn)
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
