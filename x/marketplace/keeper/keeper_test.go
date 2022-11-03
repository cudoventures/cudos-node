package keeper_test

import (
	"fmt"
	"testing"

	testkeeper "github.com/CudoVentures/cudos-node/testutil/keeper"
	"github.com/CudoVentures/cudos-node/testutil/sample"
	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	nfttypes "github.com/CudoVentures/cudos-node/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func genAddresses(count int) []string {
	addresses := []string{}
	for i := 0; i < count; i++ {
		addresses = append(addresses, sample.AccAddress())
	}
	return addresses
}

func TestDistributeRoyalties(t *testing.T) {
	addresses := genAddresses(4)

	firstRoyaltyPrcent, err := sdk.NewDecFromStr("22.9")
	require.NoError(t, err)

	secondRoyaltyPercent, err := sdk.NewDecFromStr("0.01")
	require.NoError(t, err)

	thirdRoyaltyPercent, err := sdk.NewDecFromStr("30")
	require.NoError(t, err)

	kp, _, bankKeeper, ctx := testkeeper.MarketplaceKeeper(t)

	price := int64(10000)
	require.NoError(t, bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin("acudos", sdk.NewInt(price)))))

	err = kp.DistributeRoyalties(ctx, sdk.NewCoin("acudos", sdk.NewInt(price)), addresses[0], []types.Royalty{
		{
			Address: addresses[1],
			Percent: firstRoyaltyPrcent,
		},
		{
			Address: addresses[2],
			Percent: secondRoyaltyPercent,
		},
		{
			Address: addresses[3],
			Percent: thirdRoyaltyPercent,
		},
	})
	require.NoError(t, err)

	accountsBalances := bankKeeper.GetAccountsBalances(ctx)
	accBalancesMap := make(map[string]int64)

	distributedRoyalties := int64(0)

	for _, accBalance := range accountsBalances {
		accBalancesMap[accBalance.Address] = accBalance.Coins[0].Amount.Int64()

		distributedRoyalties += accBalance.Coins[0].Amount.Int64()
	}

	require.Equal(t, accBalancesMap[addresses[0]], int64(4709))
	require.Equal(t, accBalancesMap[addresses[1]], int64(2290))
	require.Equal(t, accBalancesMap[addresses[2]], int64(1))
	require.Equal(t, accBalancesMap[addresses[3]], int64(3000))

	require.Equal(t, distributedRoyalties, price)
}

func TestDistributeRoyaltiesShouldSkipIfEmptyRoyalties(t *testing.T) {
	kp, _, _, ctx := testkeeper.MarketplaceKeeper(t)
	require.NoError(t, kp.DistributeRoyalties(ctx, sdk.NewCoin("acudos", sdk.ZeroInt()), "", []types.Royalty{}))
}

func TestDistributeRoyaltiesShouldNotFailIfPriceIs1acudos(t *testing.T) {
	addresses := genAddresses(4)

	firstRoyaltyPrcent, err := sdk.NewDecFromStr("22.9")
	require.NoError(t, err)

	secondRoyaltyPercent, err := sdk.NewDecFromStr("0.01")
	require.NoError(t, err)

	thirdRoyaltyPercent, err := sdk.NewDecFromStr("30")
	require.NoError(t, err)

	kp, _, bankKeeper, ctx := testkeeper.MarketplaceKeeper(t)

	price := int64(1)
	require.NoError(t, bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin("acudos", sdk.NewInt(price)))))

	err = kp.DistributeRoyalties(ctx, sdk.NewCoin("acudos", sdk.NewInt(price)), addresses[0], []types.Royalty{
		{
			Address: addresses[1],
			Percent: firstRoyaltyPrcent,
		},
		{
			Address: addresses[2],
			Percent: secondRoyaltyPercent,
		},
		{
			Address: addresses[3],
			Percent: thirdRoyaltyPercent,
		},
	})
	require.NoError(t, err)

	accountsBalances := bankKeeper.GetAccountsBalances(ctx)
	accBalancesMap := make(map[string]int64)

	distributedRoyalties := int64(0)

	for _, accBalance := range accountsBalances {
		accBalancesMap[accBalance.Address] = accBalance.Coins[0].Amount.Int64()

		distributedRoyalties += accBalance.Coins[0].Amount.Int64()
	}

	require.Equal(t, int64(1), accBalancesMap[addresses[0]])
	require.Equal(t, int64(0), accBalancesMap[addresses[1]])
	require.Equal(t, int64(0), accBalancesMap[addresses[2]])
	require.Equal(t, int64(0), accBalancesMap[addresses[3]])
}

