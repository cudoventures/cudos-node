if [[ -z "${CUDOS_HOME}" ]]; then
    CUDOS_HOME="./cudos-data"
fi

rm -R $CUDOS_HOME
# rm -R ./data/.*

# chain parameters
MONIKER="cudos-root-node"
CHAIN_ID="cudos-network"
TIMEOUT_COMMIT="5s" #5s originally 

MIN_SELF_DELEGATION="1" # minimum tokens sto stake multiplyer by 1 000 000 for validator01

# slashing parameters
JAIL_DURATION="600s" #600s originally

# staking parameters
BOND_DENOM="ucudos" # stake originally
UNBONDING_TIME="1814400s" #1814400s originally
MAX_VALIDATORS="100" # 100 originally

# government parameters
GOV_PROPOSAL_MIN_DEPOSIT_DENOM="ucudos" # stake orginally
GOV_PROPOSAL_MIN_DEPOSIT_AMOUNT="10000000" # 10000000 originally   
GOV_PROPOSAL_MAX_DEPOSIT_PERIOD="172800s" # 172800s originally
GOV_PROPOSAL_VOTING_PERIOD="172800s" # 172800s originally
GOV_QUORUM="0.334000000000000000" # 0.334000000000000000 originally
GOV_THRESHOLD="0.500000000000000000" # 0.500000000000000000 originally
GOV_VETO_THRESHOLD="0.334000000000000000" # 0.334000000000000000 originally

# mint parameters
MINT_DENOM="ucudos" # stake originally
MINT_INFLATION="0.0000000013" # 0.130000000000000000 originally
MINT_INFLATION_RATE_CHANGE="0.0000000013" # 0.130000000000000000 originally
MINT_INFLATION_MAX="0.0000000013" # 0.200000000000000000 originally
MINT_INFLATION_MIN="0.0000000013" # 0.070000000000000000 originally
MINT_GOAL_BONDED="0.670000000000000000" # 0.670000000000000000 originally
BLOCKS_PER_YEAR="6311520" # 6311520 originally


DENOM_METADATA_DESC="The native staking token of the Cudos Hub." 
DENOM1="ucudos" EXP1="0" ALIAS1="microcudos"
DENOM2="mcudos" EXP2="3" ALIAS2="millicudos"
DENOM3="cudos" EXP3="6"
BASE="ucudos"
DISPLAY="cudos"

cudos-noded init $MONIKER --chain-id=$CHAIN_ID
sed -i "104s/enable = false/enable = true/" "${CUDOS_HOME}/config/app.toml"
sed -i "s/laddr = \"tcp:\/\/127.0.0.1:26657\"/laddr = \"tcp:\/\/0.0.0.0:26657\"/" "${CUDOS_HOME}/config/config.toml"

# enable cors origin for local testing
sed -i "s/enabled-unsafe-cors = false/enabled-unsafe-cors = true/" ${CUDOS_HOME}/config/app.toml
sed -i "s/cors_allowed_origins = \[\]/cors_allowed_origins = \[\"\*\"\]/" ${CUDOS_HOME}/config/config.toml

# setting time after commit before proposing a new block
sed -i "s/timeout_commit = \"5s\"/timeout_commit = \"$TIMEOUT_COMMIT\"/" "${CUDOS_HOME}/config/config.toml"

# setting slashing time
cat "${CUDOS_HOME}/config/genesis.json" | jq --arg JAIL_DURATION "$JAIL_DURATION" '.app_state.slashing.params.downtime_jail_duration = $JAIL_DURATION' > "${CUDOS_HOME}/config/tmp_genesis.json" && mv "${CUDOS_HOME}/config/tmp_genesis.json" "${CUDOS_HOME}/config/genesis.json"

# setting staking params
cat "${CUDOS_HOME}/config/genesis.json" | jq --arg UNBONDING_TIME "$UNBONDING_TIME" '.app_state.staking.params.unbonding_time = $UNBONDING_TIME' > "${CUDOS_HOME}/config/tmp_genesis.json" && mv "${CUDOS_HOME}/config/tmp_genesis.json" "${CUDOS_HOME}/config/genesis.json"
cat "${CUDOS_HOME}/config/genesis.json" | jq --arg BOND_DENOM "$BOND_DENOM" '.app_state.staking.params.bond_denom = $BOND_DENOM' > "${CUDOS_HOME}/config/tmp_genesis.json" && mv "${CUDOS_HOME}/config/tmp_genesis.json" "${CUDOS_HOME}/config/genesis.json"
cat "${CUDOS_HOME}/config/genesis.json" | jq --arg MAX_VALIDATORS "$MAX_VALIDATORS" '.app_state.staking.params.max_validators = $MAX_VALIDATORS' > "${CUDOS_HOME}/config/tmp_genesis.json" && mv "${CUDOS_HOME}/config/tmp_genesis.json" "${CUDOS_HOME}/config/genesis.json"
cat "${CUDOS_HOME}/config/genesis.json" | jq --arg BOND_DENOM "$BOND_DENOM" '.app_state.crisis.constant_fee.denom = $BOND_DENOM' > "${CUDOS_HOME}/config/tmp_genesis.json" && mv "${CUDOS_HOME}/config/tmp_genesis.json" "${CUDOS_HOME}/config/genesis.json"

