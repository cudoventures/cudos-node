# project
This project uses the default blog project from starport


# local build

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