func TestPublishCollectionShouldBeSuccessfulWithoutRoyalties(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	kp, nftKeeper, _, ctx := testkeeper.MarketplaceKeeper(t)

	nftKeeper.SetDenom(ctx, nfttypes.NewDenom("testdenom", "testname", "{}", "testsym", "", "", "", "", owner))
	id, err := kp.PublishCollection(ctx, types.NewCollection("testdenom", []types.Royalty{}, []types.Royalty{}, owner.String(), false))
	require.NoError(t, err)
	require.Equal(t, uint64(0), id)

	_, found := kp.GetCollection(ctx, id)
	require.True(t, found)
}

func TestPublishCollectionShouldBeSuccessfulWithRoyalties(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	kp, nftKeeper, _, ctx := testkeeper.MarketplaceKeeper(t)

	err = nftKeeper.SetDenom(ctx, nfttypes.NewDenom("testdenom", "testname", "{}", "testsym", "", "", "", "", owner))
	require.NoError(t, err)

	addresses := genAddresses(2)

	id, err := kp.PublishCollection(ctx, types.NewCollection("testdenom", []types.Royalty{
		{
			Address: addresses[0],
			Percent: sdk.MustNewDecFromStr("0.01"),
		},
	}, []types.Royalty{
		{
			Address: addresses[1],
			Percent: sdk.MustNewDecFromStr("51.98"),
		},
	}, owner.String(), false))
	require.NoError(t, err)
	require.Equal(t, uint64(0), id)

	_, found := kp.GetCollection(ctx, id)
	require.True(t, found)
}

func TestPublishCollectionShouldFailIfNotDenomOwner(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	publisher, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	kp, nftKeeper, _, ctx := testkeeper.MarketplaceKeeper(t)

	err = nftKeeper.SetDenom(ctx, nfttypes.NewDenom("testdenom", "testname", "{}", "testsym", "", "", "", "", owner))
	require.NoError(t, err)

	_, err = kp.PublishCollection(ctx, types.NewCollection("testdenom", []types.Royalty{}, []types.Royalty{}, publisher.String(), false))
	require.Equal(t, fmt.Sprintf("Owner of denom testdenom is %s: not denom owner", owner.String()), err.Error())
}

func TestPublishCollectionShouldFailIfAlreadyPublished(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	kp, nftKeeper, _, ctx := testkeeper.MarketplaceKeeper(t)

	err = nftKeeper.SetDenom(ctx, nfttypes.NewDenom("testdenom", "testname", "{}", "testsym", "", "", "", "", owner))
	require.NoError(t, err)

	_, err = kp.PublishCollection(ctx, types.NewCollection("testdenom", []types.Royalty{}, []types.Royalty{}, owner.String(), false))
	require.NoError(t, err)

	_, err = kp.PublishCollection(ctx, types.NewCollection("testdenom", []types.Royalty{}, []types.Royalty{}, owner.String(), false))
	require.Equal(t, "Collection for denom testdenom is already published: collection already published", err.Error())
}

func TestPublishCollectionShouldFailIfDenomDoesntExist(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	kp, _, _, ctx := testkeeper.MarketplaceKeeper(t)

	_, err = kp.PublishCollection(ctx, types.NewCollection("testdenom", []types.Royalty{}, []types.Royalty{}, owner.String(), false))
	require.Equal(t, "not found denomID: testdenom: invalid denom", err.Error())
}

