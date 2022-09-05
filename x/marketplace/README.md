# NFT Marketplace Module Specification

## Overview

A module for trading NFTs on the CUDOS network. The module supports publishing collection for sale which is optional step to set mint royalties and resale royalties 
which will be deducted from the payments for all the NFTs traded from the given collection.
If user tries to sell NFT that has no royalties the seller will be paid the full price.
The NFTs can be minted via the marketplace which allows distribution of royalties on mint.
When NFT is published for sale, it will be soft locked in the NFT module, so the seller will remain owner but wont be able to transfer it.

## Module Interface

#### Transaction

| Command                                          | Description                                                                                                                                                            |
| ------------------------------------------------ | -----------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [`publish-collection`](#publish-collection)      | Publishes [`denom@NFT Module`](../../readme.md#denom) from NFT module for sale with specified royalties, creates [`Collection`](#collection) in the marketplace store. |
| [`publish-nft`](#publish-nft)                    | Publish [`NFT@NFT Module`](../../readme.md#nft) from NFT module for sale with given price, creates [`NFT`](#nft) in marketplace store.                                 |
| [`mint-nft`](#mint-nft)                          | Mint [`NFT@NFT Module`](../../readme.md#nft) via NFT module, state of marketplace is not affected anyhow.                                                              |
| [`buy-nft`](#buy-nft)                            | Buy [`NFT@NFT Module`](../../readme.md#nft) and removes the [`NFT`](#nft) from marketplace store.                                                                      |
| [`remove-nft`](#remove-nft)                      | Remove [`NFT`](#nft) from marketplace store.                                                                                                            |

#### Query

| Command                                 | Description                                   |
| --------------------------------------- | ----------------------------------------------|
| [`list-collections`](#list-collections)       | Queries all [`Collection`](#collection) |
| [`list-nfts`](#list-nfts)                     | Queries all  [`NFT`](#nft)              |
| [`show-collection`](#show-collection)        | Show [`Collection`](#collection) by Id   |
| [`show-nft`](#show-nft)                      | Show [`NFT`](#nft) by Id                 |

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

## Full commands info

### Transactions

### `publish-collection`

> Publish collection for sale with optional royalties. Only owner of collection can publish it for sale.

- arguments:
  - `denom-id` `string` `Denom to publish for sale` `required: true`
- flags:
  - `--mint-royalties` `string` `The unique name of the denom.` `required: false`
  - `--resale-royalties` `string` `The unique symbol of the denom.` `required: false`

```bash
cudos-noded tx marketplace publish-collection <denom-id> --mint-royalties="cudos1kztarv5vlckzt5j7z3y5u0dj6q6q6axyz4pe60;0.01,cudos14vjzkqs505xvs4tp3kdkzq3mzh6vutngnlqamz:11.22" --resale-royalties="cudos18687hmplu9mfxr47um0adne6ml29turydgm64j:50" --keyring-backend=<keyring> --chain-id=<chain-id> --gas=auto --gas-adjustment=1.3 --gas-prices=5000000000000acudos --from=<from-key>
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
  - `--uri` `string` `The unique name of the denom.` `required: false`
  - `--data` `string` `The unique name of the denom.` `required: false`

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
