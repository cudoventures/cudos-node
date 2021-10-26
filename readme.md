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


0884 188 072

# NFT Module Specification


## Overview

A module for operating with Non-Fungible Tokens on the CUDOS network. The methods that are exposed by it are mainly based on [ERC721 interface](https://ethereum.org/en/developers/docs/standards/tokens/erc-721/) from the Ethereum network and not so much on the [CW-721](https://github.com/CosmWasm/cw-nfts) from the Cosmos network. The reason for this is that the main idea of the module is to transfer tokens through [GravityBridge](https://github.com/CudoVentures/cosmos-gravity-bridge) between CUDOS network and Ethereum and thus it is better to follow the ERC721 standard. 

## Module Interface
The module gives the user to either write(via transaction) or read(via query) to/from the network.

The following transaction commands are available (click on them for further info) :

| Command                                               | Description                                                                                                                            |
| ----------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------- |
| [`Issue Denom`](#issue)                           | Issues a new [`Denom`](#addDenomLink) to the specified owner                                                   |
| [`Mint NFT`](#mint)                             | Mints a new NFT to the specified owner                                                                       |
| [`Edit NFT`](#edit-nft)             | Edits an already existing NFT  |
| [`Transfer NFT`](#transfer-nft)                        | Transfers an existing NFT from one owner to another                                                                                                   |
| [`Burn NFT`](#burn-nft)                  | Burns the specified NFT                                                                                        |
| [`Approve NFT`](#cudos-run)                        |  Adds an approved operator that can transfer the NFT                                                                                                |
| [`Revoke NFT`](#cudos-keys)                 | Removes an approved operated for the NFT 
| [`Approve All`](#cudos-keys)                 | Approves an operator on user level - the operator can transfer all of the user tokens



## Full commands info

### `issue`

> Issues a new denom that will be used for minting new NFTs. Only the denom creator can issue new NFTs

- arguments:
  - `denom-id` `string` `Unique Id that identifies the denom. Must be all lowercase` `required: true`
- options: 
  - `--name` `string` `The unique name of the denom.` `required: true`
  - `--from` `string` `The address that is issuing the denom. Will be set as denom creator. Can be either an addresso or alias to that address` `required: true`
  - `--schema` `string` `!!!Insert what is schema here!!! Schema-content or path to schema.json.` `required: false`
  - `--chain-id` `string` `The name of the network.` `required`
  - `--fees` `string` `!!!Inert more info here!!! The specified fee for the operation` `required: true`

**Example:**

``` bash
$ cudos-noded tx nft issue <denom-id> --from=<key-name> --name=<denom-name> --schema=<schema-content or path to schema.json> --chain-id=<chain-id> --fees=<fee>
```
### `mint`

> Mints a new NFT to the specified owner. Only the denom creator can mint a NFT

- arguments:
    - `denom-id` `string` `Unique Id that identifies the denom. Must be all lowercase` `required: true`
- options:
    - `--name` `string` `The unique name of the denom.` `required: true`
    - `--from` `string` `The address that is issuing the denom. Will be set as denom creator. Can be either an addresso or alias to that address` `required: true`
    - `--schema` `string` `!!!Insert what is schema here!!! Schema-content or path to schema.json.` `required: false`
    - `--chain-id` `string` `The name of the network.` `required`
    - `--fees` `string` `!!!Inert more info here!!! The specified fee for the operation` `required: true`

**Example:**

``` bash
$ cudos-noded tx nft issue <denom-id> --from=<key-name> --name=<denom-name> --schema=<schema-content or path to schema.json> --chain-id=<chain-id> --fees=<fee>
```



## Create denom
```
$ <appd> tx nft issue <denom-id> --from=<key-name> --name=<denom-name> --schema=<schema-content or path to schema.json> --chain-id=<chain-id> --fees=<fee>
```

## Mint token
```
$ <appd> tx nft mint <denom-id> <token-id> --uri=<uri> --recipient=<recipient> --from=<key-name> --chain-id=<chain-id> --fees=<fee>
```

## Query Collection
```
$ <appd> query nft collection <denom-id>
```

## Approve NFT
```
$ <appd> tx nft approve <approvedAddress> <denom-id> <token-id> --from=<key-name> --chain-id=<chain-id> --fees=<fee>
```

## Revoke Approval
```
$ <appd> tx nft revoke <addressToRevoke> <denom-id> <token-id>--uri=<uri> --from=<key-name> --chain-id=<chain-id> --fees=<fee>
```

## Approve All
```
$ <appd> tx nft approveAll <operator> <true/false> --from=<key-name> --chain-id=<chain-id> --fees=<fee>
```