func TestPublishNftShouldFailIfDenomDoesntExist(t *testing.T) {
	kp, _, _, ctx := testkeeper.MarketplaceKeeper(t)

	_, err := kp.PublishNFT(ctx, types.NewNft("1", "denom", "owner", sdk.NewCoin("acudos", sdk.NewInt(10000))))
	require.Equal(t, "not found denomID: denom: invalid denom", err.Error())
}

func TestPublishNftShouldFailIfNftDoesntExist(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	kp, nftKeeper, _, ctx := testkeeper.MarketplaceKeeper(t)

	err = nftKeeper.SetDenom(ctx, nfttypes.NewDenom("testdenom", "testname", "{}", "testsym", "", "", "", "", owner))
	require.NoError(t, err)

	_, err = kp.PublishNFT(ctx, types.NewNft("1", "testdenom", "owner", sdk.NewCoin("acudos", sdk.NewInt(10000))))
	require.Equal(t, "not found NFT: testdenom: unknown nft collection", err.Error())
}

func TestAfterPublishNftShouldNotBeAbleToTransferTheNft(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	kp, nftKeeper, _, ctx := testkeeper.MarketplaceKeeper(t)

	err = nftKeeper.SetDenom(ctx, nfttypes.NewDenom("testdenom", "testname", "{}", "testsym", "", "", "", "", owner))
	require.NoError(t, err)

	id, err := nftKeeper.MintNFT(ctx, "testdenom", "first", "", "", owner, owner)
	require.NoError(t, err)
	require.Equal(t, "1", id)

	publishedID, err := kp.PublishNFT(ctx, types.NewNft("1", "testdenom", owner.String(), sdk.NewCoin("acudos", sdk.NewInt(10000))))
	require.NoError(t, err)
	require.Equal(t, uint64(0), publishedID)

	newOwner, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	err = nftKeeper.TransferOwner(ctx, "testdenom", "1", owner, newOwner, owner)
	require.Equal(t, "token id 1 from denom with id testdenom is soft locked by marketplace: soft locked", err.Error())
}

func TestPublishNftShouldFailIfAlreadyPublished(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	kp, nftKeeper, _, ctx := testkeeper.MarketplaceKeeper(t)

	err = nftKeeper.SetDenom(ctx, nfttypes.NewDenom("testdenom", "testname", "{}", "testsym", "", "", "", "", owner))
	require.NoError(t, err)

	id, err := nftKeeper.MintNFT(ctx, "testdenom", "first", "", "", owner, owner)
	require.NoError(t, err)
	require.Equal(t, "1", id)

	publishedID, err := kp.PublishNFT(ctx, types.NewNft("1", "testdenom", owner.String(), sdk.NewCoin("acudos", sdk.NewInt(10000))))
	require.NoError(t, err)
	require.Equal(t, uint64(0), publishedID)

	_, err = kp.PublishNFT(ctx, types.NewNft("1", "testdenom", owner.String(), sdk.NewCoin("acudos", sdk.NewInt(10000))))
	require.Equal(t, "nft with token id (1) from denom (testdenom) already published for sale: nft already published", err.Error())
}

func TestPublishNftShouldFailIfNotNftOwner(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	kp, nftKeeper, _, ctx := testkeeper.MarketplaceKeeper(t)

	err = nftKeeper.SetDenom(ctx, nfttypes.NewDenom("testdenom", "testname", "{}", "testsym", "", "", "", "", owner))
	require.NoError(t, err)

	id, err := nftKeeper.MintNFT(ctx, "testdenom", "first", "", "", owner, owner)
	require.NoError(t, err)
	require.Equal(t, "1", id)

	notOwner, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	_, err = kp.PublishNFT(ctx, types.NewNft("1", "testdenom", notOwner.String(), sdk.NewCoin("acudos", sdk.NewInt(10000))))
	require.Equal(t, fmt.Sprintf("%s not nft owner or approved operator for token id (1) from denom (testdenom): not nft owner", notOwner.String()), err.Error())
}

