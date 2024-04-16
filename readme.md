# Cudos Testnet

## Pre-requirements
Starport: https://github.com/tendermint/starport/releases (linux only) <br />
or <br />
make (different ways to install it depending on OS) 

## Requirements
This project uses the default blog project from starport. It's purpose is to prove POC 1 requirements, such as:
 - Ability to launch a "vanilla" Cosmos blockchain and create a closed public test environment. 
 - Ability of the network to support Tendermint BFT with Ethereum type of wallets. 
 - The network should start with ~20 accounts/validators that have pre-configured vested balance in the genesis block.

## Useful links:
 - https://docs.cosmos.network/master/modules/auth/
 - https://docs.cosmos.network/master/basics/accounts.html


## Ability to launch a "vanilla" Cosmos blockchain and create a closed public test environment.
Launching a vanilla Cosmos blockchain is possible by starting this project using methods below.
  
## Ability of the network to support Tendermint BFT with Ethereum type of wallets

The network supports Tendermint BFT by default. The wallets private keys are generated using secp256k1 also by default. The retionale of using Ethereum type of wallets is to ensure that the users will be able to import their ethereum wallets into cudos blockchain using theirs seed phase. After the import the users will expect to see that their balance from ethereum blockchain is transferred to cudos. Although the cryptography is the same it is used in a slightly different manner so a converted is developed. Its usage is described below. It can convert ethereum public key to cudos wallet address. Using cudos wallet address a wallet can be pre-funded with required tokes so when a user import his wallet, using his seed, the balance will be correct.

## Creating accounts/validators with preconfigured vested balance in the genesis block
There are three ways to add preconfigured accounts with/without vested balance in the genesis block.
1. Modifying genesis.json after the blockchain is initialized, but before it is first started. This is not recommended method.
2. Using config.yml. This method works if the blockchain is initialized with starport utility.
3. Using commands from the binary itself. This method works if the blockchain is manually initialized without startport utility.

## Manual build
You can build and run cudos-node on your local machine. 
To do it please follow guide [here](https://www.notion.so/cudo/How-to-Start-a-one-node-Cudos-network-locally-no-Docker-4fc5a7ea9f054a2daebb196c5c14a84c)

# Starport build
Configure accounts and validators in config.yml after that just start the blockchain

    starport serve

# Docker run
You can run empty blockchain for test purposes.
Prerequisites: you have installed docker and docker-compose. 

Staying in root fo directory, run 
```bash
  make docker-run
```

It will run docker-container with single cudos-node with preset state.
It will include validator and 4 test accounts with 2000000000000000000000000000 acudos balance. 
The node **will not be synced** with cudos-testnet. 
You can interact with it using standard cudos-node port 26657. 

To stop container run
```bash
  make docker-stop
```

# Build proto-files
```bash
  make proto-gen
```


# Converting ethereum public keys to cosmos wallet address
Run the converter and pass a ethereum public key as argument.

    go run ./converter 0x03139bb3b92e99d034ee38674a0e29c4aad83dd09b3fa465a265da310f9948fbe6

<b>Example ethereum mnemonic:</b> battle erosion opinion city birth modify scale hood caught menu risk rather<br >
<b>Example ethereum public key (32 bytes, compressed form):</b> 0x03139bb3b92e99d034ee38674a0e29c4aad83dd09b3fa465a265da310f9948fbe6

This mnemonic could be imported into cudos blockchain in order to verify that resulting account access will be the same as generated from the converter.

    cudos-noded keys add ruser02 --recover --hd-path="m/44'/60'/0'/0/0"

# Useful commands
## Send currency
    cudos-noded tx bank send $VALIDATOR_ADDRESS $RECIPIENT 51000000stake --chain-id=cudos-network --keyring-backend test

## Check balances
    cudos-noded query bank balances $RECIPIENT --chain-id=cudos-network

## Create validator
    cudos-noded tx staking create-validator --amount=2000000000000000000000000acudos \
    --from=val-2 \
    --pubkey=$(cudos-noded tendermint show-validator) \
    --moniker=cudos-node-02 \
    --chain-id=cudos-local-network \
    --commission-rate="0.10" \
    --commission-max-rate="0.20" \
    --commission-max-change-rate="0.01" \
    --min-self-delegation="50000000000000000000000" \
    --gas="auto" \
    --gas-prices="5000000000000acudos" \
    --gas-adjustment="1.80" \
    --keyring-backend test

# Reset the blockchain
All data of the blockchain is store at ~/.blog folder. By deleting it the entire blockchain is completely reset and it must be initialized again.

# Static compile
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-extldflags "-static"' ./cmd/cudos-noded/

export CGO_LDFLAGS="-lpthread -ldl"
go build -v -a -tags netgo,osusergo -ldflags='-lpthread -extldflags "-lpthread -static"' ./cmd/cudos-noded/
