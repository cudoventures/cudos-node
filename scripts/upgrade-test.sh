#!/bin/bash

# the upgrade is a fork, "true" otherwise
FORK=${FORK:-"false"}

OLD_VERSION=v1.2.0
UPGRADE_WAIT=${UPGRADE_WAIT:-20}
HOME=cudos-data
ROOT=$(pwd)
DENOM=acudos
CHAIN_ID=cudos-node
SOFTWARE_UPGRADE_NAME="v1.3"
ADDITIONAL_PRE_SCRIPTS=${ADDITIONAL_PRE_SCRIPTS:-""}
ADDITIONAL_AFTER_SCRIPTS=${ADDITIONAL_AFTER_SCRIPTS:-""}
NEW_VERSION=v1.3.0
if [[ "$FORK" == "true" ]]; then
    export cudos-noded_HALT_HEIGHT=20
fi

# underscore so that go tool will not take gocache into account
mkdir -p _build/gocache
export GOMODCACHE=$ROOT/_build/gocache

# install old binary if not exist
if [ ! -f "_build/$OLD_VERSION.zip" ] &> /dev/null
then
    mkdir -p _build/old
    wget -c "https://github.com/CudoVentures/cudos-node/archive/refs/tags/${OLD_VERSION}.zip" -O _build/${OLD_VERSION}.zip
    unzip _build/${OLD_VERSION}.zip -d _build
fi

# reinstall old binary
if [ $# -eq 1 ] && [ $1 == "--reinstall-old" ] || ! command -v _build/old/cudos-noded &> /dev/null; then
    cd ./_build/cudos-node-${OLD_VERSION:1}
    GOBIN="$ROOT/_build/old" go install -mod=readonly ./...
    cd ../..
fi

# install new binary
if ! command -v _build/new/cudos-noded &> /dev/null
then
    echo "Installing new binary"
    mkdir -p ./_build/cudos-node-${NEW_VERSION:1} && cd ./_build/cudos-node-${NEW_VERSION:1}
    GOBIN="$ROOT/_build/new" go install -mod=readonly ../../...
    cd ../..
fi

# run old node
if [[ "$OSTYPE" == "darwin"* ]]; then
    screen -L -dmS node1 bash scripts/run-node.sh _build/old/cudos-noded $DENOM
else
    screen -L -Logfile $HOME/log-screen.txt -dmS node1 bash scripts/run-node.sh _build/old/cudos-noded $DENOM
fi

sleep 20

# execute additional pre scripts
if [ ! -z "$ADDITIONAL_PRE_SCRIPTS" ]; then
    # slice ADDITIONAL_SCRIPTS by ,
    SCRIPTS=($(echo "$ADDITIONAL_PRE_SCRIPTS" | tr ',' ' '))
    for SCRIPT in "${SCRIPTS[@]}"; do
         # check if SCRIPT is a file
        if [ -f "$SCRIPT" ]; then
            echo "executing additional pre scripts from $SCRIPT"
            source $SCRIPT
            sleep 5
        else
            echo "$SCRIPT is not a file"
        fi
    done
fi

run_fork () {
    echo "forking"

    while true; do 
        BLOCK_HEIGHT=$(./_build/old/cudos-noded status 2> >(jq '.SyncInfo.latest_block_height' -r))
        # if BLOCK_HEIGHT is not empty
        if [ ! -z "$BLOCK_HEIGHT" ]; then
            echo "BLOCK_HEIGHT = $BLOCK_HEIGHT"
            sleep 10
        else
            echo "BLOCK_HEIGHT is empty, forking"
            break
        fi
    done
}

run_upgrade () {
    echo "upgrading"

    STATUS_INFO=$(./_build/old/cudos-noded status --home $HOME 2>&1 | jq -r '.SyncInfo.latest_block_height')
    UPGRADE_HEIGHT=$((STATUS_INFO + 20))
    echo "creating tar file for new binary"
    tar -cf ./_build/new/cudos-noded.tar -C ./_build/new cudos-noded
    SUM=$(shasum -a 256 ./_build/new/cudos-noded.tar | cut -d ' ' -f1)
    UPGRADE_INFO=$(jq -n '
    {
        "binaries": {
            "linux/amd64": "file://'$(pwd)'/_build/new/cudos-noded.tar?checksum=sha256:'"$SUM"'",
        }
    }')

    ./_build/old/cudos-noded keys list --home $HOME --keyring-backend test
    echo "submitting upgrade proposal"
    ./_build/old/cudos-noded tx gov submit-proposal software-upgrade "$SOFTWARE_UPGRADE_NAME" --upgrade-height $UPGRADE_HEIGHT --upgrade-info "$UPGRADE_INFO" --title "upgrade" --description "upgrade"  --from test1 --keyring-backend test --chain-id $CHAIN_ID --home $HOME  -y

    sleep 5

    ./_build/old/cudos-noded tx gov deposit 1 "40000000000000000000000000${DENOM}" --from test1 --keyring-backend test --chain-id $CHAIN_ID --home $HOME -y

    sleep 5

    ./_build/old/cudos-noded tx gov vote 1 yes --from test0 --keyring-backend test --chain-id $CHAIN_ID --home $HOME -y

    sleep 5

    ./_build/old/cudos-noded tx gov vote 1 yes --from test1 --keyring-backend test --chain-id $CHAIN_ID --home $HOME -y

    sleep 5

    # determine block_height to halt
    while true; do 
        BLOCK_HEIGHT=$(./_build/old/cudos-noded status --home $HOME 2>&1 | jq -r '.SyncInfo.latest_block_height')
        if [ $BLOCK_HEIGHT = "$UPGRADE_HEIGHT" ]; then
            # assuming running only 1 cudos-noded
            echo "BLOCK HEIGHT = $UPGRADE_HEIGHT REACHED, KILLING OLD ONE"
            pkill cudos-noded
            break
        else
            ./_build/old/cudos-noded q gov proposal 1 --output=json | jq ".status"
            echo "BLOCK_HEIGHT = $BLOCK_HEIGHT"
            sleep 10
        fi
    done
}

# if FORK = true
if [[ "$FORK" == "true" ]]; then
    run_fork
    unset cudos-noded_HALT_HEIGHT
else
    run_upgrade
fi

sleep 5

# run new node
if [[ "$OSTYPE" == "darwin"* ]]; then
    CONTINUE="true" screen -L -dmS node1 bash scripts/run-node.sh _build/new/cudos-noded $DENOM
else
    CONTINUE="true" screen -L -Logfile $HOME/log-screen.txt -dmS node1 bash scripts/run-node.sh _build/new/cudos-noded $DENOM
fi

sleep 20

# execute additional after scripts
if [ ! -z "$ADDITIONAL_AFTER_SCRIPTS" ]; then
    # slice ADDITIONAL_SCRIPTS by ,
    SCRIPTS=($(echo "$ADDITIONAL_AFTER_SCRIPTS" | tr ',' ' '))
    for SCRIPT in "${SCRIPTS[@]}"; do
         # check if SCRIPT is a file
        if [ -f "$SCRIPT" ]; then
            echo "executing additional after scripts from $SCRIPT"
            source $SCRIPT
            sleep 5
        else
            echo "$SCRIPT is not a file"
        fi
    done
fi