func TestPublishNftShouldBeSuccessfulByNftApprovedAddresses(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	kp, nftKeeper, _, ctx := testkeeper.MarketplaceKeeper(t)

	err = nftKeeper.SetDenom(ctx, nfttypes.NewDenom("testdenom", "testname", "{}", "testsym", "", "", "", "", owner))
	require.NoError(t, err)

	id, err := nftKeeper.MintNFT(ctx, "testdenom", "first", "", "", owner, owner)
	require.NoError(t, err)
	require.Equal(t, "1", id)

	approvedOperator, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	require.NoError(t, nftKeeper.AddApproval(ctx, "testdenom", "1", owner, approvedOperator))

	_, err = kp.PublishNFT(ctx, types.NewNft("1", "testdenom", approvedOperator.String(), sdk.NewCoin("acudos", sdk.NewInt(10000))))
	require.NoError(t, err)
}

func TestPublishNftShouldBeSuccessfulByNftApprovedAddress(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	kp, nftKeeper, _, ctx := testkeeper.MarketplaceKeeper(t)

	err = nftKeeper.SetDenom(ctx, nfttypes.NewDenom("testdenom", "testname", "{}", "testsym", "", "", "", "", owner))
	require.NoError(t, err)

	id, err := nftKeeper.MintNFT(ctx, "testdenom", "first", "", "", owner, owner)
	require.NoError(t, err)
	require.Equal(t, "1", id)

	approvedOperator, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	require.NoError(t, nftKeeper.AddApproval(ctx, "testdenom", "1", owner, approvedOperator))

	_, err = kp.PublishNFT(ctx, types.NewNft("1", "testdenom", approvedOperator.String(), sdk.NewCoin("acudos", sdk.NewInt(10000))))
	require.NoError(t, err)
}

func TestPublishNftShouldBeSuccessfulByApprovedOperator(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	kp, nftKeeper, _, ctx := testkeeper.MarketplaceKeeper(t)

	err = nftKeeper.SetDenom(ctx, nfttypes.NewDenom("testdenom", "testname", "{}", "testsym", "", "", "", "", owner))
	require.NoError(t, err)

	id, err := nftKeeper.MintNFT(ctx, "testdenom", "first", "", "", owner, owner)
	require.NoError(t, err)
	require.Equal(t, "1", id)

	approvedOperator, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	require.NoError(t, nftKeeper.AddApprovalForAll(ctx, owner, approvedOperator, true))

	_, err = kp.PublishNFT(ctx, types.NewNft("1", "testdenom", approvedOperator.String(), sdk.NewCoin("acudos", sdk.NewInt(10000))))
	require.NoError(t, err)
}

func TestPublishNftShouldBeSuccessfulByOwner(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	kp, nftKeeper, _, ctx := testkeeper.MarketplaceKeeper(t)

	err = nftKeeper.SetDenom(ctx, nfttypes.NewDenom("testdenom", "testname", "{}", "testsym", "", "", "", "", owner))
	require.NoError(t, err)

	id, err := nftKeeper.MintNFT(ctx, "testdenom", "first", "", "", owner, owner)
	require.NoError(t, err)
	require.Equal(t, "1", id)

	_, err = kp.PublishNFT(ctx, types.NewNft("1", "testdenom", owner.String(), sdk.NewCoin("acudos", sdk.NewInt(10000))))
	require.NoError(t, err)
}

func TestBuyNftShouldFailForNotExistingId(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	kp, _, _, ctx := testkeeper.MarketplaceKeeper(t)

	err = kp.BuyNFT(ctx, 0, owner)
	require.Equal(t, "nft with id (0) is not found for sale: nft not found", err.Error())
}

