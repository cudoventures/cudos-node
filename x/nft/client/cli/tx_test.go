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
	"github.com/CudoVentures/cudos-node/x/nft/client/cli"
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

func (s *TxTestSuite) TestNftCommands() {
	network, err := testutil.RunNetwork(s.T(), s.config)
	require.NoError(s.T(), err)

	clientCtx := network.Validators[0].ClientCtx
	valAddr := network.Validators[0].Address.String()

	nftRecieverAccAddrRecord, _, err := clientCtx.Keyring.NewMnemonic("nft_receiver", keyring.English, hd.CreateHDPath(118, 0, 0).String(), keyring.DefaultBIP39Passphrase, hd.Secp256k1)
	require.NoError(s.T(), err)
	nftRecieverAccAddr, err := nftRecieverAccAddrRecord.GetAddress()
	require.NoError(s.T(), err)
	nftRecieverAddr := nftRecieverAccAddr.String()

	nftTransferToAccAddrRecord, _, err := clientCtx.Keyring.NewMnemonic("nft_transfer_to", keyring.English, hd.CreateHDPath(118, 0, 0).String(), keyring.DefaultBIP39Passphrase, hd.Secp256k1)
	require.NoError(s.T(), err)
	nftTransferToAccAddr, err := nftTransferToAccAddrRecord.GetAddress()
	require.NoError(s.T(), err)
	nftTransferToAddr := nftTransferToAccAddr.String()

	denomId := "testdenom"
	tokenId := ""

	s.T().Run("IssueDenom", func(t *testing.T) {
		args := testutil.AppendDefaultTxFlags(valAddr, network.Config.BondDenom, []string{
			denomId,
			"--name=testdenom_name",
			"--symbol=testdenom_symbol",
			"--schema=testdenom_schema",
			// "--traits=NotEditable",
			"--description=testdenom_description",
			fmt.Sprintf("--minter=%s", valAddr),
			"--data=data",
		})

		txRes, txErr := clitestutil.ExecTestCLICmd(clientCtx, cli.GetCmdIssueDenom(), args)
		require.NoError(s.T(), txErr)
		testutil.WaitForBlock()

		txResp, txErr := testutil.QueryJustBroadcastedTx(clientCtx, txRes)
		require.NoError(t, txErr)
		require.Equal(t, uint32(0), txResp.Code)
	})

	s.T().Run("MintNFT", func(t *testing.T) {
		args := testutil.AppendDefaultTxFlags(valAddr, network.Config.BondDenom, []string{
			denomId,
			fmt.Sprintf("--recipient=%s", nftRecieverAddr),
			"--name=testnft_name",
			"--uri=testnft_uri",
		})

		txRes, txErr := clitestutil.ExecTestCLICmd(clientCtx, cli.GetCmdMintNFT(), args)
		require.NoError(s.T(), txErr)
		testutil.WaitForBlock()

		txResp, txErr := testutil.QueryJustBroadcastedTx(clientCtx, txRes)
		require.NoError(t, txErr)
		require.Equal(t, uint32(0), txResp.Code)

		tokenId, _ = testutil.GetEventValue(txResp, "mint_nft", "token_id")
		require.NotEqual(t, "", tokenId)
	})

	s.T().Run("EditNFT", func(t *testing.T) {
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
			denomId,
			tokenId,
			"--uri=testnft_uri_edited",
		})

		txRes, txErr := clitestutil.ExecTestCLICmd(clientCtx, cli.GetCmdEditNFT(), args)
		require.NoError(s.T(), txErr)
		testutil.WaitForBlock()

		txResp, txErr := testutil.QueryJustBroadcastedTx(clientCtx, txRes)
		require.NoError(t, txErr)
		require.Equal(t, uint32(0), txResp.Code)
	})

	s.T().Run("TransferNft", func(t *testing.T) {
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
			nftRecieverAddr,
			nftTransferToAddr,
			denomId,
			tokenId,
		})
		txRes, txErr := clitestutil.ExecTestCLICmd(clientCtx, cli.GetCmdTransferNft(), args)
		require.NoError(s.T(), txErr)
		testutil.WaitForBlock()

		txResp, txErr := testutil.QueryJustBroadcastedTx(clientCtx, txRes)
		require.NoError(t, txErr)
		require.Equal(t, uint32(0), txResp.Code)
	})

	s.T().Run("ApproveNft", func(t *testing.T) {
		// transfer enough funds to cover the tx costs
		args := testutil.AppendDefaultTxFlags(valAddr, network.Config.BondDenom, []string{
			valAddr,
			nftTransferToAddr,
			fmt.Sprintf("%d%s", testutil.TxFees, network.Config.BondDenom),
		})
		_, txErr := clitestutil.ExecTestCLICmd(clientCtx, bankcli.NewSendTxCmd(), args)
		require.NoError(s.T(), txErr)
		testutil.WaitForBlock()

		args = testutil.AppendDefaultTxFlags(nftTransferToAddr, network.Config.BondDenom, []string{
			nftRecieverAddr,
			denomId,
			tokenId,
		})
		txRes, txErr := clitestutil.ExecTestCLICmd(clientCtx, cli.GetCmdApproveNft(), args)
		require.NoError(s.T(), txErr)
		testutil.WaitForBlock()

		txResp, txErr := testutil.QueryJustBroadcastedTx(clientCtx, txRes)
		require.NoError(t, txErr)
		require.Equal(t, uint32(0), txResp.Code)
	})

	s.T().Run("ApproveAllNFT", func(t *testing.T) {
		// transfer enough funds to cover the tx costs
		args := testutil.AppendDefaultTxFlags(valAddr, network.Config.BondDenom, []string{
			valAddr,
			nftTransferToAddr,
			fmt.Sprintf("%d%s", testutil.TxFees, network.Config.BondDenom),
		})
		_, txErr := clitestutil.ExecTestCLICmd(clientCtx, bankcli.NewSendTxCmd(), args)
		require.NoError(s.T(), txErr)
		testutil.WaitForBlock()

		args = testutil.AppendDefaultTxFlags(nftTransferToAddr, network.Config.BondDenom, []string{
			nftRecieverAddr,
			"true",
		})
		txRes, txErr := clitestutil.ExecTestCLICmd(clientCtx, cli.GetCmdApproveAllNFT(), args)
		require.NoError(s.T(), txErr)
		testutil.WaitForBlock()

		txResp, txErr := testutil.QueryJustBroadcastedTx(clientCtx, txRes)
		require.NoError(t, txErr)
		require.Equal(t, uint32(0), txResp.Code)
	})

	s.T().Run("RevokeNft", func(t *testing.T) {
		// transfer enough funds to cover the tx costs
		args := testutil.AppendDefaultTxFlags(valAddr, network.Config.BondDenom, []string{
			valAddr,
			nftTransferToAddr,
			fmt.Sprintf("%d%s", testutil.TxFees, network.Config.BondDenom),
		})
		_, txErr := clitestutil.ExecTestCLICmd(clientCtx, bankcli.NewSendTxCmd(), args)
		require.NoError(s.T(), txErr)
		testutil.WaitForBlock()

		args = testutil.AppendDefaultTxFlags(nftTransferToAddr, network.Config.BondDenom, []string{
			nftRecieverAddr,
			denomId,
			tokenId,
		})
		txRes, txErr := clitestutil.ExecTestCLICmd(clientCtx, cli.GetCmdRevokeNft(), args)
		require.NoError(s.T(), txErr)
		testutil.WaitForBlock()

		txResp, txErr := testutil.QueryJustBroadcastedTx(clientCtx, txRes)
		require.NoError(t, txErr)
		require.Equal(t, uint32(0), txResp.Code)
	})

	s.T().Run("BurnNFT", func(t *testing.T) {
		// transfer enough funds to cover the tx costs
		args := testutil.AppendDefaultTxFlags(valAddr, network.Config.BondDenom, []string{
			valAddr,
			nftTransferToAddr,
			fmt.Sprintf("%d%s", testutil.TxFees, network.Config.BondDenom),
		})
		_, txErr := clitestutil.ExecTestCLICmd(clientCtx, bankcli.NewSendTxCmd(), args)
		require.NoError(s.T(), txErr)
		testutil.WaitForBlock()

		args = testutil.AppendDefaultTxFlags(nftTransferToAddr, network.Config.BondDenom, []string{
			denomId,
			tokenId,
		})
		txRes, txErr := clitestutil.ExecTestCLICmd(clientCtx, cli.GetCmdBurnNFT(), args)
		require.NoError(s.T(), txErr)
		testutil.WaitForBlock()

		txResp, txErr := testutil.QueryJustBroadcastedTx(clientCtx, txRes)
		require.NoError(t, txErr)
		require.Equal(t, uint32(0), txResp.Code)
	})

	s.T().Run("TransferDenom", func(t *testing.T) {
		args := testutil.AppendDefaultTxFlags(valAddr, network.Config.BondDenom, []string{
			nftTransferToAddr,
			denomId,
		})
		txRes, txErr := clitestutil.ExecTestCLICmd(clientCtx, cli.GetCmdTransferDenom(), args)
		require.NoError(s.T(), txErr)
		testutil.WaitForBlock()

		txResp, txErr := testutil.QueryJustBroadcastedTx(clientCtx, txRes)
		require.NoError(t, txErr)
		require.Equal(t, uint32(0), txResp.Code)
	})

	network.Cleanup()
}
