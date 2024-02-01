#!/bin/bash

BINARY=$1
CONTINUE=${CONTINUE:-"false"}
HOME_DIR=cudos-data
ENV=${ENV:-""}

if [ "$CONTINUE" == "true" ]; then
    $BINARY start --home $HOME_DIR --log_level debug
    exit 0
fi

rm -rf cudos-data
pkill cudos-noded

# check DENOM is set. If not, set to acudos
DENOM=${2:-"acudos"}

COMMISSION_RATE=0.01
COMMISSION_MAX_RATE=0.02
MIN_SELF_DELEGATION=50000000000000000000000
SED_BINARY=sed
# check if this is OS X
if [[ "$OSTYPE" == "darwin"* ]]; then
    # check if gsed is installed
    if ! command -v gsed &> /dev/null
    then
        echo "gsed could not be found. Please install it with 'brew install gnu-sed'"
        exit
    else
        SED_BINARY=gsed
    fi
fi

# check BINARY is set. If not, build cudos-noded and set BINARY
if [ -z "$BINARY" ]; then
    make build
    BINARY=build/cudos-noded
fi

CHAIN_ID="cudos-node"
KEYRING="test"
KEY="test0"
KEY1="test1"
KEY2="test2"

# Function updates the config based on a jq argument as a string
update_test_genesis () {
    # update_test_genesis '.consensus_params["block"]["max_gas"]="100000000"'
    cat $HOME_DIR/config/genesis.json | jq "$1" > $HOME_DIR/config/tmp_genesis.json && mv $HOME_DIR/config/tmp_genesis.json $HOME_DIR/config/genesis.json
}

$BINARY init --chain-id $CHAIN_ID moniker --home $HOME_DIR

$BINARY keys add $KEY --keyring-backend $KEYRING --home $HOME_DIR
$BINARY keys add $KEY1 --keyring-backend $KEYRING --home $HOME_DIR
$BINARY keys add $KEY2 --keyring-backend $KEYRING --home $HOME_DIR

# Get the generated addresses
TEST0_ADDRESS=$($BINARY keys show $KEY -a --keyring-backend $KEYRING --home $HOME_DIR)
TEST1_ADDRESS=$($BINARY keys show $KEY1 -a --keyring-backend $KEYRING --home $HOME_DIR)
TEST2_ADDRESS=$($BINARY keys show $KEY2 -a --keyring-backend $KEYRING --home $HOME_DIR)

# Allocate genesis accounts (cosmos formatted addresses)
$BINARY genesis add-genesis-account $TEST0_ADDRESS "500000000000000000000000000${DENOM}" --home $HOME_DIR
$BINARY genesis add-genesis-account $TEST1_ADDRESS "500000000000000000000000000${DENOM}" --home $HOME_DIR
$BINARY genesis add-genesis-account $TEST2_ADDRESS "500000000000000000000000000${DENOM}" --home $HOME_DIR 

update_test_genesis '.app_state["gov"]["params"]["voting_period"]="50s"'
update_test_genesis '.app_state["mint"]["params"]["mint_denom"]="'$DENOM'"'
update_test_genesis '.app_state["gov"]["params"]["min_deposit"]=[{"denom":"'$DENOM'","amount": "300"}]'
update_test_genesis '.app_state["gov"]["params"]["min_initial_deposit_ratio"]="0.0"'

update_test_genesis '.app_state["crisis"]["constant_fee"]={"denom":"'$DENOM'","amount":"1000"}'
update_test_genesis '.app_state["staking"]["params"]["bond_denom"]="'$DENOM'"'

# enable rest server and swagger
$SED_BINARY -i '0,/enable = false/s//enable = true/' $HOME_DIR/config/app.toml
$SED_BINARY -i 's/swagger = false/swagger = true/' $HOME_DIR/config/app.toml

# Sign genesis transaction
# TEST0 is the validator
$BINARY genesis gentx $KEY "500000000000000000000000${DENOM}" --commission-rate=$COMMISSION_RATE --min-self-delegation=$MIN_SELF_DELEGATION --commission-max-rate=$COMMISSION_MAX_RATE --keyring-backend $KEYRING --chain-id $CHAIN_ID --home $HOME_DIR

# Collect genesis tx
echo "Collecting genesis tx..."
$BINARY genesis collect-gentxs --home $HOME_DIR

# Run this to ensure everything worked and that the genesis file is setup correctly
# This raises an error since Cudos has an additional genesis Tx : MsgSetOrchestratorAddress
$BINARY genesis validate-genesis --home $HOME_DIR

$BINARY start --home $HOME_DIR