# setting government proposal params
cat "${CUDOS_HOME}/config/genesis.json" | jq --arg GOV_PROPOSAL_MIN_DEPOSIT_DENOM "$GOV_PROPOSAL_MIN_DEPOSIT_DENOM" '.app_state.gov.deposit_params.min_deposit[0].denom = $GOV_PROPOSAL_MIN_DEPOSIT_DENOM' > "${CUDOS_HOME}/config/tmp_genesis.json" && mv "${CUDOS_HOME}/config/tmp_genesis.json" "${CUDOS_HOME}/config/genesis.json"
cat "${CUDOS_HOME}/config/genesis.json" | jq --arg GOV_PROPOSAL_MIN_DEPOSIT_AMOUNT "$GOV_PROPOSAL_MIN_DEPOSIT_AMOUNT" '.app_state.gov.deposit_params.min_deposit[0].amount = $GOV_PROPOSAL_MIN_DEPOSIT_AMOUNT' > "${CUDOS_HOME}/config/tmp_genesis.json" && mv "${CUDOS_HOME}/config/tmp_genesis.json" "${CUDOS_HOME}/config/genesis.json"
cat "${CUDOS_HOME}/config/genesis.json" | jq --arg GOV_PROPOSAL_MAX_DEPOSIT_PERIOD "$GOV_PROPOSAL_MAX_DEPOSIT_PERIOD" '.app_state.gov.deposit_params.max_deposit_period = $GOV_PROPOSAL_MAX_DEPOSIT_PERIOD' > "${CUDOS_HOME}/config/tmp_genesis.json" && mv "${CUDOS_HOME}/config/tmp_genesis.json" "${CUDOS_HOME}/config/genesis.json"
cat "${CUDOS_HOME}/config/genesis.json" | jq --arg GOV_PROPOSAL_VOTING_PERIOD "$GOV_PROPOSAL_VOTING_PERIOD" '.app_state.gov.voting_params.voting_period = $GOV_PROPOSAL_VOTING_PERIOD' > "${CUDOS_HOME}/config/tmp_genesis.json" && mv "${CUDOS_HOME}/config/tmp_genesis.json" "${CUDOS_HOME}/config/genesis.json"
cat "${CUDOS_HOME}/config/genesis.json" | jq --arg GOV_QUORUM "$GOV_QUORUM" '.app_state.gov.tally_params.quorum = $GOV_QUORUM' > "${CUDOS_HOME}/config/tmp_genesis.json" && mv "${CUDOS_HOME}/config/tmp_genesis.json" "${CUDOS_HOME}/config/genesis.json"
cat "${CUDOS_HOME}/config/genesis.json" | jq --arg GOV_THRESHOLD "$GOV_THRESHOLD" '.app_state.gov.tally_params.threshold = $GOV_THRESHOLD' > "${CUDOS_HOME}/config/tmp_genesis.json" && mv "${CUDOS_HOME}/config/tmp_genesis.json" "${CUDOS_HOME}/config/genesis.json"
cat "${CUDOS_HOME}/config/genesis.json" | jq --arg GOV_VETO_THRESHOLD "$GOV_VETO_THRESHOLD" '.app_state.gov.tally_params.veto_threshold = $GOV_VETO_THRESHOLD' > "${CUDOS_HOME}/config/tmp_genesis.json" && mv "${CUDOS_HOME}/config/tmp_genesis.json" "${CUDOS_HOME}/config/genesis.json"

