package cli_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/CudoVentures/cudos-node/simapp"
	"github.com/CudoVentures/cudos-node/testutil"
	"github.com/CudoVentures/cudos-node/testutil/sample"
	"github.com/CudoVentures/cudos-node/x/marketplace/client/cli"
	nftcli "github.com/CudoVentures/cudos-node/x/nft/client/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	bankcli "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
)

type TxTestSuite struct {
	suite.Suite
	config network.Config
}

func TestTxTestSuite(t *testing.T) {
	suite.Run(t, new(TxTestSuite))
}

func (s *TxTestSuite) SetupSuite() {
	s.T().Log("setting up tx test suite")

	s.config = simapp.NewConfig(s.T().TempDir())
	s.config.NumValidators = 1
}

func (s *TxTestSuite) TearDownSuite() {
	s.T().Log("tearing down tx test suite")
}

func (s *TxTestSuite) TestMarketplaceCommands() {
	network, err := testutil.RunNetwork(s.T(), s.config)
	require.NoError(s.T(), err)

	clientCtx := network.Validators[0].ClientCtx
	valAddr := network.Validators[0].Address.String()

	nftRecieverAccAddrRecord, _, err := clientCtx.Keyring.NewMnemonic("nft_receiver", keyring.English, hd.CreateHDPath(118, 0, 0).String(), keyring.DefaultBIP39Passphrase, hd.Secp256k1)
	require.NoError(s.T(), err)
	nftRecieverAccAddr, err := nftRecieverAccAddrRecord.GetAddress()
	require.NoError(s.T(), err)
	nftRecieverAddr := nftRecieverAccAddr.String()

	accountAddress := sample.AccAddress()

	denomId := "testdenom"
	tokenId := ""
	collectionId := ""
	nftId := ""

	s.T().Run("AddAdmin", func(t *testing.T) {
		args := testutil.AppendDefaultTxFlags(valAddr, network.Config.BondDenom, []string{
			valAddr,
		})

		txRes, txErr := clitestutil.ExecTestCLICmd(clientCtx, cli.CmdAddAdmin(), args)
		require.NoError(s.T(), txErr)
		testutil.WaitForBlock()

		txResp, txErr := testutil.QueryJustBroadcastedTx(clientCtx, txRes)
		require.NoError(t, txErr)
		require.Equal(t, uint32(0), txResp.Code)
	})

	s.T().Run("CreateCollection", func(t *testing.T) {
		args := testutil.AppendDefaultTxFlags(valAddr, network.Config.BondDenom, []string{
			denomId,
			"--name=testdenom_name",
			"--symbol=testdenom_symbol",
			"--schema=testdenom_schema",
			"--description=testdenom_description",
			fmt.Sprintf("--minter=%s", valAddr),
			"--data=data",
			fmt.Sprintf("--mint-royalties=%s:100", valAddr),
			fmt.Sprintf("--resale-royalties=%s:100", valAddr),
			"--verified=true",
		})

		txRes, txErr := clitestutil.ExecTestCLICmd(clientCtx, cli.CmdCreateCollection(), args)
		require.NoError(s.T(), txErr)
		testutil.WaitForBlock()

		txResp, txErr := testutil.QueryJustBroadcastedTx(clientCtx, txRes)
		require.NoError(t, txErr)
		require.Equal(t, uint32(0), txResp.Code)

		collectionId, _ = testutil.GetEventValue(txResp, "create_collection", "collection_id")
		require.NotEqual(t, "", collectionId)
	})

	s.T().Run("UnverifyCollection", func(t *testing.T) {
		args := testutil.AppendDefaultTxFlags(valAddr, network.Config.BondDenom, []string{
			collectionId,
		})

		txRes, txErr := clitestutil.ExecTestCLICmd(clientCtx, cli.CmdUnverifyCollection(), args)
		require.NoError(s.T(), txErr)
		testutil.WaitForBlock()

		txResp, txErr := testutil.QueryJustBroadcastedTx(clientCtx, txRes)
		require.NoError(t, txErr)
		require.Equal(t, uint32(0), txResp.Code)
	})

	s.T().Run("VerifyCollection", func(t *testing.T) {
		args := testutil.AppendDefaultTxFlags(valAddr, network.Config.BondDenom, []string{
			collectionId,
		})

		txRes, txErr := clitestutil.ExecTestCLICmd(clientCtx, cli.CmdVerifyCollection(), args)
		require.NoError(s.T(), txErr)
		testutil.WaitForBlock()

		txResp, txErr := testutil.QueryJustBroadcastedTx(clientCtx, txRes)
		require.NoError(t, txErr)
		require.Equal(t, uint32(0), txResp.Code)
	})

	s.T().Run("UpdateRoyalties", func(t *testing.T) {
		args := testutil.AppendDefaultTxFlags(valAddr, network.Config.BondDenom, []string{
			collectionId,
			fmt.Sprintf("--mint-royalties=%s:99,%s:1", valAddr, accountAddress),
			fmt.Sprintf("--resale-royalties=%s:99,%s:1", valAddr, accountAddress),
		})

		txRes, txErr := clitestutil.ExecTestCLICmd(clientCtx, cli.CmdUpdateRoyalties(), args)
		require.NoError(s.T(), txErr)
		testutil.WaitForBlock()

		txResp, txErr := testutil.QueryJustBroadcastedTx(clientCtx, txRes)
		require.NoError(t, txErr)
		require.Equal(t, uint32(0), txResp.Code)
	})

	s.T().Run("MintNft", func(t *testing.T) {
		args := testutil.AppendDefaultTxFlags(valAddr, network.Config.BondDenom, []string{
			denomId,
			nftRecieverAddr,
			fmt.Sprintf("1%s", network.Config.BondDenom),
			"nftforsale1_name",
		})

		txRes, txErr := clitestutil.ExecTestCLICmd(clientCtx, cli.CmdMintNft(), args)
		require.NoError(s.T(), txErr)
		testutil.WaitForBlock()

		txResp, txErr := testutil.QueryJustBroadcastedTx(clientCtx, txRes)
		require.NoError(t, txErr)
		require.Equal(t, uint32(0), txResp.Code)

		tokenId, _ = testutil.GetEventValue(txResp, "marketplace_mint_nft", "token_id")
		require.NotEqual(t, "", tokenId)
	})

	s.T().Run("PublishNft", func(t *testing.T) {
		// transfer enough funds to cover the tx costs
		args := testutil.AppendDefaultTxFlags(valAddr, network.Config.BondDenom, []string{
			valAddr,
			nftRecieverAddr,
			fmt.Sprintf("%d%s", testutil.TxFees, network.Config.BondDenom),
		})

		_, txErr := clitestutil.ExecTestCLICmd(clientCtx, bankcli.NewSendTxCmd(), args)
		require.NoError(s.T(), txErr)
		testutil.WaitForBlock()

		args = testutil.AppendDefaultTxFlags(nftRecieverAddr, network.Config.BondDenom, []string{
			tokenId,
			denomId,
			fmt.Sprintf("1%s", network.Config.BondDenom),
		})

		txRes, txErr := clitestutil.ExecTestCLICmd(clientCtx, cli.CmdPublishNft(), args)
		require.NoError(s.T(), txErr)
		testutil.WaitForBlock()

		txResp, txErr := testutil.QueryJustBroadcastedTx(clientCtx, txRes)
		require.NoError(t, txErr)
		require.Equal(t, uint32(0), txResp.Code)

		nftId, _ = testutil.GetEventValue(txResp, "publish_nft", "nft_id")
		require.NotEqual(t, "", nftId)
	})

	s.T().Run("BuyNft", func(t *testing.T) {
		args := testutil.AppendDefaultTxFlags(valAddr, network.Config.BondDenom, []string{
			nftId,
		})

		txRes, txErr := clitestutil.ExecTestCLICmd(clientCtx, cli.CmdBuyNft(), args)
		require.NoError(s.T(), txErr)
		testutil.WaitForBlock()

		txResp, txErr := testutil.QueryJustBroadcastedTx(clientCtx, txRes)
		require.NoError(t, txErr)
		require.Equal(t, uint32(0), txResp.Code)
	})

	// publish it again so we can check "UpdatePrice" and "RemoveNft"
	s.T().Run("PublishNft", func(t *testing.T) {
		args := testutil.AppendDefaultTxFlags(valAddr, network.Config.BondDenom, []string{
			tokenId,
			denomId,
			fmt.Sprintf("1%s", network.Config.BondDenom),
		})

		txRes, txErr := clitestutil.ExecTestCLICmd(clientCtx, cli.CmdPublishNft(), args)
		require.NoError(s.T(), txErr)
		testutil.WaitForBlock()

		txResp, txErr := testutil.QueryJustBroadcastedTx(clientCtx, txRes)
		require.NoError(t, txErr)
		require.Equal(t, uint32(0), txResp.Code)

		nftId, _ = testutil.GetEventValue(txResp, "publish_nft", "nft_id")
		require.NotEqual(t, "", nftId)
	})

	s.T().Run("UpdatePrice", func(t *testing.T) {
		args := testutil.AppendDefaultTxFlags(valAddr, network.Config.BondDenom, []string{
			nftId,
			fmt.Sprintf("2%s", network.Config.BondDenom),
		})

		txRes, txErr := clitestutil.ExecTestCLICmd(clientCtx, cli.CmdUpdatePrice(), args)
		require.NoError(s.T(), txErr)
		testutil.WaitForBlock()

		txResp, txErr := testutil.QueryJustBroadcastedTx(clientCtx, txRes)
		require.NoError(t, txErr)
		require.Equal(t, uint32(0), txResp.Code)
	})

	s.T().Run("RemoveNft", func(t *testing.T) {
		args := testutil.AppendDefaultTxFlags(valAddr, network.Config.BondDenom, []string{
			nftId,
		})

		txRes, txErr := clitestutil.ExecTestCLICmd(clientCtx, cli.CmdRemoveNft(), args)
		require.NoError(s.T(), txErr)
		testutil.WaitForBlock()

		txResp, txErr := testutil.QueryJustBroadcastedTx(clientCtx, txRes)
		require.NoError(t, txErr)
		require.Equal(t, uint32(0), txResp.Code)
	})

	s.T().Run("PublishCollection", func(t *testing.T) {
		publishDenomId := "testdenom2"

		// issue denom from nft module
		args := testutil.AppendDefaultTxFlags(valAddr, network.Config.BondDenom, []string{
			publishDenomId,
			"--name=testdenom2_name",
			"--symbol=testdenom2_symbol",
			"--schema=testdenom2_schema",
			"--description=testdenom2_description",
			fmt.Sprintf("--minter=%s", valAddr),
			"--data=data",
		})

		txRes, txErr := clitestutil.ExecTestCLICmd(clientCtx, nftcli.GetCmdIssueDenom(), args)
		require.NoError(s.T(), txErr)
		testutil.WaitForBlock()

		txResp, txErr := testutil.QueryJustBroadcastedTx(clientCtx, txRes)
		require.NoError(t, txErr)
		require.Equal(t, uint32(0), txResp.Code)

		// publish the denomo from marketplace module
		args = testutil.AppendDefaultTxFlags(valAddr, network.Config.BondDenom, []string{
			publishDenomId,
			fmt.Sprintf("--mint-royalties=%s:100", valAddr),
			fmt.Sprintf("--resale-royalties=%s:100", valAddr),
		})

		txRes, txErr = clitestutil.ExecTestCLICmd(clientCtx, cli.CmdPublishCollection(), args)
		require.NoError(s.T(), txErr)
		testutil.WaitForBlock()

		txResp, txErr = testutil.QueryJustBroadcastedTx(clientCtx, txRes)
		require.NoError(t, txErr)
		require.Equal(t, uint32(0), txResp.Code)
	})

	s.T().Run("RemoveAdmin", func(t *testing.T) {
		args := testutil.AppendDefaultTxFlags(valAddr, network.Config.BondDenom, []string{
			valAddr,
		})

		txRes, txErr := clitestutil.ExecTestCLICmd(clientCtx, cli.CmdRemoveAdmin(), args)
		require.NoError(s.T(), txErr)
		testutil.WaitForBlock()

		txResp, txErr := testutil.QueryJustBroadcastedTx(clientCtx, txRes)
		require.NoError(t, txErr)
		require.Equal(t, uint32(0), txResp.Code)
	})

	network.Cleanup()
}
