# NFT Marketplace Module Specification

## Overview

A module for trading NFTs on the CUDOS network. The module supports publishing collection for sale which is optional step to set mint royalties and resale royalties 
which will be deducted from the payments for all the NFTs traded from the given collection.
If user tries to sell NFT that has no royalties the seller will be paid the full price.
The NFTs can be minted via the marketplace which allows distribution of royalties on mint.
When NFT is published for sale, it will be soft locked in the NFT module, so the seller will remain owner but wont be able to transfer it.

## Module Interface

#### Transaction

| Command                                                   | Description                                                                                                                                                            |
| ----------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [`publish-collection`](#publish-collection)               | Publishes [`denom@NFT Module`](../../readme.md#denom) from NFT module for sale with specified royalties, creates [`Collection`](#collection) in the marketplace store. |
| [`create-collection`](#create-collection)                 | Creates [`denom@NFT Module`](../../readme.md#denom) from NFT module for sale with specified royalties, creates [`Collection`](#collection) in the marketplace store and verifies it |
| [`update-royalties`](#update-royalties)                   | Update royalties of already published collection. |
| [`publish-nft`](#publish-nft)                             | Publish [`NFT@NFT Module`](../../readme.md#nft) from NFT module for sale with given price, creates [`NFT`](#nft) in marketplace store.                                 |
| [`update-price`](#update-price)                           | Updates price of already published or sale NFT.
| [`mint-nft`](#mint-nft)                                   | Mint [`NFT@NFT Module`](../../readme.md#nft) via NFT module, state of marketplace is not affected anyhow.                                                              |
| [`buy-nft`](#buy-nft)                                     | Buy [`NFT@NFT Module`](../../readme.md#nft) and removes the [`NFT`](#nft) from marketplace store.                                                                      |
| [`remove-nft`](#remove-nft)                               | Remove [`NFT`](#nft) from marketplace store.                                                                                                                           |
| [`verify-collection`](#verify-collection)                 | Verify [`Collection`](#collection).                                                                                                                                    |
| [`unverify-collection`](#unverify-collection)             | Unverify [`Collection`](#collection).                                                                                                                                  |
| [`add-admin`](#add-admin)                                 | Add admin
| [`remove-admin`](#remove-admin)                           | Remove admin

#### Query

| Command                                              | Description                                   |
| -----------------------------------------------------|-----------------------------------------------|
| [`list-collections`](#list-collections)              | Queries all [`Collection`](#collection)       |
| [`list-nfts`](#list-nfts)                            | Queries all  [`NFT`](#nft)                    |
| [`show-collection`](#show-collection)                | Show [`Collection`](#collection) by Id        |
| [`show-nft`](#show-nft)                              | Show [`NFT`](#nft) by Id                      |
| [`collection-by-denom-id`](#collection-by-denom-id)  | Show [`Collection`](#collection) by denom Id  |
| [`list-admins`](#list-admins)                        | Queries all marketplace admins                |

## Object types:

### Collection

References denom created by NFT module and has details about royalties on mint and resale for all NFTs traded from the referenced denom.

```go
type Collection struct {
	Id              uint64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	DenomId         string    `protobuf:"bytes,2,opt,name=denomId,proto3" json:"denomId,omitempty"`
	MintRoyalties   []Royalty `protobuf:"bytes,3,rep,name=mintRoyalties,proto3" json:"mintRoyalties"`
	ResaleRoyalties []Royalty `protobuf:"bytes,4,rep,name=resaleRoyalties,proto3" json:"resaleRoyalties"`
	Verified        bool      `protobuf:"varint,5,opt,name=verified,proto3" json:"verified,omitempty"`
	Owner           string    `protobuf:"bytes,6,opt,name=owner,proto3" json:"owner,omitempty"`
}

type Royalty struct {
	Address string                                 `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Percent github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=percent,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"percent"`
}
```

### NFT

References nft created by NFT module and has details about its sale price.

```go
type Nft struct {
	Id      uint64     `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	TokenId string     `protobuf:"bytes,2,opt,name=tokenId,proto3" json:"tokenId,omitempty"`
	DenomId string     `protobuf:"bytes,3,opt,name=denomId,proto3" json:"denomId,omitempty"`
	Price   types.Coin `protobuf:"bytes,4,opt,name=price,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"price"`
	Owner   string     `protobuf:"bytes,5,opt,name=owner,proto3" json:"owner,omitempty"`
}
```

### Events
> The events that are emitted after certain operations
```go
	EventPublishCollectionType       = "publish_collection"
	EventPublishNftType              = "publish_nft"
	EventBuyNftType                  = "buy_nft"
	EventMintNftType                 = "mint_nft"
	EventRemoveNftType               = "remove_nft"
	EventVerifyCollectionType        = "verify_collection"
	EventUnverifyCollectionType      = "unverify_collection"
  EventCreateCollectionType        = "create_collection"
	EventUpdateRoyaltiesType         = "update_royalties"
	EventUpdatePriceType             = "update_price"
	EventAddAdminType                = "add_admin"
	EventRemoveAdminType             = "remove_admin"
```

## Full commands info

### Transactions

### `publish-collection`

> Publish collection for sale with optional royalties. Only owner of collection can publish it for sale.

- arguments:
  - `denom-id` `string` `Denom to publish for sale` `required: true`
- flags:
  - `--mint-royalties` `string` `Royalties that will be distributed when NFTs are minted on demand via the marketplace.` `required: false`
  - `--resale-royalties` `string` `Royalties that will be distributed when reselling NFTs on the marketplace.` `required: false`

Royalties are represented in the format `"address1:percent,address2:percent"`. For resale royalties first royalties are paid and whatever is left is paid to the seller. If there are no royalties set, the full amount is paid to the seller. Mint royalties are required to sum to 100% because for some cases we could have one owner of collection onchain that manages it and someone else who should receive the bigger part of the amount (ex. CUDOS Markets).

```bash
cudos-noded tx marketplace publish-collection <denom-id> --mint-royalties="cudos1kztarv5vlckzt5j7z3y5u0dj6q6q6axyz4pe60;0.01,cudos14vjzkqs505xvs4tp3kdkzq3mzh6vutngnlqamz:11.22" --resale-royalties="cudos18687hmplu9mfxr47um0adne6ml29turydgm64j:50" --keyring-backend=<keyring> --chain-id=<chain-id> --gas=auto --gas-adjustment=1.3 --gas-prices=5000000000000acudos --from=<from-key>
```

### `create-collection`

> Creates denom, publishes for sale with optional royalties and can also verify it if executed by marketplace admin.

- arguments:
  - `denom-id` `string` `Denom to publish for sale` `required: true`
- flags:
  - `--name` `string` `The unique name of the denom.` `required: true`
  - `--symbol` `string` `The unique symbol of the denom.` `required: true`
  - `--from` `string` `The address that is issuing the denom. Will be set as denom creator. Can be either an address or alias to that address` `required: true`
  - `--schema` `string` `Metadata about the NFT. Schema-content or path to schema.json.` `required: false`
  - `--traits` `string` [`Traits`](../nft/types/traits.go)` that define the denom behavior and restrictions` `required: false`
  - `--description` `string` `Text description of the denom` `required: false`
  - `--data` `string` `Denom metadata` `required: false`
  - `--minter` `string` `Address that will be allowed to mint NFTs from this denom other than the owner` `required: false`
  - `--mint-royalties` `string` `Royalties that will be distributed when NFTs are minted on demand via the marketplace.` `required: false`
  - `--resale-royalties` `string` `Royalties that will be distributed when reselling NFTs on the marketplace.` `required: false`
  - `--verified` `bool` `Specifies if collection is verified - can be set to true only by admins` `required: false`

```bash
cudos-noded tx marketplace create-collection <denom-id> --name=<denom-name> --symbol=<denom-symbol> --verified=true --mint-royalties="cudos1kztarv5vlckzt5j7z3y5u0dj6q6q6axyz4pe60;0.01,cudos14vjzkqs505xvs4tp3kdkzq3mzh6vutngnlqamz:11.22" --resale-royalties="cudos18687hmplu9mfxr47um0adne6ml29turydgm64j:50" --keyring-backend=<keyring> --chain-id=<chain-id> --gas=auto --gas-adjustment=1.3 --gas-prices=5000000000000acudos --from=<from-key>
```
  
### `publish-nft`

> Publish NFT for sale with price. Only owner, approved operator or approved address can publish it for sale.

- arguments:
  - `token-id` `string` `Token id to publish for sale` `required: true`
  - `denom-id` `string` `Denom to which the token id belongs` `required: true`
  - `price` `string` `Price for which to publish the NFT for sale` `required: true`
- flags:
  none

```bash
cudos-noded tx marketplace publish-nft <token-id> <denom-id> 10000000acudos --keyring-backend=<keyring> --chain-id=<chain-id> --gas=auto --gas-adjustment=1.3 --gas-prices=5000000000000acudos --from=<from-key>
```

### `mint-nft`

> Mint NFT via the NFT module and distribute royalties if any.

- arguments:
  - `denom-id` `string` `Denom from which to mint the token` `required: true`
  - `recipient` `string` `Recipient to receive the minted NFT` `required: true`
  - `price` `string` `Amount to be paid to mint the NFT` `required: true`
  - `token-name` `string` `Name of token to be minted` `required: true`
- flags:
  - `--uri` `string` `Uri for NFT metadata stored offchain.` `required: false`
  - `--data` `string` `NFT metdata stored onchain.` `required: false`

```bash
cudos-noded tx marketplace mint-nft <denom-id> <recipient> 11111acudos <token-name> --uri=<token-uri> --data=<token-data> --keyring-backend=<keyring> --chain-id=<chain-id> --gas=auto --gas-adjustment=1.3 --gas-prices=5000000000000acudos --from=<from-key>
```

### `buy-nft`

> Buy NFT published for sale.

- arguments:
  - `nft-id` `string` `Nft id in the marketplace` `required: true`
- flags:
  none

```bash
cudos-noded tx marketplace buy-nft <nft-id> --keyring-backend=<keyring> --chain-id=cudos-local-network --gas=auto --gas-adjustment=1.3 --gas-prices=5000000000000acudos --from=<from-key>
```

### `remove-nft`

> Remove NFT from marketplace that was previously published for sale. Only owner of the NFT can remove it.

- arguments:
  - `nft-id` `string` `Nft id in the marketplace` `required: true`
- flags:
  none

```bash
cudos-noded tx marketplace remove-nft <nft-id> --keyring-backend=<keyring> --chain-id=cudos-local-network --gas=auto --gas-adjustment=1.3 --gas-prices=5000000000000acudos --from=<from-key>
```

### `verify-collection`

> Verify collection in the marketplace.

- arguments:
  - `collection-id` `string` `Collection id in the marketplace` `required: true`
- flags:
  none

```bash
cudos-noded tx marketplace verify-collection <collection-id> --keyring-backend=<keyring> --chain-id=cudos-local-network --gas=auto --gas-adjustment=1.3 --gas-prices=5000000000000acudos --from=<from-key>
```

### `unverify-collection`

> Unverify collection in the marketplace.

- arguments:
  - `collection-id` `string` `Collection id in the marketplace` `required: true`
- flags:
  none

```bash
cudos-noded tx marketplace unverify-collection <collection-id> --keyring-backend=<keyring> --chain-id=cudos-local-network --gas=auto --gas-adjustment=1.3 --gas-prices=5000000000000acudos --from=<from-key>
```

### `transfer-admin-permission`

> Transfer marketplace admin permission to different address.

- arguments:
  - `new-admin` `string` `New marketplace admin address` `required: true`
- flags:
  none

```bash
cudos-noded tx marketplace transfer-admin-permission <new-admin> --keyring-backend=<keyring> --chain-id=cudos-local-network --gas=auto --gas-adjustment=1.3 --gas-prices=5000000000000acudos --from=<from-key>
```

### `update-royalties`

> Update collection royalties.

- arguments:
  - `collection-id` `string` `Collection id in the marketplace` `required: true`
- flags:
  - `--mint-royalties` `string` `Royalties that will be distributed when NFTs are minted on demand via the marketplace.` `required: false`
  - `--resale-royalties` `string` `Royalties that will be distributed when reselling NFTs on the marketplace.` `required: false`

```bash
cudos-noded tx marketplace update-royalties <collection-id> --mint-royalties="cudos18x9glvtqk0x43xnjdx7w9lzqm0ganc950ur8n5:50" --resale-royalties="cudos18x9glvtqk0x43xnjdx7w9lzqm0ganc950ur8n5:50" --keyring-backend=<keyring> --chain-id=cudos-local-network --gas=auto --gas-adjustment=1.3 --gas-prices=5000000000000acudos --from=<from-key>
```

### `update-price`

> Update NFT price.

- arguments:
  - `nft-id` `string` `NFT id in the marketplace` `required: true`
  - `price` `string` `New price for the NFT` `required: true`
- flags:
  none

```bash
cudos-noded tx marketplace update-price <nft-id> <price> --keyring-backend=<keyring> --chain-id=cudos-local-network --gas=auto --gas-adjustment=1.3 --gas-prices=5000000000000acudos --from=<from-key>
```

### `add-admin`

> Add admin.

- arguments:
  - `address` `string` `Address to join the admin set` `required: true`
- flags:
  none

```bash
cudos-noded tx marketplace add-admin <admin-address> --keyring-backend=<keyring> --chain-id=cudos-local-network --gas=auto --gas-adjustment=1.3 --gas-prices=5000000000000acudos --from=<from-key>
```

### `remove-admin`

> Remove admin.

- arguments:
  - `address` `string` `Address to be removed from the admin set` `required: true`
- flags:
  none

```bash
cudos-noded tx marketplace remove-admin <admin-address> --keyring-backend=<keyring> --chain-id=cudos-local-network --gas=auto --gas-adjustment=1.3 --gas-prices=5000000000000acudos --from=<from-key>
```

### Queries

### `list-collections`

> List all collections published for sale.

- arguments:
  - none
- flags:
  - none

```bash
cudos-noded query marketplace list-collections
```

### `list-nfts`

> List all NFTs published for sale.

- arguments:
  - none
- flags:
  - none

```bash
cudos-noded query marketplace list-nfts
```

### `show-collection`

> Get collection published for sale by its Id in the marketplace.

- arguments:
  - `collection-id` `string` `Collection id in the marketplace` `required: true`
- flags:
  none

```bash
cudos-noded query marketplace show-collection <collection-id>
```

### `show-nft`

> Get NFT published for sale by its Id in the marketplace.

- arguments:
  - `nft-id` `string` `Nft id in the marketplace` `required: true`
- flags:
  none

```bash
cudos-noded query marketplace show-nft <nft-id>
```

### `collection-by-denom-id`

> Get collection published for sale by its denom Id.

- arguments:
  - `denom-id` `string` `Denom Id of the collection` `required: true`
- flags:
  none

```bash
cudos-noded query marketplace collection-by-denom-id <denom-id>
```

### `list-admins`

> List admins.

- arguments:
  - none
- flags:
  - none

```bash
cudos-noded query marketplace list-admins
```