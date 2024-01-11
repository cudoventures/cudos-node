package apptesting

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/CudoVentures/cudos-node/app"
	appparams "github.com/CudoVentures/cudos-node/app/params"
	cudoMintKeeper "github.com/CudoVentures/cudos-node/x/cudoMint/keeper"
	cudoMinttypes "github.com/CudoVentures/cudos-node/x/cudoMint/types"
	gravitytypes "github.com/althea-net/cosmos-gravity-bridge/module/x/gravity/types"
	"github.com/cosmos/cosmos-sdk/baseapp"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	adminkeeper "github.com/CudoVentures/cudos-node/x/admin/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

// SimAppChainID hardcoded chainID for simulation
const (
	SimAppChainID        = "cudos-app"
	AccountAddressPrefix = "cudos"
	DefaultEthAddress    = "0x4838B106FCe9647Bdf1E7877BF73cE8B0BAD5f97"
	CudosMainnetEthAddr  = "0x817bbDbC3e8A1204f3691d14bB44992841e3dB35"
)

var (
	DefaultPowerReduction, _ = sdk.NewIntFromString("1000000000000000000")
	MinSelfDelegation, _     = sdk.NewIntFromString("2000000000000000000000000")
)

var (
	AccountPubKeyPrefix    = AccountAddressPrefix + "pub"
	ValidatorAddressPrefix = AccountAddressPrefix + "valoper"
	ValidatorPubKeyPrefix  = AccountAddressPrefix + "valoperpub"
	ConsNodeAddressPrefix  = AccountAddressPrefix + "valcons"
	ConsNodePubKeyPrefix   = AccountAddressPrefix + "valconspub"
)

type KeeperTestHelper struct {
	suite.Suite

	App         *app.App
	Ctx         sdk.Context
	CheckCtx    sdk.Context
	QueryHelper *baseapp.QueryServiceTestHelper
	TestAccs    []sdk.AccAddress

	CudoMintKeeper cudoMintKeeper.Keeper
	AdminKeeper    adminkeeper.Keeper
}

func (s *KeeperTestHelper) Setup(_ *testing.T, chainID string) {
	s.App = SetupApp(s.T())
	s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: chainID, Time: time.Now().UTC()})
	s.CheckCtx = s.App.BaseApp.NewContext(true, tmproto.Header{Height: 1, ChainID: chainID, Time: time.Now().UTC()})
	s.QueryHelper = &baseapp.QueryServiceTestHelper{
		GRPCQueryRouter: s.App.GRPCQueryRouter(),
		Ctx:             s.Ctx,
	}

	s.TestAccs = s.RandomAccountAddresses(3)

	// CudoMintKeeper is already initialized after SetupApp, but it is set as private field, so we need to set it again in the test helper for testing
	keyCudoMint := sdk.NewKVStoreKey("someKey")
	memstoreCudoMint := sdk.NewKVStoreKey("anotherKey")
	s.CudoMintKeeper = *cudoMintKeeper.NewKeeper(s.App.AppCodec(), keyCudoMint, memstoreCudoMint, s.App.BankKeeper, s.App.AccountKeeper, s.App.ParamsKeeper.Subspace("cudomint2"), authtypes.FeeCollectorName)

	// AdminKeeper is already initialized after SetupApp, but it is set as private field, so we need to set it again in the test helper for testing
	keyAdmin := sdk.NewKVStoreKey("someKey1")
	memstoreAdmin := sdk.NewKVStoreKey("anotherKey1")
	s.AdminKeeper = *adminkeeper.NewKeeper(s.App.AppCodec(), keyAdmin, memstoreAdmin, s.App.DistrKeeper, s.App.BankKeeper)
}

// DefaultConsensusParams defines the default Tendermint consensus params used
// in app testing.
var DefaultConsensusParams = &abci.ConsensusParams{
	Block: &abci.BlockParams{
		MaxBytes: 200000,
		MaxGas:   2000000,
	},
	Evidence: &tmproto.EvidenceParams{
		MaxAgeNumBlocks: 302400,
		MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration
		MaxBytes:        10000,
	},
	Validator: &tmproto.ValidatorParams{
		PubKeyTypes: []string{
			tmtypes.ABCIPubKeyTypeEd25519,
		},
	},
}