func TestBuyNftShouldFailWhenBuyingOwnNft(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	kp, nftKeeper, _, ctx := testkeeper.MarketplaceKeeper(t)

	err = nftKeeper.SetDenom(ctx, nfttypes.NewDenom("testdenom", "testname", "{}", "testsym", "", "", "", "", owner))
	require.NoError(t, err)

	id, err := nftKeeper.MintNFT(ctx, "testdenom", "first", "", "", owner, owner)
	require.NoError(t, err)
	require.Equal(t, "1", id)

	publishedId, err := kp.PublishNFT(ctx, types.NewNft("1", "testdenom", owner.String(), sdk.NewCoin("acudos", sdk.NewInt(10000))))
	require.NoError(t, err)

	err = kp.BuyNFT(ctx, publishedId, owner)
	require.Equal(t, "cannot buy own nft: cannot buy own nft", err.Error())
}

func TestBuyNftShouldBeSuccessfulWithResaleRoyalties(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	kp, nftKeeper, bankKeeper, ctx := testkeeper.MarketplaceKeeper(t)

	err = nftKeeper.SetDenom(ctx, nfttypes.NewDenom("testdenom", "testname", "{}", "testsym", "", "", "", "", owner))
	require.NoError(t, err)

	id, err := nftKeeper.MintNFT(ctx, "testdenom", "first", "", "", owner, owner)
	require.NoError(t, err)
	require.Equal(t, "1", id)

	royaltyReceiver := genAddresses(1)[0]

	kp.PublishCollection(ctx, types.NewCollection("testdenom", []types.Royalty{}, []types.Royalty{
		{
			Address: royaltyReceiver,
			Percent: sdk.MustNewDecFromStr("50"),
		},
	}, owner.String(), false))

	price := int64(10000)
	publishedId, err := kp.PublishNFT(ctx, types.NewNft("1", "testdenom", owner.String(), sdk.NewCoin("acudos", sdk.NewInt(price))))
	require.NoError(t, err)

	buyer, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	require.NoError(t, bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin("acudos", sdk.NewInt(price)))))
	err = bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, buyer, sdk.NewCoins(sdk.NewCoin("acudos", sdk.NewInt(price))))
	require.NoError(t, err)

	require.NoError(t, kp.BuyNFT(ctx, publishedId, buyer))

	royaltyReceiverAddr, err := sdk.AccAddressFromBech32(royaltyReceiver)
	require.NoError(t, err)

	require.Equal(t, true, bankKeeper.GetBalance(ctx, royaltyReceiverAddr, "acudos").Amount.Equal(sdk.NewInt(5000)))

	baseNft, err := nftKeeper.GetBaseNFT(ctx, "testdenom", "1")
	require.NoError(t, err)

	require.True(t, nftKeeper.IsOwner(baseNft, buyer))
}

func TestBuyNftShouldBeSuccessfulWithoutResaleRoyalties(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	kp, nftKeeper, bankKeeper, ctx := testkeeper.MarketplaceKeeper(t)

	err = nftKeeper.SetDenom(ctx, nfttypes.NewDenom("testdenom", "testname", "{}", "testsym", "", "", "", "", owner))
	require.NoError(t, err)

	id, err := nftKeeper.MintNFT(ctx, "testdenom", "first", "", "", owner, owner)
	require.NoError(t, err)
	require.Equal(t, "1", id)

	kp.PublishCollection(ctx, types.NewCollection("testdenom", []types.Royalty{}, []types.Royalty{}, owner.String(), false))

	price := int64(10000)
	publishedId, err := kp.PublishNFT(ctx, types.NewNft("1", "testdenom", owner.String(), sdk.NewCoin("acudos", sdk.NewInt(price))))
	require.NoError(t, err)

	buyer, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	require.NoError(t, bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin("acudos", sdk.NewInt(price)))))
	err = bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, buyer, sdk.NewCoins(sdk.NewCoin("acudos", sdk.NewInt(price))))
	require.NoError(t, err)

	require.NoError(t, kp.BuyNFT(ctx, publishedId, buyer))

	require.Equal(t, true, bankKeeper.GetBalance(ctx, owner, "acudos").Amount.Equal(sdk.NewInt(10000)))

	baseNft, err := nftKeeper.GetBaseNFT(ctx, "testdenom", "1")
	require.NoError(t, err)

	require.True(t, nftKeeper.IsOwner(baseNft, buyer))
}

