package keeper_test

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/CudoVentures/cudos-node/simapp"
	"github.com/CudoVentures/cudos-node/x/nft/keeper"
	"github.com/CudoVentures/cudos-node/x/nft/types"
)

//nolint:gosec // these are not hard-coded credentials
var (
	denomID          = "denomid"
	denomNm          = "denomnm"
	denomSymbol      = "denomSymbol"
	schema           = "{a:a,b:b}"
	denomTraits      = ""
	denomMinter      = ""
	denomDescription = "denom Description"
	denomData        = "somedata"

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

func (s *IntegrationTestKeeperSuite) SetupTest() {
	app := simapp.Setup(isCheckTx)

	s.app = app
	s.legacyAmino = app.LegacyAmino()
	s.ctx = app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	s.keeper = app.NftKeeper

	queryHelper := baseapp.NewQueryServerTestHelper(s.ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.NftKeeper)
	s.queryClient = types.NewQueryClient(queryHelper)
}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestKeeperSuite))
}

func (s *IntegrationTestKeeperSuite) TestIssueDenom_ShouldError_WhenDenomIDAlreadyExists() {
	err := s.keeper.IssueDenom(s.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address)
	s.NoError(err)

	err = s.keeper.IssueDenom(s.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address)
	s.ErrorIs(err, types.ErrInvalidDenom)
}

func (s *IntegrationTestKeeperSuite) TestIssueDenom_ShouldError_WhenDenomNameAlreadyExists() {
	err := s.keeper.IssueDenom(s.ctx, denomID2, denomNm, denomSymbol, schema, denomTraits, denomMinter, denomDescription, denomData, address)
	s.NoError(err)

	err = s.keeper.IssueDenom(s.ctx, denomID2, denomNm, denomSymbol, schema, denomTraits, denomMinter, denomDescription, denomData, address)
	s.ErrorIs(err, types.ErrInvalidDenom)
}

func (s *IntegrationTestKeeperSuite) TestIssueDenom_ShouldError_WhenDenomSymbolAlreadyExists() {
	err := s.keeper.IssueDenom(s.ctx, denomID2, denomNm, denomSymbol, schema, denomTraits, denomMinter, denomDescription, denomData, address)
	s.NoError(err)

	err = s.keeper.IssueDenom(s.ctx, denomID2, denomNm, denomSymbol, schema, denomTraits, denomMinter, denomDescription, denomData, address)
	s.ErrorIs(err, types.ErrInvalidDenom)
}

func (s *IntegrationTestKeeperSuite) TestIssueDenom_ShouldCorrectly_SetDenomIDAndNameAndSymbol() {
	err := s.keeper.IssueDenom(s.ctx, denomID, denomNm, denomSymbol, schema, denomTraits, denomMinter, denomDescription, denomData, address)
	s.NoError(err)
}

