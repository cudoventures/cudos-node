rm -R ./cudos-data
# rm -R ./data/.*

cudos-noded init $MONIKER

if [ -z "$PERSISTENT_NODE_ID" ]
then
    PERSISTENT_NODE_ID=$(cudos-noded tendermint show-node-id --home=/usr/cudos-root-data)@cudos-root-node:26656
fi

cp /usr/cudos-root-data/config/genesis.json ./cudos-data/config/genesis.json

sed -i "s/persistent_peers = \"\"/persistent_peers = \"$PERSISTENT_NODE_ID\"/g" ./cudos-data/config/config.toml