func TestBuyNftShouldFailWhenInsufficientFunds(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	kp, nftKeeper, _, ctx := testkeeper.MarketplaceKeeper(t)

	err = nftKeeper.SetDenom(ctx, nfttypes.NewDenom("testdenom", "testname", "{}", "testsym", "", "", "", "", owner))
	require.NoError(t, err)

	id, err := nftKeeper.MintNFT(ctx, "testdenom", "first", "", "", owner, owner)
	require.NoError(t, err)
	require.Equal(t, "1", id)

	kp.PublishCollection(ctx, types.NewCollection("testdenom", []types.Royalty{}, []types.Royalty{}, owner.String(), false))

	publishedId, err := kp.PublishNFT(ctx, types.NewNft("1", "testdenom", owner.String(), sdk.NewCoin("acudos", sdk.NewInt(10000))))
	require.NoError(t, err)

	buyer, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	err = kp.BuyNFT(ctx, publishedId, buyer)
	require.Equal(t, "0acudos is smaller than 10000acudos: insufficient funds", err.Error())
}

func TestMintNftShouldFailIfDenomNotFound(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	kp, _, _, ctx := testkeeper.MarketplaceKeeper(t)
	_, err = kp.MintNFT(ctx, "testdenom", "testname", "testuri", "", sdk.NewCoin("acudos", sdk.NewInt(10000)), owner, owner)
	require.Equal(t, "not found denomID: testdenom: invalid denom", err.Error())
}

func TestMintNftShouldFailIfCollectionNotPublished(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	kp, nftKeeper, _, ctx := testkeeper.MarketplaceKeeper(t)
	err = nftKeeper.IssueDenom(ctx, "testdenom", "testname", "{}", "testsymbol", "", "", "", "", owner)
	require.NoError(t, err)

	_, err = kp.MintNFT(ctx, "testdenom", "testname", "testuri", "", sdk.NewCoin("acudos", sdk.NewInt(10000)), owner, owner)
	require.Equal(t, "collection testdenom not published for sale: collection not published for sale", err.Error())
}

func TestMintNftShouldFailIfCollectionNotVerified(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	kp, nftKeeper, _, ctx := testkeeper.MarketplaceKeeper(t)
	err = nftKeeper.IssueDenom(ctx, "testdenom", "testname", "{}", "testsymbol", "", "", "", "", owner)
	require.NoError(t, err)

	_, err = kp.PublishCollection(ctx, types.NewCollection("testdenom", []types.Royalty{}, []types.Royalty{}, owner.String(), false))
	require.NoError(t, err)

	_, err = kp.MintNFT(ctx, "testdenom", "testname", "testuri", "", sdk.NewCoin("acudos", sdk.NewInt(10000)), owner, owner)
	require.Equal(t, "collection 0 is not verified: collection is unverified", err.Error())
}

func TestMintNftShouldBeSuccessfulWithRoyalties(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	minter, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	kp, nftKeeper, bankKeeper, ctx := testkeeper.MarketplaceKeeper(t)
	err = nftKeeper.IssueDenom(ctx, "testdenom", "testname", "{}", "testsymbol", "", minter.String(), "", "", owner)
	require.NoError(t, err)

	_, err = kp.PublishCollection(ctx, types.NewCollection("testdenom", []types.Royalty{
		{
			Address: owner.String(),
			Percent: sdk.MustNewDecFromStr("100"),
		},
	}, []types.Royalty{}, owner.String(), true))
	require.NoError(t, err)

	price := int64(10000)
	require.NoError(t, bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin("acudos", sdk.NewInt(price)))))
	err = bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, minter, sdk.NewCoins(sdk.NewCoin("acudos", sdk.NewInt(price))))
	require.NoError(t, err)

	_, err = kp.MintNFT(ctx, "testdenom", "testname", "testuri", "", sdk.NewCoin("acudos", sdk.NewInt(10000)), minter, minter)
	require.NoError(t, err)

	require.Equal(t, true, bankKeeper.GetBalance(ctx, owner, "acudos").Amount.Equal(sdk.NewInt(10000)))

	baseNft, err := nftKeeper.GetBaseNFT(ctx, "testdenom", "1")
	require.NoError(t, err)

	require.True(t, nftKeeper.IsOwner(baseNft, minter))
}

