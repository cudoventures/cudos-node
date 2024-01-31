#!/bin/bash
BINARY=build/cudos-noded

CHAIN_ID=cudos-node
VAL_KEY=test0
VOTER=
VAL=$($BINARY keys show -a $VAL_KEY --keyring-backend test)
CONTRACT=cw20_base
PROPOSAL=1
DENOM=${2:-"acudos"}
HOME=cudos-data

echo "submit wasm store proposal..."
$BINARY tx wasm submit-proposal wasm-store scripts/wasm/$CONTRACT.wasm --title "Add $CONTRACT" \
  --summary "Let's upload this contract 3" \
  --from $VAL_KEY --keyring-backend test --chain-id $CHAIN_ID -y \
    --gas auto --gas-adjustment 1.3 > /dev/null


echo "deposit to proposal..."
sleep 15
# $BINARY query gov proposal $PROPOSAL
$BINARY tx gov deposit $PROPOSAL 40000000000000000000$DENOM --from $VAL_KEY --keyring-backend test \
    --chain-id $CHAIN_ID -y -b sync --gas auto --gas-adjustment 1.3 --home $HOME > /dev/null

echo "process to vote..."
sleep 5
$BINARY tx gov vote $PROPOSAL yes --from $VAL_KEY --keyring-backend test \
    --chain-id $CHAIN_ID -y -b sync --gas auto  --gas-adjustment 1.3  --home $HOME > /dev/null

echo "Waiting for voting periods to finish..."
COUNTER=0
while ((COUNTER < 12)); do
    # Capture output of $BINARY command and extract "status" using jq
    status=$($BINARY q gov proposal 1 --output=json | jq '.status')
    sleep 5
    echo "Current proposal status: $status"
    # Increment COUNTER using arithmetic expansion
    ((COUNTER++))
    if [ "$status" == "PROPOSAL_STATUS_PASSED" ]; then
        break
    fi
done


$BINARY query wasm list-code