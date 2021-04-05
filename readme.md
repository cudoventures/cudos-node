# POC 1 - Create Cosmos Chain with BFT, Wallets and preconfigured accounts
This project uses the default blog project from starport. It's purpose is to prove POC 1 requirements, such as:
 - Ability to launch a "vanilla" Cosmos blockchain and create a closed public test environment. 
 - Ability of the network to support Tendermint BFT with Ethereum type of wallets. 
 - The network should start with ~20 accounts/validators that have pre-configured vested balance in the genesis block.

## Links:
 - https://docs.cosmos.network/master/modules/auth/
 - https://docs.cosmos.network/master/basics/accounts.html
 - https://docs.cosmos.network/master/modules/auth/05_vesting.html


## Ability to launch a "vanilla" Cosmos blockchain and create a closed public test environment.
Launching a vanilla Cosmos blockchain is possible using Cosmos tools like Tendermint or Starport.
  
## Ability of the network to support Tendermint BFT with Ethereum type of wallets
The network supports Tendermint BFT by default. The wallets private keys are generated using secp256k1 also by default.

## Creating accounts/validators with preconfigured vested balance in the genesis block
Accounts can be created in the config.yml by specifying a name and the coins held. This would automatically add them to the genesis block upon launching the network. Another way to do this is by manually writing the genesis json, from which the genesis block is built.



# Local build of this project
--home param on each commad indicates the blockchain storage directory.

## make
make

## init
cudos-poc-01d init cudos-poc-01-network --chain-id=cudos-poc-01-network --home=./HOME

## create staking account
cudos-poc-01d keys add validator01 --keyring-backend test --home=./HOME

## get validator's address
cudos-poc-01d keys show validator01 -a --keyring-backend test --home=./HOME

## add stacking account
cudos-poc-01d add-genesis-account $MY_VALIDATOR_ADDRESS 100000000000stake --home=./HOME

## create gen tx
cudos-poc-01d gentx validator01 100000000stake --chain-id cudos-poc-01-network --keyring-backend test --home=./HOME

## add tx to genesis
cudos-poc-01d collect-gentxs --home=./HOME

## start
cudos-poc-01d start --home=./HOME

# docker
1. Build persistent-node
cd ./docker
docker-compose -f .\persistent-node.yml -p cudos-network-persistent-node up --build

2. After node starts copy its it and paste it into full-node.yml
Peer node looks like:
P2P Node ID ID=de14a2005d220171c7133efb31b3f3e1d7ba776a file=/root/.blog/config/node_key.json module=p2p

3. Run full-node
cd ./docker
docker-compose -f .\full-node.yml -p cudos-network-full-node up --build

# Add local account
cudos-poc-01d keys add user01 --keyring-backend test --home=./HOME

# Send currency
cudos-poc-01d tx bank send $MY_VALIDATOR_ADDRESS $RECIPIENT 10stake --chain-id=cudos-poc-01-network --keyring-backend test --home=./HOME

# Check that the recipient account did receive the tokens.
cudos-poc-01d query bank balances $RECIPIENT --chain-id=cudos-poc-01-network --home=./HOME

# Add vesting account (only during genesis)
cudos-poc-01d keys add user-vesting --keyring-backend test --home=./HOME
cudos-poc-01d add-genesis-account $VESTING_ACCOUNT 1000stake --vesting-amount 500stake --vesting-end-time 1617613800 --home=./HOME