func TestMintNftShouldBeSuccessfulWithoutRoyalties(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	minter, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	kp, nftKeeper, bankKeeper, ctx := testkeeper.MarketplaceKeeper(t)
	err = nftKeeper.IssueDenom(ctx, "testdenom", "testname", "{}", "testsymbol", "", minter.String(), "", "", owner)
	require.NoError(t, err)

	_, err = kp.PublishCollection(ctx, types.NewCollection("testdenom", []types.Royalty{}, []types.Royalty{}, owner.String(), true))
	require.NoError(t, err)

	price := int64(10000)
	require.NoError(t, bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin("acudos", sdk.NewInt(price)))))
	err = bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, minter, sdk.NewCoins(sdk.NewCoin("acudos", sdk.NewInt(price))))
	require.NoError(t, err)

	_, err = kp.MintNFT(ctx, "testdenom", "testname", "testuri", "", sdk.NewCoin("acudos", sdk.NewInt(10000)), minter, minter)
	require.NoError(t, err)

	baseNft, err := nftKeeper.GetBaseNFT(ctx, "testdenom", "1")
	require.NoError(t, err)

	require.True(t, nftKeeper.IsOwner(baseNft, minter))
}

func TestMintNftShouldFailWhenInsufficientFunds(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	kp, nftKeeper, _, ctx := testkeeper.MarketplaceKeeper(t)
	err = nftKeeper.IssueDenom(ctx, "testdenom", "testname", "{}", "testsymbol", "", "", "", "", owner)
	require.NoError(t, err)

	_, err = kp.PublishCollection(ctx, types.NewCollection("testdenom", []types.Royalty{}, []types.Royalty{}, owner.String(), true))
	require.NoError(t, err)

	_, err = kp.MintNFT(ctx, "testdenom", "testname", "testuri", "", sdk.NewCoin("acudos", sdk.NewInt(10000)), owner, owner)
	require.Equal(t, "0acudos is smaller than 10000acudos: insufficient funds", err.Error())
}

func TestCreateCollection(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	kp, _, _, ctx := testkeeper.MarketplaceKeeper(t)
	id, err := kp.CreateCollection(ctx, owner, "testdenom", "testname", "{}", "symbol", "", "", "", "", []types.Royalty{}, []types.Royalty{}, false)
	require.NoError(t, err)
	require.Equal(t, uint64(0), id)
}

func TestSetCollectionStatus(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	kp, _, _, ctx := testkeeper.MarketplaceKeeper(t)
	id, err := kp.CreateCollection(ctx, owner, "testdenom", "testname", "{}", "symbol", "", "", "", "", []types.Royalty{}, []types.Royalty{}, false)
	require.NoError(t, err)
	require.Equal(t, uint64(0), id)

	coll, found := kp.GetCollection(ctx, id)
	require.True(t, found)

	require.False(t, coll.Verified)
	require.NoError(t, kp.SetCollectionStatus(ctx, id, true))

	coll, found = kp.GetCollection(ctx, id)
	require.True(t, found)

	require.True(t, coll.Verified)
}

func TestSetCollectionRoyalties(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	kp, _, _, ctx := testkeeper.MarketplaceKeeper(t)
	resaleRoyalties := []types.Royalty{
		{
			Address: "resale_royalty_receiver",
			Percent: sdk.MustNewDecFromStr("0.01"),
		},
	}
	id, err := kp.CreateCollection(ctx, owner, "testdenom", "testname", "{}", "symbol", "", "", "", "", []types.Royalty{}, resaleRoyalties, false)
	require.NoError(t, err)
	require.Equal(t, uint64(0), id)

	coll, found := kp.GetCollection(ctx, id)
	require.True(t, found)

	require.Len(t, coll.MintRoyalties, 0)
	require.Equal(t, resaleRoyalties, coll.ResaleRoyalties)

	mintRoyalties := []types.Royalty{
		{
			Address: "mint_royalty_receiver",
			Percent: sdk.MustNewDecFromStr("0.01"),
		},
	}

	err = kp.SetCollectionRoyalties(ctx, owner.String(), id, mintRoyalties, []types.Royalty{})
	require.NoError(t, err)

	coll, found = kp.GetCollection(ctx, id)
	require.True(t, found)

	require.Equal(t, mintRoyalties, coll.MintRoyalties)
	require.Len(t, coll.ResaleRoyalties, 0)
}

