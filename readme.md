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
| Command                             | Description                                                                           |
| ----------------------------------- | ------------------------------------------------------------------------------------- |
| [`issue`](#issue)                   | Issues a new [`denomination`](#Denom) to the specified owner                          |
| [`mint`](#mint)                     | Mints a new [`NFT`](#NFT) in a given denomination to the specified owner              |
| [`edit`](#edit)                     | Edits an already existing [`NFT`](#NFT)                                               |
| [`transfer`](#transfer)             | Transfers an existing NFT from one owner to another                                   |
| [`transfer-denom`](#transfer-denom) | Transfers an existing NFT classification from one owner to another                    |
| [`burn`](#burn)                     | Burns the specified [`NFT`](#NFT) . Only the owner can burn the NFT                   |
| [`approve`](#approve)               | Adds an approved operator that can transfer the [`NFT`](#NFT)                         |
| [`revoke`](#revoke)                 | Removes an approved operated for the [`NFT`](#NFT)                                    |
| [`approveAll`](#approveall)         | Approves an operator on user level - the operator can transfer all of the user tokens |

#### Query

| Command                                 | Description                                                                          |
| --------------------------------------- | ------------------------------------------------------------------------------------ |
| [`denom`](#denom)                       | Queries for a [`denomination`](#Denom) by denomination Id                            |
| [`denom-by-name`](#denom-by-name)       | Queries for a  [`denomination`](#Denom) by denomination name                         |
| [`denom-by-symbol`](#denom-by-symbol)   | Queries for a  [`denomination`](#Denom) by denomination symbol                       |
| [`denoms`](#denoms)                     | Query for all denominations of all collections of NFTs                               |
| [`collection`](#collection)             | Get all the NFTs from a given [`collection`](#Collections).                          |
| [`supply`](#supply)                     | Returns the total supply of a collection or owner of NFTs.                           |
| [`owner`](#owner)                       | Queries for the [`owner`](#Owners) and returns the NFTs owned by an account address. |
| [`token`](#token)                       | Query a single [`NFT`](#NFT) from a [`collection`](#Collections).                    |
| [`approvals`](#approvals)               | Get the approved addresses for the [`NFT`](#NFT)                                     |
| [`isApprovedForAll`](#isapprovedforall) | Gets whether the address is approved for all                                         |


## Usage from inside a CosmWasm smart contract
You can check how to use the module from a rust smart contract in the [`cudos-cosmwasm-bindings`](https://github.com/CudoVentures/cudos-cosmwasm-bindings)

## Full commands info:

### Transaction commands

### `issue`

> Issues a new denom that will be used for minting new NFTs. Only the denom creator can issue new NFTs

- arguments:
  - `denom-id` `string` `Unique Id that identifies the denom. Must be all lowercase` `required: true`
- flags: 
  - `--name` `string` `The unique name of the denom.` `required: true`
  - `--symbol` `string` `The unique symbol of the denom.` `required: true`
  - `--from` `string` `The address that is issuing the denom. Will be set as denom creator. Can be either an address or alias to that address` `required: true`
  - `--schema` `string` `Metadata about the NFT. Schema-content or path to schema.json.` `required: false`
  - `--chain-id` `string` `The name of the network.` `required`
  - `--fees` `string` `The specified fee for the operation` `required: true`

**Example:**

``` bash
$ cudos-noded tx nft issue <denom-id> --from=<key-name> --name=<denom-name> --symbol=<denom-symbol> --schema=<schema-content or path to schema.json> --chain-id=<chain-id> --fees=<fee>
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

### `transfer-denom`

> Transfers the ownership of the NFT classification to others. The sender must be the owner.

- arguments:
  - `recipient` `string` `The address of the new NFT classification owner` `required: true`
  - `denom-id` `string` `The denomId of the transferred NFT classification` `required: true`
- flags:
  - `--from` `string` `The address that is requesting the transfer of the NFT collection. Must be the owner.` `required: true`
  - `--chain-id` `string` `The name of the network.` `required: true`
  - `--fees` `string` `The specified fee for the operation` `required: true`

**Example:**

``` bash
$ cudos-noded tx nft transfer-denom <recipient> <denom-id> --from=<key-name> --chain-id=<chain-id> --fees=<fee>
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

### `denom-by-symbol`

> Query the denom by the specified denom symbol.

- arguments:
  - `symbol` `string` `The denom symbol to search for` `required: true`
- flags:
  - none

**Example:**

``` bash
$ cudos-noded query nft denom <symbol>
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

> Get the NFTs owned by an account address for a given denom.

- arguments:
  - `address`: `The address of the owner.` `required:true`
- flags:
  - `--denom-id`: `The id of the denom` `required:false`
  
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
$ cudos-noded query nft isApprovedForAll <owner> <operator>
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
	Symbol    string `protobuf:"bytes,2,opt,name=name,proto3" json:"symbol,omitempty"`
	Schema  string `protobuf:"bytes,3,opt,name=schema,proto3" json:"schema,omitempty"`
	Creator string `protobuf:"bytes,4,opt,name=creator,proto3" json:"creator,omitempty"`
}
```

## Events
> The events that are emitted after certain operations
```go
	EventTypeIssueDenom    = "issue_denom"
	EventTypeTransferNft   = "transfer_nft"
    EventTypeTransferDenom = "transfer_denom"
	EventTypeApproveNft    = "approve_nft"
	EventTypeApproveAllNft = "approve_all_nft"
	EventTypeRevokeNft     = "revoke_nft"
	EventTypeEditNFT       = "edit_nft"
	EventTypeMintNFT       = "mint_nft"
	EventTypeBurnNFT       = "burn_nft"
```

## API Endpoints
> default API local url: localhost:1317

All the requests/response below are used as an example, for the full capabilities and parameters, consult the full command specifications

### Transactions:

### Issue Denom : POST
http://localhost:1317/nft/nfts/denoms/issue

Request:
```json
{
  "owner": "test",
  "id": "testdenom",
  "name": "testname",
  "symbol": "testDenomSymbol",
  "base_req": {
    "from":"cudos1qy7a8qvmqtqrscz7rf9l3xlllm0l6x3xnmarze",
    "chain_id":"cudos-network"
  }
}
```

Response: 
```json
{
  "type": "cosmos-sdk/StdTx",
  "value": {
    "msg": [
      {
        "type": "github.com/CudoVentures/cudos-node/nft/MsgIssueDenom",
        "value": {
          "id": "testdenom",
          "name": "testname",
          "test_denom_symbol": "testDenomSymbol",
          "sender": "cudos1qy7a8qvmqtqrscz7rf9l3xlllm0l6x3xnmarze"
        }
      }
    ],
    "fee": {
      "amount": [],
      "gas": "200000"
    },
    "signatures": [],
    "memo": "",
    "timeout_height": "0"
  }
}
```

### Mint NFT : POST
http://localhost:1317/nft/nfts/mint

Request:
```json
{
  "denom_id": "testdenom",
  "name": "testTokenName",
  "uri": "testuri",
  "data": "testdata",
  "recipient": "cudos1s609vqsnwxpm2t4scjq70770kph7uaz53lg89v",
  "base_req": {
    "from":"cudos1qy7a8qvmqtqrscz7rf9l3xlllm0l6x3xnmarze",
    "chain_id":"cudos-network"
  }
}
```

Response:
```json
{
  "type": "cosmos-sdk/StdTx",
  "value": {
    "msg": [
      {
        "type": "github.com/CudoVentures/cudos-node/nft/MsgMintNFT",
        "value": {
          "denom_id": "testdenom",
          "name": "testTokenName",
          "uri": "testuri",
          "data": "testdata",
          "sender": "cudos1qy7a8qvmqtqrscz7rf9l3xlllm0l6x3xnmarze",
          "recipient": "cudos1s609vqsnwxpm2t4scjq70770kph7uaz53lg89v"
        }
      }
    ],
    "fee": {
      "amount": [],
      "gas": "200000"
    },
    "signatures": [],
    "memo": "",
    "timeout_height": "0"
  }
}
```

### Edit NFT : PUT
http://localhost:1317/nft/nfts/edit/{denomId}/{tokenId}

Request:
```json
{
  "uri": "testuri",
  "data": "testdata",
  "name": "testname",
  "base_req": {
    "from":"cudos1qy7a8qvmqtqrscz7rf9l3xlllm0l6x3xnmarze",
    "chain_id":"cudos-network"
  }
}
```

Response:
```json
{
  "type": "cosmos-sdk/StdTx",
  "value": {
    "msg": [
      {
        "type": "github.com/CudoVentures/cudos-node/nft/MsgEditNFT",
        "value": {
          "id": "1",
          "denom_id": "testdenom",
          "name": "testname",
          "uri": "testuri",
          "data": "testdata",
          "sender": "cudos1qy7a8qvmqtqrscz7rf9l3xlllm0l6x3xnmarze"
        }
      }
    ],
    "fee": {
      "amount": [],
      "gas": "200000"
    },
    "signatures": [],
    "memo": "",
    "timeout_height": "0"
  }
}
```

### Transfer NFT : POST
http://localhost:1317/nft/nfts/transfer/{denomId}/{tokenId}

Request:
```json
{
  "from": "cudos1s609vqsnwxpm2t4scjq70770kph7uaz53lg89v",
  "to": "cudos18vpe7dfn6038ceae0ndlxdyuvgafrk2y6klzkx",
  "recipient": "cudos1s609vqsnwxpm2t4scjq70770kph7uaz53lg89v",
  "base_req": {
    "from":"cudos1qy7a8qvmqtqrscz7rf9l3xlllm0l6x3xnmarze",
    "chain_id":"cudos-network"
  }
}
```

Response:
```json
{
  "type": "cosmos-sdk/StdTx",
  "value": {
    "msg": [
      {
        "type": "github.com/CudoVentures/cudos-node/nft/MsgTransferNft",
        "value": {
          "denom_id": "testdenom",
          "token_id": "1",
          "from": "cudos1s609vqsnwxpm2t4scjq70770kph7uaz53lg89v",
          "to": "cudos18vpe7dfn6038ceae0ndlxdyuvgafrk2y6klzkx",
          "sender": "cudos1qy7a8qvmqtqrscz7rf9l3xlllm0l6x3xnmarze"
        }
      }
    ],
    "fee": {
      "amount": [],
      "gas": "200000"
    },
    "signatures": [],
    "memo": "",
    "timeout_height": "0"
  }
}
```

### Transfer Denom : POST
http://localhost:1317/nft/nfts/denoms/transfer/{denomId}

Request:
```json
{
  "recipient": "cudos1s609vqsnwxpm2t4scjq70770kph7uaz53lg89v",
  "denom-id": "testdenom",
  "base_req": {
    "from":"cudos1qy7a8qvmqtqrscz7rf9l3xlllm0l6x3xnmarze",
    "chain_id":"cudos-network"
  }
}
```

Response:
```json
{
    "type": "cosmos-sdk/StdTx",
    "value": {
        "msg": [
            {
                "type": "github.com/CudoVentures/cudos-node/nft/MsgTransferDenom",
                "value": {
                    "id": "whathever",
                    "sender": "cudos1detu83m7rd9ygvuzg3mee53sgwdae852ve5xav",
                    "recipient": "cudos1xlvmmvvmwvkcugkmjxc06du8qs2vrw337w5jda"
                }
            }
        ],
        "fee": {
            "amount": [],
            "gas": "200000"
        },
        "signatures": [],
        "memo": "",
        "timeout_height": "0"
    }
}
```

### Approve NFT : POST
http://localhost:1317/nft/nfts/approve/{denomId}/{tokenId}

Request:
```json
{
  "address_to_approve": "cudos1s609vqsnwxpm2t4scjq70770kph7uaz53lg89v",
  "base_req": {
    "from":"cudos1qy7a8qvmqtqrscz7rf9l3xlllm0l6x3xnmarze",
    "chain_id":"cudos-network"
  }
}
```

Response:
```json
{
  "type": "cosmos-sdk/StdTx",
  "value": {
    "msg": [
      {
        "type": "github.com/CudoVentures/cudos-node/nft/MsgApproveNft",
        "value": {
          "id": "1",
          "denom_id": "testdenom",
          "sender": "cudos1qy7a8qvmqtqrscz7rf9l3xlllm0l6x3xnmarze",
          "approvedAddress": "cudos1s609vqsnwxpm2t4scjq70770kph7uaz53lg89v"
        }
      }
    ],
    "fee": {
      "amount": [],
      "gas": "200000"
    },
    "signatures": [],
    "memo": "",
    "timeout_height": "0"
  }
}
```

### Revoke NFT : POST
http://localhost:1317/nft/nfts/revoke/{denomId}/{tokenId}

Request:
```json
{
  "address_to_revoke": "cudos1s609vqsnwxpm2t4scjq70770kph7uaz53lg89v",
  "base_req": {
    "from":"cudos1qy7a8qvmqtqrscz7rf9l3xlllm0l6x3xnmarze",
    "chain_id":"cudos-network"
  }
}
```

Response:
```json
{
  "type": "cosmos-sdk/StdTx",
  "value": {
    "msg": [
      {
        "type": "github.com/CudoVentures/cudos-node/nft/MsgRevokeNft",
        "value": {
          "addressToRevoke": "cudos1s609vqsnwxpm2t4scjq70770kph7uaz53lg89v",
          "denom_id": "testdenom",
          "token_id": "1",
          "sender": "cudos1qy7a8qvmqtqrscz7rf9l3xlllm0l6x3xnmarze"
        }
      }
    ],
    "fee": {
      "amount": [],
      "gas": "200000"
    },
    "signatures": [],
    "memo": "",
    "timeout_height": "0"
  }
}
```


### Burn NFT : POST
http://localhost:1317/nft/nfts/revoke/{denomId}/{tokenId}

Request:
```json
{
  "base_req": {
    "from":"cudos1qy7a8qvmqtqrscz7rf9l3xlllm0l6x3xnmarze",
    "chain_id":"cudos-network"
  }
}
```

Response:
```json
{
  "type": "cosmos-sdk/StdTx",
  "value": {
    "msg": [
      {
        "type": "github.com/CudoVentures/cudos-node/nft/MsgBurnNFT",
        "value": {
          "id": "1",
          "denom_id": "testdenom",
          "sender": "cudos1qy7a8qvmqtqrscz7rf9l3xlllm0l6x3xnmarze"
        }
      }
    ],
    "fee": {
      "amount": [],
      "gas": "200000"
    },
    "signatures": [],
    "memo": "",
    "timeout_height": "0"
  }
}
```

### Approve All : POST
http://localhost:1317/nft/nfts/revoke/approveAll

Request:
```json
{
  "approved": true,
  "approved_operator": "cudos1s609vqsnwxpm2t4scjq70770kph7uaz53lg89v",
  "base_req": {
    "from":"cudos1qy7a8qvmqtqrscz7rf9l3xlllm0l6x3xnmarze",
    "chain_id":"cudos-network"
  }
}
```

Response:
```json
{
  "type": "cosmos-sdk/StdTx",
  "value": {
    "msg": [
      {
        "type": "github.com/CudoVentures/cudos-node/nft/MsgApproveAllNft",
        "value": {
          "operator": "cudos1s609vqsnwxpm2t4scjq70770kph7uaz53lg89v",
          "sender": "cudos1qy7a8qvmqtqrscz7rf9l3xlllm0l6x3xnmarze",
          "approved": true
        }
      }
    ],
    "fee": {
      "amount": [],
      "gas": "200000"
    },
    "signatures": [],
    "memo": "",
    "timeout_height": "0"
  }
}
```

### Queries:

### Query Denom By Id : GET
http://localhost:1317/nft/denoms/{{denomId}}

Response:
```json
{
  "height": "4774",
  "result": {
    "denom": {
      "id": "testdenom",
      "name": "TESTDENOM",
      "schema": "testschema",
      "creator": "cudos13kkzjnz9t30dtkcevcvk2n2xu2n8mnzxuwnuur"
    }
  }
}
```

### Query Denom By Name: GET
http://localhost:1317/nft/denoms/name/{{denomName}}

Response:
```json
{
  "height": "4994",
  "result": {
    "denom": {
      "id": "testdenom1",
      "name": "testDenomNewName",
      "schema": "testschema",
      "creator": "cudos13kkzjnz9t30dtkcevcvk2n2xu2n8mnzxuwnuur"
    }
  }
}
```
### Query Denom By Symbol: GET
http://localhost:1317/nft/denoms/symbol/{{symbol}}

Response:
```json
{
  "height": "23",
  "result": {
    "denom": {
      "id": "testdenom",
      "name": "testName",
      "schema": "",
      "creator": "cudos1wye475erldt37cgj3kf4j35w24emhh0cdddg7z",
      "symbol": "testSymbol"
    }
  }
}
```

### Query All Denoms: POST
http://localhost:1317/nft/denoms/

Request:
```json
{}
```

Optional pagination:
```json
{ 
    "pagination": {
        "offset": "0",
        "limit":  "5",
        "count_total": true
    }
}
```

Response:
```json
{
  "height": "5061",
  "result": {
    "denoms": [
      {
        "id": "testdenom",
        "name": "TESTDENOM",
        "schema": "testschema",
        "creator": "cudos13kkzjnz9t30dtkcevcvk2n2xu2n8mnzxuwnuur"
      },
      {
        "id": "testdenom1",
        "name": "testDenomNewName",
        "schema": "testschema",
        "creator": "cudos13kkzjnz9t30dtkcevcvk2n2xu2n8mnzxuwnuur"
      }
    ],
    "pagination": {
      "total": "2"
    }
  }
}
```



### Query Collection for a Denom: POST
http://localhost:1317/nft/collection/{{denomId}}

Request: Pagination is optional
```json
  {
  "denom_id": "testdenom",
  "pagination": {
    "offset": "0",
    "limit":  "5",
    "count_total": true
  }
}
```


Response:
```json
{
  "height": "5189",
  "result": {
    "collection": {
      "denom": {
        "id": "testdenom",
        "name": "TESTDENOM",
        "schema": "testschema",
        "creator": "cudos13kkzjnz9t30dtkcevcvk2n2xu2n8mnzxuwnuur"
      },
      "nfts": [
        {
          "id": "1",
          "name": "testtoken",
          "uri": "",
          "data": "testData",
          "owner": "cudos1ztwjs6cp8t369l6ckgv5fg8vzleqn3qdkycll4"
        }
      ]
    },
    "pagination": {
      "total": "1"
    }
  }
}
```

### Query Supply for a Denom : GET
> Gets the total NFT count for a given denomId


http://localhost:1317/nft/collections/supply/{{denomId}}
 
> TODO: Must add pagination support to request and handle in the node

Response:
```json
{
  "height": "5221",
  "result": {
    "amount": "1"
  }
}
```

### Query Owner  GET
> Gets the NFTs for a given Owner ( Optional denomId )


http://localhost:1317/nft/owners/{{ownerAddress}}/{denomId}
Request:

```json
{
    "denom_id": "testdenom",
    "owner_address": "cudos1s6ncz2gyy0cgzzk5yctjx7yx7tyjzxnmnx9xlj",
    "pagination": {
        "offset": "1",
        "limit":  "5",
        "count_total": true
    }
}
```

Response:
```json
{
  "height": "10001",
  "result": {
    "owner": {
      "address": "cudos1ztwjs6cp8t369l6ckgv5fg8vzleqn3qdkycll4",
      "id_collections": [
        {
          "denom_id": "testdenom",
          "token_ids": [
            "1"
          ]
        }
      ]
    },
    "pagination": {
      "total": "1"
    }
  }
}
```

### Query NFT : GET
> Gets the NFT by a given denomId and tokenId

http://localhost:1317/nft/nfts/{{denomId}}/{{tokenId}}

Response:
```json
{
  "height": "10001",
  "result": {
    "nft": {
      "id": "1",
      "name": "testtoken",
      "uri": "",
      "data": "testData",
      "owner": "cudos1ztwjs6cp8t369l6ckgv5fg8vzleqn3qdkycll4"
    }
  }
}
```

Response with approved addresses:
```json
{
    "height": "161",
    "result": {
        "nft": {
            "id": "1",
            "name": "testtoken",
            "uri": "",
            "data": "testData",
            "owner": "cudos17najx40kq4f6yrslw6ggr5qm4meqs8p72jhv2f",
            "approved_addresses": {
                "cudos1yg8et80vfyjetyafqcpr8geyp2ypmnd3rk54z2": true
            }
        }
    }
}
```

### Query Approvals NFT : GET
> Gets the approvals for a NFT

http://localhost:1317/nft/approvals/{{denomId}}/{{tokenId}}

Response:
```json
{
  "height": "238",
  "result": {
    "approved_addresses": {
      "cudos1yg8et80vfyjetyafqcpr8geyp2ypmnd3rk54z2": true,
      "cudos1y43rjgknmk2hv3cpcu007crucwsgrv4n4rhmcs": true
    }
  }
}
```

### Query IsApprovedForAll : POST
> Gets the approvals for a NFT

http://localhost:1317/nft/isApprovedForAll

Request:
```json
{
        "owner": "cudos1y43rjgknmk2hv3cpcu007crucwsgrv4n4rhmcs",
        "operator": "cudos17najx40kq4f6yrslw6ggr5qm4meqs8p72jhv2f"
}
```

Response:
```json
{
  "height": "409",
  "result": {
    "is_approved": true
  }
}
```
