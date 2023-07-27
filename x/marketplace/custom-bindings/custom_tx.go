package custom_bindings

import (
	"encoding/json"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	"github.com/CudoVentures/cudos-node/x/marketplace/keeper"
	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ wasmkeeper.Messenger = (*MarketplaceWasmMessengerDecorator)(nil)

type MarketplaceWasmMessengerDecorator struct {
	old    wasmkeeper.Messenger
	msgSrv types.MsgServer
}

func NewMarketplaceWasmMessageDecorator(marketplaceKeeper keeper.Keeper) func(wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &MarketplaceWasmMessengerDecorator{
			old:    old,
			msgSrv: keeper.NewMsgServerImpl(marketplaceKeeper),
		}
	}
}

func (m *MarketplaceWasmMessengerDecorator) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) (events []sdk.Event, data [][]byte, err error) {
	if msg.Custom == nil {
		return m.old.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
	}

	var customMsg marketplaceCustomMsg
	if err := json.Unmarshal(msg.Custom, &customMsg); err != nil {
		return nil, nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	switch {
	case customMsg.PublishCollection != nil:
		c := customMsg.PublishCollection
		msg := types.NewMsgPublishCollection(c.Creator, c.DenomId, c.MintRoyalties, c.ResaleRoyalties)
		if err = msg.ValidateBasic(); err != nil {
			return nil, nil, err
		}
		_, err := m.msgSrv.PublishCollection(ctx, msg)
		return nil, nil, err
	case customMsg.PublishNft != nil:
		c := customMsg.PublishNft
		msg := types.NewMsgPublishNft(c.Creator, c.TokenId, c.DenomId, c.Price)
		if err = msg.ValidateBasic(); err != nil {
			return nil, nil, err
		}
		_, err := m.msgSrv.PublishNft(ctx, msg)
		return nil, nil, err
	case customMsg.BuyNft != nil:
		c := customMsg.BuyNft
		msg := types.NewMsgBuyNft(c.Creator, c.Id)
		if err = msg.ValidateBasic(); err != nil {
			return nil, nil, err
		}
		_, err := m.msgSrv.BuyNft(ctx, msg)
		return nil, nil, err
	case customMsg.MintNft != nil:
		c := customMsg.MintNft
		msg := types.NewMsgMintNft(c.Creator, c.DenomId, c.Recipient, c.Name, c.URI, c.Data, c.Uid, c.Price)
		if err = msg.ValidateBasic(); err != nil {
			return nil, nil, err
		}
		_, err := m.msgSrv.MintNft(ctx, msg)
		return nil, nil, err
	case customMsg.RemoveNft != nil:
		c := customMsg.RemoveNft
		msg := types.NewMsgRemoveNft(c.Creator, c.Id)
		if err = msg.ValidateBasic(); err != nil {
			return nil, nil, err
		}
		_, err := m.msgSrv.RemoveNft(ctx, msg)
		return nil, nil, err
	case customMsg.VerifyCollection != nil:
		c := customMsg.VerifyCollection
		msg := types.NewMsgVerifyCollection(c.Creator, c.Id)
		if err = msg.ValidateBasic(); err != nil {
			return nil, nil, err
		}
		_, err := m.msgSrv.VerifyCollection(ctx, msg)
		return nil, nil, err
	case customMsg.UnverifyCollection != nil:
		c := customMsg.UnverifyCollection
		msg := types.NewMsgUnverifyCollection(c.Creator, c.Id)
		if err = msg.ValidateBasic(); err != nil {
			return nil, nil, err
		}
		_, err := m.msgSrv.UnverifyCollection(ctx, msg)
		return nil, nil, err
	case customMsg.CreateCollection != nil:
		c := customMsg.CreateCollection
		msg := types.NewMsgCreateCollection(c.Creator, c.Id, c.Name, c.Schema, c.Symbol, c.Traits, c.Description, c.Minter, c.Data, c.MintRoyalties, c.ResaleRoyalties, c.Verified)
		if err = msg.ValidateBasic(); err != nil {
			return nil, nil, err
		}
		_, err := m.msgSrv.CreateCollection(ctx, msg)
		return nil, nil, err
	case customMsg.UpdateRoyalties != nil:
		c := customMsg.UpdateRoyalties
		msg := types.NewMsgUpdateRoyalties(c.Creator, c.Id, c.MintRoyalties, c.ResaleRoyalties)
		if err = msg.ValidateBasic(); err != nil {
			return nil, nil, err
		}
		_, err := m.msgSrv.UpdateRoyalties(ctx, msg)
		return nil, nil, err
	case customMsg.UpdatePrice != nil:
		c := customMsg.UpdatePrice
		msg := types.NewMsgUpdatePrice(c.Creator, c.Id, c.Price)
		if err = msg.ValidateBasic(); err != nil {
			return nil, nil, err
		}
		_, err := m.msgSrv.UpdatePrice(ctx, msg)
		return nil, nil, err
	case customMsg.AddAdmin != nil:
		c := customMsg.AddAdmin
		msg := types.NewMsgAddAdmin(c.Creator, c.Address)
		if err = msg.ValidateBasic(); err != nil {
			return nil, nil, err
		}
		_, err := m.msgSrv.AddAdmin(ctx, msg)
		return nil, nil, err
	case customMsg.RemoveAdmin != nil:
		c := customMsg.RemoveAdmin
		msg := types.NewMsgRemoveAdmin(c.Creator, c.Address)
		if err = msg.ValidateBasic(); err != nil {
			return nil, nil, err
		}
		_, err := m.msgSrv.RemoveAdmin(ctx, msg)
		return nil, nil, err
	default:
		return m.old.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
	}
}

type marketplaceCustomMsg struct {
	PublishCollection  *PublishCollectionRequest  `json:"publish_collection_msg,omitempty"`
	PublishNft         *PublishNftRequest         `json:"publish_nft_msg,omitempty"`
	BuyNft             *BuyNftRequest             `json:"buy_nft_msg,omitempty"`
	MintNft            *MintNftRequest            `json:"mint_nft_marketplace_msg,omitempty"`
	RemoveNft          *RemoveNftRequest          `json:"remove_nft_msg,omitempty"`
	VerifyCollection   *VerifyCollectionRequest   `json:"verify_collection_msg,omitempty"`
	UnverifyCollection *UnverifyCollectionRequest `json:"unverify_collection_msg,omitempty"`
	CreateCollection   *CreateCollectionRequest   `json:"create_collection_msg,omitempty"`
	UpdateRoyalties    *UpdateRoyaltiesRequest    `json:"update_royalties_msg,omitempty"`
	UpdatePrice        *UpdatePriceRequest        `json:"update_price_msg,omitempty"`
	AddAdmin           *AddAdminRequest           `json:"add_admin_msg,omitempty"`
	RemoveAdmin        *RemoveAdminRequest        `json:"remove_admin_msg,omitempty"`
}

type PublishCollectionRequest struct {
	Creator         string          `json:"creator"`
	DenomId         string          `json:"denom_id"`
	MintRoyalties   []types.Royalty `json:"mint_royalties"`
	ResaleRoyalties []types.Royalty `json:"resale_royalties"`
}

type PublishNftRequest struct {
	Creator string   `json:"creator"`
	DenomId string   `json:"denom_id"`
	TokenId string   `json:"token_id"`
	Price   sdk.Coin `json:"price"`
}

type BuyNftRequest struct {
	Creator string `json:"creator"`
	Id      uint64 `json:"id"`
}

type MintNftRequest struct {
	Creator   string   `json:"creator"`
	DenomId   string   `json:"denom_id"`
	Recipient string   `json:"recipient"`
	Price     sdk.Coin `json:"price"`
	Name      string   `json:"name"`
	URI       string   `json:"uri,omitempty"`
	Data      string   `json:"data,omitempty"`
	Uid       string   `json:"uid"`
}

type RemoveNftRequest struct {
	Creator string `json:"creator"`
	Id      uint64 `json:"id"`
}

type VerifyCollectionRequest struct {
	Creator string `json:"creator"`
	Id      uint64 `json:"id"`
}

type UnverifyCollectionRequest struct {
	Creator string `json:"creator"`
	Id      uint64 `json:"id"`
}

type CreateCollectionRequest struct {
	Creator         string          `json:"creator"`
	Id              string          `json:"id"`
	Name            string          `json:"name"`
	Schema          string          `json:"schema"`
	Symbol          string          `json:"symbol"`
	Traits          string          `json:"traits"`
	Description     string          `json:"description"`
	Minter          string          `json:"minter"`
	Data            string          `json:"data"`
	MintRoyalties   []types.Royalty `json:"mint_royalties"`
	ResaleRoyalties []types.Royalty `json:"resale_royalties"`
	Verified        bool            `json:"verified"`
}

type UpdateRoyaltiesRequest struct {
	Creator         string          `json:"creator"`
	Id              uint64          `json:"id"`
	MintRoyalties   []types.Royalty `json:"mint_royalties"`
	ResaleRoyalties []types.Royalty `json:"resale_royalties"`
}

type UpdatePriceRequest struct {
	Creator string   `json:"creator"`
	Id      uint64   `json:"id"`
	Price   sdk.Coin `json:"price"`
}

type AddAdminRequest struct {
	Creator string `json:"creator"`
	Address string `json:"address"`
}

type RemoveAdminRequest struct {
	Creator string `json:"creator"`
	Address string `json:"address"`
}
