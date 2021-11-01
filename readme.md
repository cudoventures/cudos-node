# Cudos Testnet

## Prerequirements
Starport: https://github.com/tendermint/starport/releases (linux only)<br />
or<br />
make (different ways to install it depending on OS) 


## Requirements
This project uses the default blog project from starport. It's purpose is to prove POC 1 requirements, such as:
 - Ability to launch a "vanilla" Cosmos blockchain and create a closed public test environment. 
 - Ability of the network to support Tendermint BFT with Ethereum type of wallets. 
 - The network should start with ~20 accounts/validators that have pre-configured vested balance in the genesis block.

## Links:
 - https://docs.cosmos.network/master/modules/auth/
 - https://docs.cosmos.network/master/basics/accounts.html
 - https://docs.cosmos.network/master/modules/auth/05_vesting.html


## Ability to launch a "vanilla" Cosmos blockchain and create a closed public test environment.
Launching a vanilla Cosmos blockchain is possible by starting this project using methods below.
  
## Ability of the network to support Tendermint BFT with Ethereum type of wallets
The network supports Tendermint BFT by default. The wallets private keys are generated using secp256k1 also by default. The retionale of using Ethereum type of wallets is to ensure that the users will be able to import their ethereum wallets into cudos blockchain using theirs seed phase. After the import the users will expect to see that their balance from ethereum blockchain is transferred to cudos. Although the cryptography is the same it is used in a slightly different manner so a converted is developed. Its usage is described below. It can convert ethereum public key to cudos wallet address. Using cudos wallet address a wallet can be pre-funded with required tokes so when a user import his wallet, using his seed, the balance will be correct.

## Creating accounts/validators with preconfigured vested balance in the genesis block
There are three ways to add preconfigured accounts with/without vested balance in the genesis block.
1. Modifying genesis.json after the blockchian is initialized, but before it is first started. This is not recommented method.
2. Using config.yml. This method works if the blockchain is initialized with starport utility.
3. Using commands from the binary itself. This method works if the blockchain is manually initialized without startport utility.

<br />
<br />
<br />

# Manual build
Build the blockchian binary into $GOPATH directory using "cudos-noded" name. All these steps are combined into init.sh/init.cmd 

    make

Initialize the blockchain.

    cudos-noded init cudos-node --chain-id=cudos-node

Creating accounts. Each account can have vesting balance as well.

    cudos-noded keys add validator01 --keyring-backend test --vesting-amount 500stake --vesting-end-time 1617615300

Add balance in the genesis block to an account.

    cudos-noded add-genesis-account $MY_VALIDATOR_ADDRESS 100000000000stake

Add validator 

    cudos-noded gentx validator01 100000000stake --chain-id cudos-node --keyring-backend test

Collect genesis transaction and start the blockchain

    cudos-noded collect-gentxs
    cudos-noded start

<br />
<br />
<br />

# Starport build
Configure accounts and validators in config.yml after that just start the blockchain

    starport serve

<br />
<br />
<br />

# Docker build
1. Build persistent-node
cd ./docker
docker-compose -f ./persistent-node.yml -p cudos-network-persistent-node up --build

2. After node starts copy its it and paste it into full-node.yml
Peer node looks like:
P2P Node ID ID=de14a2005d220171c7133efb31b3f3e1d7ba776a file=/root/.blog/config/node_key.json module=p2p

3. Run full-node
cd ./docker
docker-compose -f ./full-node.yml -p cudos-network-full-node up --build

<br />
<br />
<br />

# Converting ethereum public keys to cosmos wallet address
Run the converter and pass a ethereum public key as argument.

    go run ./converter 0x03139bb3b92e99d034ee38674a0e29c4aad83dd09b3fa465a265da310f9948fbe6

<b>Example ethereum mnemonic:</b> battle erosion opinion city birth modify scale hood caught menu risk rather<br >
<b>Example ethereum public key (32 bytes, compressed form):</b> 0x03139bb3b92e99d034ee38674a0e29c4aad83dd09b3fa465a265da310f9948fbe6

