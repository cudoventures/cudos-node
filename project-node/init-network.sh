rm -R ./data/*
rm -R ./data/.*

cudos-noded init $MONIKER

if [ -z "$PERSISTENT_NODE_ID" ]
then
    PERSISTENT_NODE_ID=$(cudos-noded tendermint show-node-id --home=/usr/cudos-root/.cudos-node)@cudos-root-node:26656
fi

cp /usr/cudos-root/.cudos-node/config/genesis.json /usr/cudos/data/.cudos-node/config/genesis.json

sed -i "s/persistent_peers = \"\"/persistent_peers = \"$PERSISTENT_NODE_ID\"/g" /usr/cudos/data/.cudos-node/config/config.toml