package types

import (
	"github.com/CudoVentures/cudos-node/x/nft/exported"
	nfttypes "github.com/CudoVentures/cudos-node/x/nft/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

type NftKeeper interface {
	// Methods imported from nft should be defined here
	GetDenom(ctx sdk.Context, id string) (denom nfttypes.Denom, err error)
	GetNFT(ctx sdk.Context, denomID, tokenID string) (nft exported.NFT, err error)
	GetBaseNFT(ctx sdk.Context, denomID, tokenID string) (nft nfttypes.BaseNFT, err error)
	IsApprovedOperator(ctx sdk.Context, owner, operator sdk.AccAddress) bool
	TransferNftInternal(ctx sdk.Context, denomID string, tokenID string, from sdk.AccAddress, to sdk.AccAddress, nft nfttypes.BaseNFT)
	SoftLockNFT(ctx sdk.Context, lockOwner, denomID, tokenID string) error
	SoftUnlockNFT(ctx sdk.Context, lockOwner, denomID, tokenID string) error
	MintNFT(ctx sdk.Context, denomID string, tokenNm, tokenURI, tokenData string, sender, owner sdk.AccAddress) (string, error)
	GetOwners(ctx sdk.Context) (owners nfttypes.Owners, err error)
	IssueDenom(ctx sdk.Context, id, name, schema, symbol, traits, minter, description, data string, creator sdk.AccAddress) error
}

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, formModule string, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, fromAddr sdk.AccAddress, toModule string, amt sdk.Coins) error
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
}