This mnemonic could be imported into cudos blockchain in order to verify that resulting account access will be the same as generated from the converter.

    cudos-noded keys add ruser02 --recover --hd-path="m/44'/60'/0'/0/0"

<br />
<br />
<br />

# Usefull commands
## Send currency
    cudos-noded tx bank send $VALIDATOR_ADDRESS $RECIPIENT 51000000stake --chain-id=cudos-network --keyring-backend test

## Check balances
    cudos-noded query bank balances $RECIPIENT --chain-id=cudos-network

## Create validator
    cudos-noded tx staking create-validator --amount=1000stake \
    --from=validator02 \
    --pubkey=$(cudos-noded tendermint show-validator) \
    --moniker=cudos-node-01 \
    --chain-id=cudos-network \
    --commission-rate="0.10" \
    --commission-max-rate="0.20" \
    --commission-max-change-rate="0.01" \
    --min-self-delegation="1" \
    --gas="auto" \
    --gas-prices="2500acudos" \
    --gas-adjustment="1.80" \
    --keyring-backend test

<br />
<br />
<br />

# Reset the blockchain
All data of the blockchain is store at ~/.blog folder. By deleting it the entire blockchain is completely reset and it must be initialized again.

# Static compile
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-extldflags "-static"' ./cmd/cudos-noded/

export CGO_LDFLAGS="-lpthread -ldl"
go build -v -a -tags netgo,osusergo -ldflags='-lpthread -extldflags "-lpthread -static"' ./cmd/cudos-noded/


# NFT Module Specification


## Overview

