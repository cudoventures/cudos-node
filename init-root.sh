rm -R ./data/*
rm -R ./data/.*

MONIKER="cudos-root-node"
CHAIN_ID="cudos-network"

UNBONDING_TIME="1814401s" #1814400s originally
JAIL_DURATION="605s" #600s originally
TIMEOUT_COMMIT="10s" #5s originally 

GOV_PROPOSAL_DEPOSIT_DENOM="stake" # stake orginally
GOV_PROPOSAL_DEPOSIT_AMOUNT="100" # 10000000 originally   

MINT_DENOM="stake" # stake originally
STAKE_DENOM="stake" # stake originally

cudos-noded init $MONIKER --chain-id=$CHAIN_ID
sed -i "104s/enable = false/enable = true/" /usr/cudos/data/.cudos-node/config/app.toml
sed -i "s/laddr = \"tcp:\/\/127.0.0.1:26657\"/laddr = \"tcp:\/\/0.0.0.0:26657\"/" /usr/cudos/data/.cudos-node/config/config.toml

# setting time after commit before proposing a new block
sed -i "s/timeout_commit = \"5s\"/timeout_commit = \"$TIMEOUT_COMMIT\"/" /usr/cudos/data/.cudos-node/config/config.toml

# setting unbonding time
cat /usr/cudos/data/.cudos-node/config/genesis.json | jq --arg UNBONDING_TIME "$UNBONDING_TIME" '.app_state.staking.params.unbonding_time = $UNBONDING_TIME' > /usr/cudos/data/.cudos-node/config/tmp_genesis.json && mv /usr/cudos/data/.cudos-node/config/tmp_genesis.json /usr/cudos/data/.cudos-node/config/genesis.json

# setting jailtime duration
cat /usr/cudos/data/.cudos-node/config/genesis.json | jq --arg JAIL_DURATION "$JAIL_DURATION" '.app_state.slashing.params.downtime_jail_duration = $JAIL_DURATION' > /usr/cudos/data/.cudos-node/config/tmp_genesis.json && mv /usr/cudos/data/.cudos-node/config/tmp_genesis.json /usr/cudos/data/.cudos-node/config/genesis.json

# setting government proposal denom and minimal deposit
cat /usr/cudos/data/.cudos-node/config/genesis.json | jq --arg GOV_PROPOSAL_DEPOSIT_DENOM "$GOV_PROPOSAL_DEPOSIT_DENOM" '.app_state.gov.deposit_params.min_deposit[0].denom = $GOV_PROPOSAL_DEPOSIT_DENOM' > /usr/cudos/data/.cudos-node/config/tmp_genesis.json && mv /usr/cudos/data/.cudos-node/config/tmp_genesis.json /usr/cudos/data/.cudos-node/config/genesis.json
cat /usr/cudos/data/.cudos-node/config/genesis.json | jq --arg GOV_PROPOSAL_DEPOSIT_AMOUNT "$GOV_PROPOSAL_DEPOSIT_AMOUNT" '.app_state.gov.deposit_params.min_deposit[0].amount = $GOV_PROPOSAL_DEPOSIT_AMOUNT' > /usr/cudos/data/.cudos-node/config/tmp_genesis.json && mv /usr/cudos/data/.cudos-node/config/tmp_genesis.json /usr/cudos/data/.cudos-node/config/genesis.json

# setting mint coin denom
cat /usr/cudos/data/.cudos-node/config/genesis.json | jq --arg MINT_DENOM "$MINT_DENOM" '.app_state.mint.params.mint_denom = $MINT_DENOM' > /usr/cudos/data/.cudos-node/config/tmp_genesis.json && mv /usr/cudos/data/.cudos-node/config/tmp_genesis.json /usr/cudos/data/.cudos-node/config/genesis.json

# setting stake coin denom
cat /usr/cudos/data/.cudos-node/config/genesis.json | jq --arg STAKE_DENOM "$STAKE_DENOM" '.app_state.staking.params.bond_denom = $STAKE_DENOM' > /usr/cudos/data/.cudos-node/config/tmp_genesis.json && mv /usr/cudos/data/.cudos-node/config/tmp_genesis.json /usr/cudos/data/.cudos-node/config/genesis.json


cudos-noded keys add validator01 --keyring-backend test

VALIDATOR_ADDRESS=$(cudos-noded keys show validator01 -a)

cudos-noded add-genesis-account $VALIDATOR_ADDRESS 100000000000stake,10000000stake2

cudos-noded gentx validator01 10000000stake --chain-id $CHAIN_ID --keyring-backend test

cudos-noded collect-gentxs



