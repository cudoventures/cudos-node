#!bin/bash

BINARY=build/cudos-noded

CHAIN_ID=cudos-app

CONTRACT=cw20_base
PROPOSAL=1
HOME=cudos-data
DENOM="acudos"

$BINARY tx wasm store scripts/wasm/$CONTRACT.wasm \
--from test0 --keyring-backend test --chain-id $CHAIN_ID -y -b block \
  --gas 9000000 --gas-prices 0.025stake --home $HOME


$BINARY query gov proposal $PROPOSAL

$BINARY tx gov deposit 1 "40000000000000000000000000${DENOM}" --from test1 --keyring-backend test --chain-id $CHAIN_ID --home $HOME -y

$BINARY tx gov vote 1 yes --from test0 --keyring-backend test --chain-id $CHAIN_ID --home $HOME -y