A module for operating with Non-Fungible Tokens on the CUDOS network. The methods that are exposed by it are mainly based on [ERC721 interface](https://ethereum.org/en/developers/docs/standards/tokens/erc-721/) from the Ethereum network and not so much on the [CW-721](https://github.com/CosmWasm/cw-nfts) from the Cosmos network. The reason for this is that the main idea of the module is to transfer tokens through a [bridge](https://github.com/CudoVentures/cosmos-gravity-bridge) between CUDOS network and Ethereum and thus it is better to follow the ERC721 standard. 

## Module Interface
The module gives the user the ability to either write(via transaction) or read(via query) to/from the network.

### The following commands are available (click on them for further info)

#### Transaction
| Command                                               | Description                                                                                                                            |
| ----------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------- |
| [`issue`](#issue)                           | Issues a new [`denomination`](#Denom) to the specified owner                                                   |
| [`mint`](#mint)                             | Mints a new [`NFT`](#NFT) in a given denomination to the specified owner                                                                       |
| [`edit`](#edit)             | Edits an already existing [`NFT`](#NFT)  |
| [`transfer`](#transfer)                        | Transfers an existing NFT from one owner to another                                                                                                   |
| [`burn`](#burn)                  | Burns the specified [`NFT`](#NFT) . Only the owner can burn the NFT                                                                                        |
| [`approve`](#approve)                        |  Adds an approved operator that can transfer the [`NFT`](#NFT)                                                                                                |
| [`revoke`](#revoke)                 | Removes an approved operated for the [`NFT`](#NFT)
| [`approveAll`](#approveall)                 | Approves an operator on user level - the operator can transfer all of the user tokens

#### Query

| Command                                               | Description                                                                                                                            |
| ----------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------- |
| [`denom`](#denom)                           | Queries for a [`denomination`](#Denom) by denomination Id                                                  |
| [`denom-by-name`](#denom-by-name)                             | Queries for a  [`denomination`](#Denom) by denomination name                                                                  |
| [`denoms`](#denoms)             | Query for all denominations of all collections of NFTs  |
| [`collection`](#collection)                        | Get all the NFTs from a given [`collection`](#Collections).                                                                                                 |
| [`supply`](#supply)                  | Returns the total supply of a collection or owner of NFTs.                                                                                       |
| [`owner`](#owner)                        |  Queries for the [`owner`](#Owners) and returns the NFTs owned by an account address.                                                                                               |
| [`token`](#token)                 | Query a single [`NFT`](#NFT) from a [`collection`](#Collections).
| [`approvals`](#approvals)                 | Get the approved addresses for the [`NFT`](#NFT)
| [`isApprovedForAll`](#isapprovedforall)                 | Gets whether the address is approved for all

## Full commands info:

### Transaction commands

### `issue`

> Issues a new denom that will be used for minting new NFTs. Only the denom creator can issue new NFTs

- arguments:
  - `denom-id` `string` `Unique Id that identifies the denom. Must be all lowercase` `required: true`
- flags: 
  - `--name` `string` `The unique name of the denom.` `required: true`
  - `--from` `string` `The address that is issuing the denom. Will be set as denom creator. Can be either an address or alias to that address` `required: true`
  - `--schema` `string` `Metadata about the NFT. Schema-content or path to schema.json.` `required: false`
  - `--chain-id` `string` `The name of the network.` `required`
  - `--fees` `string` `The specified fee for the operation` `required: true`

**Example:**

``` bash
$ cudos-noded tx nft issue <denom-id> --from=<key-name> --name=<denom-name> --schema=<schema-content or path to schema.json> --chain-id=<chain-id> --fees=<fee>
```

### `mint`

> Mint a NFT and set the owner to the recipient. Only the denom creator can mint a new NFT

- arguments:
    - `denom-id` `string` `The denomId that this NFT will be associated` `required: true`
    - `token-id` `string` `Unique Id that identifies the token. Must be all lowercase` `required: true`
- flags:
    - `--from` `string` `The address that is minting the NFT. Must be denom creator. Can be either an address or alias to that address` `required: true`
    - `--recipient` `string` `The user(owner) that will receive the NFT` `required: true`
    - `--uri` `string` `The URI of the NFT.` `required: false`
    - `--chain-id` `string` `The name of the network.` `required: true`
    - `--fees` `string` ` The specified fee for the operation` `required: true`

**Example:**

``` bash
$ cudos-noded tx nft mint <denom-id> <token-id> --recipient=<recipient> --from=<key-name> --uri=<uri> --chain-id=<chain-id> --fees=<fee>

```

### `edit`

> Edit an NFT - can change name, uri or data. Only the owner can edit the NFT.

- arguments:
  - `denom-id` `string` `The denomId of the edited NFT` `required: true`
  - `token-id` `string` `Unique Id that identifies the token. Must be all lowercase` `required: true`
- flags:
  - `--from` `string` `The address that is editing the NFT. Can be either an address or alias to that address` `required: true`
  - `--uri` `string` `The URI of the NFT.` `required: false`
  - `--chain-id` `string` `The name of the network.` `required: true`
  - `--fees` `string` `The specified fee for the operation` `required: true`

**Example:**

``` bash
$ cudos-noded tx nft edit <denom-id> <token-id>  --from=<key-name> --uri=<uri> --chain-id=<chain-id> --fees=<fee>
```

### `burn`

> Burns the NFT - deletes it permanently

- arguments:
  - `denom-id` `string` `The denomId of the edited NFT` `required: true`
  - `token-id` `string` `Unique Id that identifies the token. Must be all lowercase` `required: true`
- flags:
  - `--from` `string` `The address that is editing the NFT. Can be either an address or alias to that address` `required: true`
  - `--chain-id` `string` `The name of the network.` `required: true`
  - `--fees` `string` `The specified fee for the operation` `required: true`

**Example:**

``` bash
$ cudos-noded tx nft burn <denom-id> <token-id> --from=<key-name> --chain-id=<chain-id> --fees=<fee>
```

### `transfer`

> Transfer an NFT - from one owner to another The sender must be either the owner, approved address on NFT or globally approved operator.

- arguments:
  - `from` `string` `The address of the NFT owner` `required: true`
  - `to` `string` `The address of the user that will receive the NFT` `required: true`
  - `denom-id` `string` `The denomId of the edited NFT` `required: true`
  - `token-id` `string` `Unique Id that identifies the token. Must be all lowercase` `required: true`
- flags:
  - `--from` `string` `The address that is requesting the transfer of the NFT. Can be either an address or alias to that address. must be either the owner, approved address on NFT or globally approved operator.` `required: true`
  - `--uri` `string` `The URI of the NFT.` `required: false`
  - `--chain-id` `string` `The name of the network.` `required: true`
  - `--fees` `string` `The specified fee for the operation` `required: true`

**Example:**

``` bash
$ cudos-noded tx nft transfer <from> <to> <denom-id> <token-id>  --from=<key-name> --uri=<uri> --chain-id=<chain-id> --fees=<fee> 
```

### `approve`

> Adds an address to the approved list. Approved address on NFT level can transfer the NFT from one owner to another. Approved addresses for the NFT are cleared upon transfer.

- arguments:
  - `approvedAddress` `string` `The address that will be approved` `required: true`
  - `denom-id` `string` `The denomId of the edited NFT` `required: true`
  - `token-id` `string` `Unique Id that identifies the token. Must be all lowercase` `required: true`
- flags:
  - `--from` `string` `The address that is requesting the approval. Can be either an address or alias to that address. must be either the owner  or globally approved operator.` `required: true`
  - `--chain-id` `string` `The name of the network.` `required: true`
  - `--fees` `string` `The specified fee for the operation` `required: true`

**Example:**

``` bash
$ cudos-noded tx nft approve <approvedAddress> <denom-id> <token-id> --from=<key-name> --chain-id=<chain-id> --fees=<fee>
```

### `revoke`

> Removes the address from the approved list. Approved address on NFT level can transfer the nft from one owner to another. Approved addresses for the NFT are cleared upon transfer.

- arguments:
  - `addressToRevoke` `string` `The address that will be removed` `required: true`
  - `denom-id` `string` `The denomId of the edited NFT` `required: true`
  - `token-id` `string` `Unique Id that identifies the token. Must be all lowercase` `required: true`
- flags:
  - `--from` `string` `The address that is requesting the removal of approval. Can be either an address or alias to that address. Must be either the owner  or globally approved operator.` `required: true`
  - `--chain-id` `string` `The name of the network.` `required: true`
  - `--fees` `string` `Тhe specified fee for the operation` `required: true`

**Example:**

``` bash
$ cudos-noded tx nft revoke <addressToRevoke> <denom-id> <token-id>--uri=<uri> --from=<key-name> --chain-id=<chain-id> --fees=<fee>
```

### `approveAll`

> Adds the address to the approved operator list for the user. Approved address on user level can transfer the nft from one owner to another. The address is automatically added to the msg.sender(--from) approved list

- arguments:
  - `operator` `string` `The address that will be approved` `required: true`
  - `approved` `string` `Boolean value indicating if the addres is approved: can be true or false` `required: true`
- flags:
  - `--from` `string` `The address that is requesting the approval. The approved address will be able to handle the transfers of --from assets. Can be either an address or alias to that address. must be either the owner  or globally approved operator.` `required: true`
  - `--chain-id` `string` `The name of the network.` `required: true`
  - `--fees` `string` `The specified fee for the operation` `required: true`

**Example:**

``` bash
$ cudos-noded tx nft approveAll <operator> <true/false> --from=<key-name> --chain-id=<chain-id> --fees=<fee>
```

### Query commands

### `denom`

> Query the denom by the specified denom id.

- arguments:
  - `denom-id` `string` `The denomId to search for` `required: true`
- flags:
  - none
  
**Example:**

``` bash
$ cudos-noded query nft denom <denomId>
```

### `denom-by-name`

> Query the denom by the specified denom name.

- arguments:
  - `denom-name` `string` `The denom name to search for` `required: true`
- flags:
  - none

**Example:**

``` bash
$ cudos-noded query nft denom <denomName>
```

### `denoms`

> Query all denominations of all collections of NFTs.

- arguments:
  - none
- flags:
  - none

**Example:**

``` bash
$ cudos-noded query nft denoms
```

### `collection`

> Query all denominations of all collections of NFTs.

- arguments:
  - `denom-id`: `The id of the denomination collection.` `required:true`
- flags:
  - none

**Example:**

``` bash
$ cudos-noded query nft collection <denom-id>
```

### `supply`

> Gets the total supply of a collection or owner of NFTs.

- arguments:
  - `denom-id`: `The id of the denomination collection.` `required:true`
- flags:
  - none

**Example:**

``` bash
$ cudos-noded query nft supply <denom-id>
```

### `owner`

> Get the NFTs owned by an account address.

- arguments:
  - `address`: `The address of the owner.` `required:true`
- flags:
  - `--denom-id`: `The id of the denom` `required:true`
  
**Example:**
``` bash
$ cudos-noded query nft owner <address> --denom-id=<denom-id>
```

### `token`

> Query a single NFT from a collection.

- arguments:
  - `denom-id`: `The id of the denom collection` `required:true`
  - `token-id`: `The id of the NFT` `required:true`
- flags:
  - none

**Example:**

``` bash
$ cudos-noded query nft token <denom-id> <token-id>
```

### `approvals`

> Get the approved addresses for the NFT

- arguments:
  - `denom-id`: `The id of the denom collection` `required:true`
  - `token-id`: `The id of the NFT` `required:true`
- flags:
  - none

**Example:**

``` bash
$ cudos-noded query nft approvals <denomId> <tokenId>
```


### `isApprovedForAll`

> Query if an address is an authorized operator for another address

- arguments:
  - `owner`: `The owner addresses to search` `required:true`
  - `operator`: `The operator address to be searched for` `required:true`
- flags:
  - none

**Example:**

``` bash
$ cudos-noded query nft approvals <owner> <operator>
```

## Object types:

### NFT
```go
// NFT non fungible token interface
type NFT interface {
    GetID() string              // unique identifier of the NFT
    GetName() string            // return the name of BaseNFT
    GetOwner() sdk.AccAddress   // gets owner account of the NFT
    GetURI() string             // tokenData field: URI to retrieve the of chain tokenData of the NFT
    GetData() string            // return the Data of BaseNFT
    GetApprovedAddresses() map[string]bool// returns the approved addresses of BaseNFT

}
```

### NFT implementation
```go
type BaseNFT struct {
	Id                string          `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name              string          `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	URI               string          `protobuf:"bytes,3,opt,name=uri,proto3" json:"uri,omitempty"`
	Data              string          `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
	Owner             string          `protobuf:"bytes,5,opt,name=owner,proto3" json:"owner,omitempty"`
	ApprovedAddresses map[string]bool `protobuf:"bytes,6,rep,name=approvedAddresses,proto3" json:"approvedAddresses,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}
```

## Collections

>As all NFTs belong to a specific `Collection` under `{denomID}/{tokenID}`

```go
// Collection of non fungible tokens
type Collection struct {
    Denom Denom     `json:"denom"`  // Denom of the collection; not exported to clients
    NFTs  []BaseNFT `json:"nfts"`   // NFTs that belongs to a collection
}
```

## Owners

>Owner holds the address of the user and his collection of NFTs

```go
// Owner of non fungible tokens
type Owner struct {
    Address       string            `json:"address"`
    IDCollections []IDCollection    `json:"id_collections"`
}
```

## IDCollection
>IDCollection holds the denomId and the Ids of the NFTs(insted of the full object)

```go
// IDCollection of non fungible tokens
type IDCollection struct {
    DenomId string   `json:"denom_id"`
    TokenIds []string `json:"token_ids"`
}

```

## Denom
> The denomination is used to group NFTs under it
```go
// Denom defines a type of NFT
type Denom struct {
	Id      string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name    string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Schema  string `protobuf:"bytes,3,opt,name=schema,proto3" json:"schema,omitempty"`
	Creator string `protobuf:"bytes,4,opt,name=creator,proto3" json:"creator,omitempty"`
}
```

## Events
> The events that are emitted after certain operations
```go
	EventTypeIssueDenom    = "issue_denom"
	EventTypeTransferNft   = "transfer_nft"
	EventTypeApproveNft    = "approve_nft"
	EventTypeApproveAllNft = "approve_all_nft"
	EventTypeRevokeNft     = "revoke_nft"
	EventTypeEditNFT       = "edit_nft"
	EventTypeMintNFT       = "mint_nft"
	EventTypeBurnNFT       = "burn_nft"
```