func SetupApp(t *testing.T) *app.App {
	t.Helper()
	privVal := NewPV()
	pubKey, err := privVal.GetPubKey()
	require.NoError(t, err)
	// create validator set with single validator
	validator := tmtypes.NewValidator(pubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	// generate genesis account
	senderPrivKey := secp256k1.GenPrivKey()
	acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(appparams.BondDenom, MinSelfDelegation.Mul(sdk.NewIntFromUint64(10)))),
	}
	return SetupWithGenesisValSet(t, valSet, []authtypes.GenesisAccount{acc}, balance)
}

// SetupWithGenesisValSet initializes a new app with a validator set and genesis accounts
// that also act as delegators. For simplicity, each validator is bonded with a delegation
// of one consensus engine unit in the default token of the app from first genesis
// account. A Nop logger is set in app.
func SetupWithGenesisValSet(t *testing.T, valSet *tmtypes.ValidatorSet, genAccs []authtypes.GenesisAccount, balances ...banktypes.Balance) *app.App {
	t.Helper()

	cudosApp, genesisState := setup(true, 5)
	genesisState = genesisStateWithValSet(t, cudosApp, genesisState, valSet, genAccs, balances...)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")

	require.NoError(t, err)
	// init chain will set the validator set and initialize the genesis accounts

	cudosApp.InitChain(
		abci.RequestInitChain{
			ChainId:         SimAppChainID,
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)

	// commit genesis changes
	cudosApp.Commit()
	cudosApp.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{
		ChainID:            SimAppChainID,
		Height:             cudosApp.LastBlockHeight() + 1,
		AppHash:            cudosApp.LastCommitID().Hash,
		ValidatorsHash:     valSet.Hash(),
		NextValidatorsHash: valSet.Hash(),
	}})

	return cudosApp
}

func setup(withGenesis bool, invCheckPeriod uint) (*app.App, app.GenesisState) {
	db := dbm.NewMemDB()
	encCdc := app.MakeEncodingConfig()

	cudosApp := app.New(
		log.NewNopLogger(),
		db,
		nil,
		true,
		map[int64]bool{},
		app.DefaultNodeHome,
		invCheckPeriod,
		encCdc,
		simapp.EmptyAppOptions{},
	)
	if withGenesis {
		return cudosApp, app.NewDefaultGenesisState(encCdc.Codec)
	}

	return cudosApp, app.GenesisState{}
}

