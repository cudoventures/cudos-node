# POC 1 - Create Cosmos Chain with BFT, Wallets and preconfigured accounts

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
Build the blockchian binary into $GOPATH directory using "cudos-noded" name. 

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
docker-compose -f .\persistent-node.yml -p cudos-network-persistent-node up --build

2. After node starts copy its it and paste it into full-node.yml
Peer node looks like:
P2P Node ID ID=de14a2005d220171c7133efb31b3f3e1d7ba776a file=/root/.blog/config/node_key.json module=p2p

3. Run full-node
cd ./docker
docker-compose -f .\full-node.yml -p cudos-network-full-node up --build

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
    cudos-noded tx bank send $MY_VALIDATOR_ADDRESS $RECIPIENT 10stake --chain-id=cudos-node --keyring-backend test

## Check balances
    cudos-noded query bank balances $RECIPIENT --chain-id=cudos-node

<br />
<br />
<br />

# Reset the blockchain
All data of the blockchain is store at ~/.blog folder. By deleting it the entire blockchain is completely reset and it must be initialized again.