func (s *IntegrationTestKeeperSuite) TestMintNFT_ShouldError_WhenSenderIsNotDenomCreator() {
	err := s.keeper.IssueDenom(s.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address)
	s.NoError(err)

	_, err = s.keeper.MintNFT(s.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	s.ErrorIs(err, types.ErrUnauthorized)
}

func (s *IntegrationTestKeeperSuite) TestMintNFT_ShouldError_WhenDenomDoesNotExist() {
	_, err := s.keeper.MintNFT(s.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	s.ErrorIs(err, types.ErrInvalidDenom)
}

func (s *IntegrationTestKeeperSuite) TestGetCollection_ShouldCorrectly_ReturnDenomCollections() {
	err := s.keeper.IssueDenom(s.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address)
	s.NoError(err)

	// MintNFT shouldn't fail when collection does not exist
	err = s.keeper.IssueDenom(s.ctx, denomID2, denomNm2, schema, denomSymbol2, denomTraits, denomMinter, denomDescription, denomData, address)
	s.NoError(err)

	err = s.keeper.IssueDenom(s.ctx, denomID3, denomNm3, schema, denomSymbol3, denomTraits, denomMinter, denomDescription, denomData, address3)
	s.NoError(err)

	// collections should equal 3
	collections := s.keeper.GetCollections(s.ctx)
	s.NotEmpty(collections)
	s.Equal(len(collections), 3)
}

func (s *IntegrationTestKeeperSuite) TestMintNFT_ShouldCorrectly_MintNewNFT() {
	err := s.keeper.IssueDenom(s.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address)
	s.NoError(err)

	tokenID, err := s.keeper.MintNFT(s.ctx, denomID, denomNm, tokenURI, tokenData, address, address2)
	s.NoError(err)

	nftSuccessfullyMinted := s.keeper.HasNFT(s.ctx, denomID, tokenID)
	assert.Equal(s.T(), true, nftSuccessfullyMinted)
}

func (s *IntegrationTestKeeperSuite) TestMintNFT_ShouldCorrectly_SetOwner() {
	err := s.keeper.IssueDenom(s.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address)
	s.NoError(err)

	tokenID, err := s.keeper.MintNFT(s.ctx, denomID, denomNm, tokenURI, tokenData, address, address2)
	s.NoError(err)

	nft, err := s.keeper.GetBaseNFT(s.ctx, denomID, tokenID)
	s.NoError(err)

	assert.Equal(s.T(), nft.Owner, address2.String())

	owner, err := s.keeper.GetOwner(s.ctx, address2, denomID)
	isOwnerCorrectlySavedInDb := false

	for _, collection := range owner.IDCollections {
		if collection.DenomID == denomID {
			for _, ownedTokenId := range collection.TokenIds {
				if ownedTokenId == tokenID {
					isOwnerCorrectlySavedInDb = true
				}
			}
		}
	}

	assert.Equal(s.T(), true, isOwnerCorrectlySavedInDb)
}

func (s *IntegrationTestKeeperSuite) TestMintNFT_ShouldCorrectly_IncreasesTotalSupply() {
	err := s.keeper.IssueDenom(s.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address)
	s.NoError(err)
	supplyBeforeMinting := s.keeper.GetTotalSupply(s.ctx, denomID)

	_, err = s.keeper.MintNFT(s.ctx, denomID, denomNm, tokenURI, tokenData, address, address2)
	s.NoError(err)
	supplyAfterMinting := s.keeper.GetTotalSupply(s.ctx, denomID)

	assert.Greater(s.T(), supplyAfterMinting, supplyBeforeMinting)
}

// TODO:
// test total count function

func (s *IntegrationTestKeeperSuite) TestEditNFT_ShouldError_WhenDenomDoesNotExist() {
	err := s.keeper.EditNFT(s.ctx, denomID, "1234", tokenNm, tokenURI, tokenData, address)
	s.ErrorIs(err, types.ErrInvalidDenom)
}

func (s *IntegrationTestKeeperSuite) TestEditNFT_ShouldError_WhenNFTDoesNotExit() {
	err := s.keeper.IssueDenom(s.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address)
	s.NoError(err)

	err = s.keeper.EditNFT(s.ctx, denomID, "1234", tokenNm, tokenURI, tokenData, address)
	s.ErrorIs(err, types.ErrNotFoundNFT)
}

func (s *IntegrationTestKeeperSuite) TestEditNFT_ShouldError_WhenSenderIsNotOwner() {
	err := s.keeper.IssueDenom(s.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address)
	s.NoError(err)

	tokenID, err := s.keeper.MintNFT(s.ctx, denomID, denomNm, tokenURI, tokenData, address, address2)
	s.NoError(err)

	err = s.keeper.EditNFT(s.ctx, denomID, tokenID, tokenNm, tokenURI, tokenData, address)
	s.ErrorIs(err, types.ErrUnauthorized)
}

func (s *IntegrationTestKeeperSuite) TestEditNFT_ShouldCorrectly_UpdateNFTProperties() {
	err := s.keeper.IssueDenom(s.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address2)
	s.NoError(err)

	tokenID, err := s.keeper.MintNFT(s.ctx, denomID, denomNm, tokenURI, tokenData, address2, address2)
	s.NoError(err)

	originalNFT, _ := s.keeper.GetNFT(s.ctx, denomID, tokenID)
	err = s.keeper.EditNFT(s.ctx, denomID, tokenID, tokenNm2, tokenURI2, tokenData2, address2)
	s.NoError(err)

	editedNFT, _ := s.keeper.GetNFT(s.ctx, denomID, tokenID)

	assert.Equal(s.T(), editedNFT.GetName(), tokenNm2)
	assert.Equal(s.T(), editedNFT.GetData(), tokenData2)
	assert.Equal(s.T(), editedNFT.GetURI(), tokenURI2)

	assert.NotEqual(s.T(), originalNFT.GetName(), editedNFT.GetName())
	assert.NotEqual(s.T(), originalNFT.GetData(), editedNFT.GetData())
	assert.NotEqual(s.T(), originalNFT.GetURI(), editedNFT.GetURI())
}

func (s *IntegrationTestKeeperSuite) TestTransferOwner_ShouldError_WhenDenomDoesNotExist() {
	err := s.keeper.TransferOwner(s.ctx, denomID, "1234", address, address2, address3)
	s.ErrorIs(err, types.ErrInvalidDenom)
}

func (s *IntegrationTestKeeperSuite) TestTransferOwner_ShouldError_WhenNFTDoesNotBelongToFromAddress() {
	err := s.keeper.IssueDenom(s.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address2)
	s.NoError(err)

	tokenID, err := s.keeper.MintNFT(s.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	s.NoError(err)

	err = s.keeper.TransferOwner(s.ctx, denomID, tokenID, address3, address2, address2)
	s.ErrorIs(err, types.ErrUnauthorized)
}

func (s *IntegrationTestKeeperSuite) TestTransferOwner_ShouldError_WhenSenderDoesNotHavePermissionForTransfer() {
	err := s.keeper.IssueDenom(s.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address2)
	s.NoError(err)

	tokenID, err := s.keeper.MintNFT(s.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	s.NoError(err)

	err = s.keeper.TransferOwner(s.ctx, denomID, tokenID, address, address2, address2)
	s.ErrorIs(err, types.ErrUnauthorized)
}

func (s *IntegrationTestKeeperSuite) TestTransferOwner_ShouldCorrectly_TransferWhenSenderIsOwner() {
	err := s.keeper.IssueDenom(s.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address2)
	s.NoError(err)

	tokenID, err := s.keeper.MintNFT(s.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	s.NoError(err)

	err = s.keeper.TransferOwner(s.ctx, denomID, tokenID, address, address2, address)
	s.NoError(err)
}

func (s *IntegrationTestKeeperSuite) TestTransferOwner_ShouldCorrectly_TransferWhenSenderIsApprovedOnNFT() {
	err := s.keeper.IssueDenom(s.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address2)
	s.NoError(err)

	tokenID, err := s.keeper.MintNFT(s.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	s.NoError(err)

	err = s.keeper.AddApproval(s.ctx, denomID, tokenID, address, address3)
	s.NoError(err)

	err = s.keeper.TransferOwner(s.ctx, denomID, tokenID, address, address2, address3)
	s.NoError(err)

	nft, err := s.keeper.GetBaseNFT(s.ctx, denomID, tokenID)
	assert.Equal(s.T(), nft.Owner, address2.String())
}

func (s *IntegrationTestKeeperSuite) TestTransferOwner_ShouldCorrectly_TransferWhenSenderIsApprovedOperatorAllForNFTOwner() {
	err := s.keeper.IssueDenom(s.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address2)
	s.NoError(err)

	tokenID, err := s.keeper.MintNFT(s.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	s.NoError(err)

	err = s.keeper.AddApprovalForAll(s.ctx, address, address3, true)
	s.NoError(err)

	err = s.keeper.TransferOwner(s.ctx, denomID, tokenID, address, address2, address3)
	s.NoError(err)
}

func (s *IntegrationTestKeeperSuite) TestTransferOwner_ShouldCorrectly_SwapOwner() {
	err := s.keeper.IssueDenom(s.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address2)
	s.NoError(err)

	tokenID, err := s.keeper.MintNFT(s.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	s.NoError(err)

	err = s.keeper.TransferOwner(s.ctx, denomID, tokenID, address, address2, address)
	s.NoError(err)

	nft, err := s.keeper.GetNFT(s.ctx, denomID, tokenID)
	s.NoError(err)

	assert.Equal(s.T(), nft.GetOwner().String(), address2.String())

	owner, err := s.keeper.GetOwner(s.ctx, address2, denomID)
	isOwnerCorrectlySwappedInDb := false

	for _, collection := range owner.IDCollections {
		if collection.DenomID == denomID {
			for _, ownedTokenId := range collection.TokenIds {
				if ownedTokenId == tokenID {
					isOwnerCorrectlySwappedInDb = true
				}
			}
		}
	}

	assert.Equal(s.T(), true, isOwnerCorrectlySwappedInDb)
}

func (s *IntegrationTestKeeperSuite) TestAddApproval_ShouldError_WhenSenderIsNotOwnerOfNftOrIsNotApproved() {
	err := s.keeper.IssueDenom(s.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address2)
	s.NoError(err)

	tokenID, err := s.keeper.MintNFT(s.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	s.NoError(err)

	err = s.keeper.AddApproval(s.ctx, denomID, tokenID, address2, address2)
	s.ErrorIs(err, types.ErrUnauthorized)
}

func (s *IntegrationTestKeeperSuite) TestAddApproval_ShouldCorrectly_AddAddressToNFTApprovedList() {
	err := s.keeper.IssueDenom(s.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address2)
	s.NoError(err)

	tokenID, err := s.keeper.MintNFT(s.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	s.NoError(err)

	err = s.keeper.AddApproval(s.ctx, denomID, tokenID, address, address2)
	s.NoError(err)

	nft, _ := s.keeper.GetBaseNFT(s.ctx, denomID, tokenID)
	isApproved := s.keeper.IsApprovedAddress(&nft, address2.String())
	assert.Equal(s.T(), isApproved, true)
}

func (s *IntegrationTestKeeperSuite) TestAddApprovalAll_ShouldError_WhenSenderAddressIsTheSameAsApproved() {
	err := s.keeper.IssueDenom(s.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address2)
	s.NoError(err)

	_, err = s.keeper.MintNFT(s.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	s.NoError(err)

	err = s.keeper.AddApprovalForAll(s.ctx, address2, address2, true)
	s.ErrorIs(err, sdkerrors.ErrInvalidRequest)
}

func (s *IntegrationTestKeeperSuite) TestAddApprovalAll_ShouldCorrectly_AddAddressToNFTApprovedList() {
	err := s.keeper.IssueDenom(s.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address2)
	s.NoError(err)

	_, err = s.keeper.MintNFT(s.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	s.NoError(err)

	err = s.keeper.AddApprovalForAll(s.ctx, address, address2, true)
	s.NoError(err)

	approvedAddresses, _ := s.keeper.GetApprovedAddresses(s.ctx, address)
	isApproved := approvedAddresses.ApprovedAddressesData[address2.String()]
	assert.Equal(s.T(), isApproved, true)
}

func (s *IntegrationTestKeeperSuite) TestRevokeApproval_ShouldError_WhenSenderIsNotOwnerOrApprovedOperator() {
	err := s.keeper.IssueDenom(s.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address2)
	s.NoError(err)

	tokenID, err := s.keeper.MintNFT(s.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	s.NoError(err)

	err = s.keeper.RevokeApproval(s.ctx, denomID, tokenID, address2, address2)
	s.ErrorIs(err, types.ErrUnauthorized)
}

func (s *IntegrationTestKeeperSuite) TestRevokeApproval_ShouldCorrectly_RevokeNFTApproval() {
	err := s.keeper.IssueDenom(s.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address2)
	s.NoError(err)

	tokenID, err := s.keeper.MintNFT(s.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	s.NoError(err)

	err = s.keeper.AddApproval(s.ctx, denomID, tokenID, address, address2)
	s.NoError(err)

	nft, _ := s.keeper.GetBaseNFT(s.ctx, denomID, tokenID)
	isApproved := s.keeper.IsApprovedAddress(&nft, address2.String())
	assert.Equal(s.T(), isApproved, true)

	err = s.keeper.RevokeApproval(s.ctx, denomID, tokenID, address, address2)

	nft, _ = s.keeper.GetBaseNFT(s.ctx, denomID, tokenID)
	isApproved = s.keeper.IsApprovedAddress(&nft, address2.String())
	assert.Equal(s.T(), isApproved, false)
}

func (s *IntegrationTestKeeperSuite) TestTransferDenom() {
	err := s.keeper.IssueDenom(s.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address)
	s.NoError(err)

	// invalid owner
	err = s.keeper.TransferDenomOwner(s.ctx, denomID, address3, address2)
	s.Error(err)

	// right
	err = s.keeper.TransferDenomOwner(s.ctx, denomID, address, address3)
	s.NoError(err)

	denom, _ := s.keeper.GetDenom(s.ctx, denomID)

	// denom.Creator should equal to address3 after transfer
	s.Equal(denom.Creator, address3.String())
}

func (s *IntegrationTestKeeperSuite) TestBurnNFT_ShouldError_WhenDenomIDDoesNotExist() {
	err := s.keeper.BurnNFT(s.ctx, denomID, "1234", address)
	s.ErrorIs(err, types.ErrInvalidDenom)
}

func (s *IntegrationTestKeeperSuite) TestBurnNFT_ShouldError_WhenSenderIsNotOwner() {
	err := s.keeper.IssueDenom(s.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address2)
	s.NoError(err)

	tokenID, err := s.keeper.MintNFT(s.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	s.NoError(err)

	err = s.keeper.BurnNFT(s.ctx, denomID, tokenID, address2)
	s.ErrorIs(err, types.ErrUnauthorized)
}

func (s *IntegrationTestKeeperSuite) TestBurnNFT_ShouldCorrectly_DeleteNFT() {
	err := s.keeper.IssueDenom(s.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address2)
	s.NoError(err)

	tokenID, err := s.keeper.MintNFT(s.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	s.NoError(err)

	nft, err := s.keeper.GetBaseNFT(s.ctx, denomID, tokenID)
	s.NoError(err)
	assert.NotNil(s.T(), nft)

	err = s.keeper.BurnNFT(s.ctx, denomID, tokenID, address)
	s.NoError(err)

	_, err = s.keeper.GetBaseNFT(s.ctx, denomID, tokenID)
	s.ErrorIs(err, types.ErrNotFoundNFT)
}

func (s *IntegrationTestKeeperSuite) TestBurnNFT_ShouldCorrectly_DeleteNFTOwner() {
	err := s.keeper.IssueDenom(s.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address2)
	s.NoError(err)

	tokenID, err := s.keeper.MintNFT(s.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	s.NoError(err)

	nft, err := s.keeper.GetBaseNFT(s.ctx, denomID, tokenID)
	s.NoError(err)
	assert.NotNil(s.T(), nft)

	err = s.keeper.BurnNFT(s.ctx, denomID, tokenID, address)
	s.NoError(err)

	_, err = s.keeper.GetBaseNFT(s.ctx, denomID, tokenID)
	s.ErrorIs(err, types.ErrNotFoundNFT)

	owner, err := s.keeper.GetOwner(s.ctx, address, denomID)
	isOwnerCorrectlySwappedInDb := false

	for _, collection := range owner.IDCollections {
		if collection.DenomID == denomID {
			for _, ownedTokenId := range collection.TokenIds {
				if ownedTokenId == tokenID {
					isOwnerCorrectlySwappedInDb = true
				}
			}
		}
	}

	assert.Equal(s.T(), false, isOwnerCorrectlySwappedInDb)
}

func (s *IntegrationTestKeeperSuite) TestBurnNFT_ShouldCorrectly_DecreaseSupply() {
	err := s.keeper.IssueDenom(s.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address2)
	s.NoError(err)

	tokenID, err := s.keeper.MintNFT(s.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	s.NoError(err)

	nft, err := s.keeper.GetBaseNFT(s.ctx, denomID, tokenID)
	s.NoError(err, types.ErrNotFoundNFT)
	assert.NotNil(s.T(), nft)

	err = s.keeper.BurnNFT(s.ctx, denomID, tokenID, address)
	s.NoError(err)

	supplyAfterBurn := s.keeper.GetTotalSupply(s.ctx, denomID)

	assert.Equal(s.T(), uint64(0), supplyAfterBurn)
}

func (s *IntegrationTestKeeperSuite) TestDenom_With_NotEditable_Trait_NftsShouldNotBeEditable() {
	err := s.keeper.IssueDenom(s.ctx, denomID, denomNm, schema, denomSymbol, "NotEditable", denomMinter, denomDescription, denomData, address2)
	s.NoError(err)

	tokenID, err := s.keeper.MintNFT(s.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	s.NoError(err)

	err = s.keeper.BurnNFT(s.ctx, denomID, tokenID, address)
	s.Equal("denom 'denomid' has not editable trait: not editable", err.Error())

	err = s.keeper.EditNFT(s.ctx, denomID, tokenID, tokenNm, "", "", address)
	s.Equal("denom 'denomid' has not editable trait: not editable", err.Error())
}

func (s *IntegrationTestKeeperSuite) TestDenom_With_Minter_ShouldAllowMinterToMintNfts() {
	err := s.keeper.IssueDenom(s.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, address.String(), denomDescription, denomData, address2)
	s.NoError(err)

	_, err = s.keeper.MintNFT(s.ctx, denomID, denomNm, tokenURI, tokenData, address, address)
	s.NoError(err)
}

func (s *IntegrationTestKeeperSuite) TestSoftLockedNftShouldBeNotTransferableBurnableEditable() {
	err := s.keeper.IssueDenom(s.ctx, denomID, denomNm, schema, denomSymbol, denomTraits, denomMinter, denomDescription, denomData, address2)
	s.NoError(err)

	tokenID, err := s.keeper.MintNFT(s.ctx, denomID, denomNm, tokenURI, tokenData, address2, address)
	s.NoError(err)

	lockOwner := "lockOwner"
	s.NoError(s.keeper.SoftLockNFT(s.ctx, lockOwner, denomID, tokenID))

	err = s.keeper.BurnNFT(s.ctx, denomID, tokenID, address)
	s.Equal("token id 1 from denom with id denomid is soft locked by lockOwner: soft locked", err.Error())

	err = s.keeper.EditNFT(s.ctx, denomID, tokenID, tokenNm, "", "", address)
	s.Equal("token id 1 from denom with id denomid is soft locked by lockOwner: soft locked", err.Error())

	err = s.keeper.TransferOwner(s.ctx, denomID, tokenID, address, address2, address)
	s.Equal("token id 1 from denom with id denomid is soft locked by lockOwner: soft locked", err.Error())

	s.NoError(s.keeper.SoftUnlockNFT(s.ctx, lockOwner, denomID, tokenID))

	s.NoError(s.keeper.EditNFT(s.ctx, denomID, tokenID, tokenNm, "", "", address))
}

// CreateTestAddrs creates test addresses
func CreateTestAddrs(numAddrs int) []sdk.AccAddress {
	var addresses []sdk.AccAddress
	var buffer bytes.Buffer

	// start at 100 so we can make up to 999 test addresses with valid test addresses
	for i := 100; i < (numAddrs + 100); i++ {
		numString := strconv.Itoa(i)

		buffer.WriteString(numString) // adding on final two digits to make addresses unique
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
