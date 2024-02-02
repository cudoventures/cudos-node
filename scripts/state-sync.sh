#!/bin/bash
VERSION=v1.1.0
SED_BINARY=sed
HOME_DIR=cudos-data
ROOT=$(pwd)
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

# install binary if not exist
if [ ! -f "_build/$VERSION.zip" ] &> /dev/null
then
    mkdir -p _build/${VERSION}
    wget -c "https://github.com/CudoVentures/cudos-node/archive/refs/tags/${VERSION}.zip" -O _build/${VERSION}.zip
    unzip _build/${VERSION}.zip -d _build/${VERSION}
fi
# reinstall old binary
if [ $# -eq 1 ] && [ $1 == "--reinstall-old" ] || ! command -v _build/$VERSION/cudos-noded &> /dev/null; then
    cd ./_build/${VERSION}/cudos-node-${VERSION:1}
    GOBIN="$ROOT/_build/${VERSION}" go install -mod=readonly ./...
    cd ../..
fi
BINARY=_build/${VERSION}/cudos-noded

# clean the data
if [ -d "cudos-data" ]; then
    rm -rf cudos-data
fi
sleep 1
# init the node
$BINARY init test --chain-id cudos-1

#######################################################
###         State Sync Configuration Options        ###
#######################################################
RPC_SERVER=https://rpc.cudos.org:443
INTERVAL=1000
LATEST_HEIGHT=$(curl -s $RPC_SERVER/block | jq -r .result.block.header.height)
BLOCK_HEIGHT=$(($LATEST_HEIGHT - $INTERVAL))
TRUST_HASH=$(curl -s $RPC_SERVER/block?height=$BLOCK_HEIGHT | jq -r .result.block_id.hash)

# For now, hardcode the persistent_peer
PERSISTENT_PEER="2cc0a12ff1038509b2ed64719fcddfdded9a04ad@198.244.179.112:26657"


# Perform to modify the configuration file to enable state sync

$SED_BINARY -i "s/\(enable *= *\).*/\1true/" $HOME_DIR/config/config.toml
$SED_BINARY -i "s|rpc_servers *=.*|rpc_servers = \"$RPC_SERVER,$RPC_SERVER\"|" $HOME_DIR/config/config.toml
$SED_BINARY -i "s/\(trust_height *= *\).*/\1$BLOCK_HEIGHT/"  $HOME_DIR/config/config.toml
$SED_BINARY -i "s/\(trust_hash *= *\).*/\1\"$TRUST_HASH\"/"  $HOME_DIR/config/config.toml

# Start the node
$BINARY start --home $HOME_DIR --p2p.persistent_peers=$PERSISTENT_PEER --log_level debug