func TestSetCollectionRoyaltiesShouldFailIfNotOwner(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	kp, _, _, ctx := testkeeper.MarketplaceKeeper(t)
	id, err := kp.CreateCollection(ctx, owner, "testdenom", "testname", "{}", "symbol", "", "", "", "", []types.Royalty{}, []types.Royalty{}, false)
	require.NoError(t, err)
	require.Equal(t, uint64(0), id)

	setter, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	err = kp.SetCollectionRoyalties(ctx, setter.String(), id, []types.Royalty{}, []types.Royalty{})
	require.Equal(t, fmt.Sprintf("owner of collection 0 is %s, not %s: not collection owner", owner.String(), setter.String()), err.Error())
}

func TestSetNftPriceShouldFailIfNftNotFound(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	kp, _, _, ctx := testkeeper.MarketplaceKeeper(t)
	err = kp.SetNftPrice(ctx, owner.String(), 0, sdk.NewCoin("acudos", sdk.NewInt(1)))
	require.Equal(t, "NFT with id 0 not found: nft not found", err.Error())
}

func TestSetNftPriceShouldFailIfNotOwner(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	kp, nftKeeper, _, ctx := testkeeper.MarketplaceKeeper(t)

	id, err := kp.CreateCollection(ctx, owner, "testdenom", "testname", "{}", "symbol", "", "", "", "", []types.Royalty{}, []types.Royalty{}, false)
	require.NoError(t, err)
	require.Equal(t, uint64(0), id)

	mintedNftId, err := nftKeeper.MintNFT(ctx, "testdenom", "first", "", "", owner, owner)
	require.NoError(t, err)
	require.Equal(t, "1", mintedNftId)

	nftId, err := kp.PublishNFT(ctx, types.NewNft("1", "testdenom", owner.String(), sdk.NewCoin("acudos", sdk.NewInt(10000))))
	require.NoError(t, err)

	setter, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	err = kp.SetNftPrice(ctx, setter.String(), nftId, sdk.NewCoin("acudos", sdk.NewInt(5)))
	require.Equal(t, fmt.Sprintf("owner of NFT 0 is %s, not %s: not collection owner", owner.String(), setter.String()), err.Error())
}

func TestSetNftPrice(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32(genAddresses(1)[0])
	require.NoError(t, err)

	kp, nftKeeper, _, ctx := testkeeper.MarketplaceKeeper(t)

	id, err := kp.CreateCollection(ctx, owner, "testdenom", "testname", "{}", "symbol", "", "", "", "", []types.Royalty{}, []types.Royalty{}, false)
	require.NoError(t, err)
	require.Equal(t, uint64(0), id)

	mintedNftId, err := nftKeeper.MintNFT(ctx, "testdenom", "first", "", "", owner, owner)
	require.NoError(t, err)
	require.Equal(t, "1", mintedNftId)

	nftId, err := kp.PublishNFT(ctx, types.NewNft("1", "testdenom", owner.String(), sdk.NewCoin("acudos", sdk.NewInt(10000))))
	require.NoError(t, err)

	err = kp.SetNftPrice(ctx, owner.String(), nftId, sdk.NewCoin("acudos", sdk.NewInt(5)))
	require.NoError(t, err)

	nft, found := kp.GetNft(ctx, nftId)
	require.True(t, found)
	require.Equal(t, int64(5), nft.Price.Amount.Int64())
}