# setting mint params
cat "${CUDOS_HOME}/config/genesis.json" | jq --arg MINT_DENOM "$MINT_DENOM" '.app_state.mint.params.mint_denom = $MINT_DENOM' > "${CUDOS_HOME}/config/tmp_genesis.json" && mv "${CUDOS_HOME}/config/tmp_genesis.json" "${CUDOS_HOME}/config/genesis.json"
cat "${CUDOS_HOME}/config/genesis.json" | jq --arg MINT_INFLATION "$MINT_INFLATION" '.app_state.mint.minter.inflation = $MINT_INFLATION' > "${CUDOS_HOME}/config/tmp_genesis.json" && mv "${CUDOS_HOME}/config/tmp_genesis.json" "${CUDOS_HOME}/config/genesis.json"
cat "${CUDOS_HOME}/config/genesis.json" | jq --arg MINT_INFLATION_RATE_CHANGE "$MINT_INFLATION_RATE_CHANGE" '.app_state.mint.params.inflation_rate_change = $MINT_INFLATION_RATE_CHANGE' > "${CUDOS_HOME}/config/tmp_genesis.json" && mv "${CUDOS_HOME}/config/tmp_genesis.json" "${CUDOS_HOME}/config/genesis.json"
cat "${CUDOS_HOME}/config/genesis.json" | jq --arg MINT_INFLATION_MAX "$MINT_INFLATION_MAX" '.app_state.mint.params.inflation_max = $MINT_INFLATION_MAX' > "${CUDOS_HOME}/config/tmp_genesis.json" && mv "${CUDOS_HOME}/config/tmp_genesis.json" "${CUDOS_HOME}/config/genesis.json"
cat "${CUDOS_HOME}/config/genesis.json" | jq --arg MINT_INFLATION_MIN "$MINT_INFLATION_MIN" '.app_state.mint.params.inflation_min = $MINT_INFLATION_MIN' > "${CUDOS_HOME}/config/tmp_genesis.json" && mv "${CUDOS_HOME}/config/tmp_genesis.json" "${CUDOS_HOME}/config/genesis.json"
cat "${CUDOS_HOME}/config/genesis.json" | jq --arg MINT_GOAL_BONDED "$MINT_GOAL_BONDED" '.app_state.mint.params.goal_bonded = $MINT_GOAL_BONDED' > "${CUDOS_HOME}/config/tmp_genesis.json" && mv "${CUDOS_HOME}/config/tmp_genesis.json" "${CUDOS_HOME}/config/genesis.json"
cat "${CUDOS_HOME}/config/genesis.json" | jq --arg BLOCKS_PER_YEAR "$BLOCKS_PER_YEAR" '.app_state.mint.params.blocks_per_year = $BLOCKS_PER_YEAR' > "${CUDOS_HOME}/config/tmp_genesis.json" && mv "${CUDOS_HOME}/config/tmp_genesis.json" "${CUDOS_HOME}/config/genesis.json"

# setting fractions metadata
cat "${CUDOS_HOME}/config/genesis.json" | jq --arg DENOM_METADATA_DESC "$DENOM_METADATA_DESC" --arg DENOM1 "$DENOM1" --arg EXP1 "$EXP1" --arg ALIAS1 "$ALIAS1" --arg DENOM2 "$DENOM2" --arg EXP2 "$EXP2" --arg ALIAS2 "$ALIAS2" --arg DENOM3 "$DENOM3" --arg EXP3 "$EXP3" --arg BASE "$BASE" --arg DISPLAY "$DISPLAY" '.app_state.bank.denom_metadata[0].description=$DENOM_METADATA_DESC | .app_state.bank.denom_metadata[0].denom_units[0].denom=$DENOM1 | .app_state.bank.denom_metadata[0].denom_units[0].exponent=$EXP1 | .app_state.bank.denom_metadata[0].denom_units[0].aliases[0]=$ALIAS1 | .app_state.bank.denom_metadata[0].denom_units[1].denom=$DENOM2 | .app_state.bank.denom_metadata[0].denom_units[1].exponent=$EXP2 | .app_state.bank.denom_metadata[0].denom_units[1].aliases[0]=$ALIAS2 | .app_state.bank.denom_metadata[0].denom_units[2].denom=$DENOM3 | .app_state.bank.denom_metadata[0].denom_units[2].exponent=$EXP3 | .app_state.bank.denom_metadata[0].base=$BASE | .app_state.bank.denom_metadata[0].display=$DISPLAY'  > "${CUDOS_HOME}/config/tmp_genesis.json" && mv "${CUDOS_HOME}/config/tmp_genesis.json" "${CUDOS_HOME}/config/genesis.json"

# add a new key entry from which to make validator
cudos-noded keys add root-validator --keyring-backend test
VALIDATOR_ADDRESS=$(cudos-noded keys show root-validator -a)

# # create validators
cudos-noded add-genesis-account $VALIDATOR_ADDRESS "100000000${BOND_DENOM}"
cudos-noded gentx root-validator "100000000${BOND_DENOM}" --chain-id $CHAIN_ID --keyring-backend test

cudos-noded keys add faucet --keyring-backend test |& tee "${CUDOS_HOME}/faucet.wallet"
FAUCET_ADDRESS=$(cudos-noded keys show faucet -a)
cudos-noded add-genesis-account $FAUCET_ADDRESS "100000000000000000000000${BOND_DENOM}"

cudos-noded collect-gentxs