func genesisStateWithValSet(t *testing.T,
	app *app.App, genesisState app.GenesisState,
	valSet *tmtypes.ValidatorSet, genAccs []authtypes.GenesisAccount,
	balances ...banktypes.Balance,
) app.GenesisState {
	// set genesis accounts
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = app.AppCodec().MustMarshalJSON(authGenesis)
	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

	bondAmt := MinSelfDelegation
	for _, val := range valSet.Validators {
		pk, err := cryptocodec.FromTmPubKeyInterface(val.PubKey)
		require.NoError(t, err)
		pkAny, err := codectypes.NewAnyWithValue(pk)
		require.NoError(t, err)
		// Get account address from pubkey
		validator := stakingtypes.Validator{
			OperatorAddress:   sdk.ValAddress(val.Address).String(),
			ConsensusPubkey:   pkAny,
			Jailed:            false,
			Status:            stakingtypes.Bonded,
			Tokens:            bondAmt,
			DelegatorShares:   sdk.OneDec(),
			Description:       stakingtypes.Description{},
			UnbondingHeight:   int64(0),
			UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingtypes.NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
			MinSelfDelegation: MinSelfDelegation,
		}

		validators = append(validators, validator)
		delegations = append(delegations, stakingtypes.NewDelegation(genAccs[0].GetAddress(), val.Address.Bytes(), sdk.OneDec()))

	}
	pk, _ := cryptocodec.FromTmPubKeyInterface(valSet.Validators[0].PubKey)

	accAddr := sdk.AccAddress(pk.Address())

	// set validators and delegations
	defaultStParams := stakingtypes.DefaultParams()
	stParams := stakingtypes.NewParams(
		defaultStParams.UnbondingTime,
		defaultStParams.MaxValidators,
		defaultStParams.MaxEntries,
		defaultStParams.HistoricalEntries,
		appparams.BondDenom,
	)

	// set validators and delegations
	stakingGenesis := stakingtypes.NewGenesisState(stParams, validators, delegations)
	genesisState[stakingtypes.ModuleName] = app.AppCodec().MustMarshalJSON(stakingGenesis)

	totalSupply := sdk.NewCoins()
	for _, b := range balances {
		// add genesis acc tokens to total supply
		totalSupply = totalSupply.Add(b.Coins...)
	}

	for range delegations {
		// add delegated tokens to total supply
		totalSupply = totalSupply.Add(sdk.NewCoin(appparams.BondDenom, bondAmt))
	}

	// add bonded amount to bonded pool module account
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(appparams.BondDenom, bondAmt)},
	})

	// update total supply
	bankGenesis := banktypes.NewGenesisState(
		banktypes.DefaultGenesisState().Params,
		balances,
		totalSupply,
		[]banktypes.Metadata{},
	)

	genesisState[banktypes.ModuleName] = app.AppCodec().MustMarshalJSON(bankGenesis)

	// Custom Cudos logic for gravity module
	// TODO: verify if true, maybe get the current snapshot of the gravity module genesis state

	gravityGenesis := &gravitytypes.GenesisState{
		Params: gravitytypes.DefaultParams(),
		StaticValCosmosAddrs: []string{
			accAddr.String(),
		},
		Erc20ToDenoms: []*gravitytypes.ERC20ToDenom{
			{
				Erc20: CudosMainnetEthAddr, // just for testing
				Denom: "acudos",
			},
		},
		DelegateKeys: []*gravitytypes.MsgSetOrchestratorAddress{
			{
				Validator:    sdk.ValAddress(valSet.Validators[0].Address).String(),
				EthAddress:   DefaultEthAddress,
				Orchestrator: genAccs[0].GetAddress().String(),
			},
		},
	}
	genesisState[gravitytypes.ModuleName] = app.AppCodec().MustMarshalJSON(gravityGenesis)

	// Add gentxs, MsgCreateValidator and MsgSetOrchestratorAddress

	return genesisState
}

func (s *KeeperTestHelper) Ed25519PubAddr() (cryptotypes.PrivKey, cryptotypes.PubKey, sdk.AccAddress) {
	key := ed25519.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

func (s *KeeperTestHelper) RandomAccountAddresses(n int) []sdk.AccAddress {
	addrsList := make([]sdk.AccAddress, n)
	for i := 0; i < n; i++ {
		_, _, addrs := testdata.KeyTestPubAddr()
		addrsList[i] = addrs
	}
	return addrsList
}

// From https://github.com/cosmos/cosmos-sdk/blob/v0.46.14/x/bank/testutil/test_helpers.go
func fundAccount(bankKeeper bankkeeper.Keeper, ctx sdk.Context, addr sdk.AccAddress, amounts sdk.Coins) error {
	if err := bankKeeper.MintCoins(ctx, cudoMinttypes.ModuleName, amounts); err != nil {
		return err
	}

	return bankKeeper.SendCoinsFromModuleToAccount(ctx, cudoMinttypes.ModuleName, addr, amounts)

}

// FundAcc funds target address with specified amount.
func (s *KeeperTestHelper) FundAcc(acc sdk.AccAddress, amounts sdk.Coins) {
	err := fundAccount(s.App.BankKeeper, s.Ctx, acc, amounts)
	s.Require().NoError(err)
}

// func (s *KeeperTestHelper) SetStaticValSet(cosmosAddress string) {
// 	s.App.GravityKeeper.SetStaticValCosmosAddr(s.Ctx, cosmosAddress)
// }
