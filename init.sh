if [[ -z "${CUDOS_HOME-}" ]]; then
    CUDOS_HOME="./cudos-data"
fi

WORKING_PATH=$(pwd) && cd $CUDOS_HOME && rm -Rf ./* && cd $WORKING_PATH

# chain parameters
 MONIKER="cudos-root-node"
CHAIN_ID="cudos-network"
TIMEOUT_COMMIT="5s" #5s originally
KEYPASSWD='123123123' #at least 8 characters

MIN_SELF_DELEGATION="1" # minimum tokens sto stake multiplyer by 1 000 000 for validator01

# slashing parameters
JAIL_DURATION="600s" #600s originally
SIGNED_BLOCKS_WINDOW="400" #100
MIN_SIGNED_PER_WINDOW="0.1" #0.5

# staking parameters
BOND_DENOM="acudos" # stake originally
UNBONDING_TIME="1814400s" #1814400s originally
MAX_VALIDATORS="100" # 100 originally

# government parameters
GOV_PROPOSAL_MIN_DEPOSIT_DENOM="acudos" # stake orginally
GOV_PROPOSAL_MIN_DEPOSIT_AMOUNT="10000000" # 10000000 originally
GOV_PROPOSAL_MAX_DEPOSIT_PERIOD="600s" # 172800s originally
GOV_PROPOSAL_VOTING_PERIOD="600s" # 172800s originally
GOV_QUORUM="0.334000000000000000" # 0.334000000000000000 originally
GOV_THRESHOLD="0.500000000000000000" # 0.500000000000000000 originally
GOV_VETO_THRESHOLD="0.334000000000000000" # 0.334000000000000000 originally

# mint parameters
MINT_DENOM="acudos" # stake originally
MINT_INFLATION="0.0000000013" # 0.130000000000000000 originally
MINT_INFLATION_RATE_CHANGE="0.0000000013" # 0.130000000000000000 originally
MINT_INFLATION_MAX="0.0000000000" # 0.200000000000000000 originally
MINT_INFLATION_MIN="0.0000000000" # 0.070000000000000000 originally
MINT_GOAL_BONDED="0.670000000000000000" # 0.670000000000000000 originally
BLOCKS_PER_YEAR="6311520" # 6311520 originally
BANK_SYMBOL="cudos"
BANK_NAME="cudos"

# bank parameters
BANK_SEND_ENABLED='[{"denom": "cudosAdmin", "enabled": false}]'

DENOM_METADATA_DESC="The native staking token of the Cudos Hub."
DENOM1="acudos" EXP1="0" ALIAS1="attocudos"
DENOM2="fcudos" EXP2="3" ALIAS2="femtocudos"
DENOM3="pcudos" EXP3="6" ALIAS3="picocudos"
DENOM4="ncudos" EXP4="9" ALIAS4="nanocudos"
DENOM5="ucudos" EXP5="12" ALIAS5="microcudos"
DENOM6="mcudos" EXP6="15" ALIAS6="millicudos"
DENOM7="cudos" EXP7="18"
BASE="acudos"
DISPLAY="cudos"

cudos-noded init $MONIKER --chain-id=$CHAIN_ID

# sed -i "104s/enable = false/enable = true/" "${CUDOS_HOME}/config/app.toml"
# sed -i "s/laddr = \"tcp:\/\/127.0.0.1:26657\"/laddr = \"tcp:\/\/0.0.0.0:26657\"/" "${CUDOS_HOME}/config/config.toml"

# enable cors origin for local testing
# sed -i "s/enabled-unsafe-cors = false/enabled-unsafe-cors = true/" ${CUDOS_HOME}/config/app.toml
# sed -i "s/cors_allowed_origins = \[\]/cors_allowed_origins = \[\"\*\"\]/" ${CUDOS_HOME}/config/config.toml

# setting time after commit before proposing a new block
sed -i '' "s/timeout_commit = \"5s\"/timeout_commit = \"$TIMEOUT_COMMIT\"/" "${CUDOS_HOME}/config/config.toml"

# setting slashing time
cat "${CUDOS_HOME}/config/genesis.json" | jq --arg JAIL_DURATION "$JAIL_DURATION" '.app_state.slashing.params.downtime_jail_duration = $JAIL_DURATION' > "${CUDOS_HOME}/config/tmp_genesis.json" && mv "${CUDOS_HOME}/config/tmp_genesis.json" "${CUDOS_HOME}/config/genesis.json"
cat "${CUDOS_HOME}/config/genesis.json" | jq --arg SIGNED_BLOCKS_WINDOW "$SIGNED_BLOCKS_WINDOW" '.app_state.slashing.params.signed_blocks_window = $SIGNED_BLOCKS_WINDOW' > "${CUDOS_HOME}/config/tmp_genesis.json" && mv "${CUDOS_HOME}/config/tmp_genesis.json" "${CUDOS_HOME}/config/genesis.json"
cat "${CUDOS_HOME}/config/genesis.json" | jq --arg MIN_SIGNED_PER_WINDOW "$MIN_SIGNED_PER_WINDOW" '.app_state.slashing.params.min_signed_per_window = $MIN_SIGNED_PER_WINDOW' > "${CUDOS_HOME}/config/tmp_genesis.json" && mv "${CUDOS_HOME}/config/tmp_genesis.json" "${CUDOS_HOME}/config/genesis.json"

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

# setting bank params
cat "${CUDOS_HOME}/config/genesis.json" | jq  --argjson BANK_SEND_ENABLED "$BANK_SEND_ENABLED" '.app_state.bank.params.send_enabled = $BANK_SEND_ENABLED' > "${CUDOS_HOME}/config/tmp_genesis.json" && mv "${CUDOS_HOME}/config/tmp_genesis.json" "${CUDOS_HOME}/config/genesis.json"

cat "${CUDOS_HOME}/config/genesis.json" | jq  --arg BANK_SYMBOL "$BANK_SYMBOL" '.app_state.bank.denom_metadata[0].symbol = $BANK_SYMBOL' > "${CUDOS_HOME}/config/tmp_genesis.json" && mv "${CUDOS_HOME}/config/tmp_genesis.json" "${CUDOS_HOME}/config/genesis.json"
cat "${CUDOS_HOME}/config/genesis.json" | jq  --arg BANK_NAME "$BANK_NAME" '.app_state.bank.denom_metadata[0].name = $BANK_NAME' > "${CUDOS_HOME}/config/tmp_genesis.json" && mv "${CUDOS_HOME}/config/tmp_genesis.json" "${CUDOS_HOME}/config/genesis.json"


# setting fractions metadata
cat "${CUDOS_HOME}/config/genesis.json" | jq --arg DENOM_METADATA_DESC "$DENOM_METADATA_DESC" --arg DENOM1 "$DENOM1" --arg EXP1 "$EXP1" --arg ALIAS1 "$ALIAS1" --arg DENOM2 "$DENOM2" --arg EXP2 "$EXP2" --arg ALIAS2 "$ALIAS2" --arg DENOM3 "$DENOM3" --arg EXP3 "$EXP3" --arg ALIAS3 "$ALIAS3" --arg DENOM4 "$DENOM4" --arg EXP4 "$EXP4" --arg ALIAS4 "$ALIAS4" --arg DENOM5 "$DENOM5" --arg EXP5 "$EXP5" --arg ALIAS5 "$ALIAS5" --arg DENOM6 "$DENOM6" --arg EXP6 "$EXP6" --arg ALIAS6 "$ALIAS6" --arg DENOM7 "$DENOM7" --arg EXP7 "$EXP7" --arg BASE "$BASE" --arg DISPLAY "$DISPLAY" '.app_state.bank.denom_metadata[0].description=$DENOM_METADATA_DESC | .app_state.bank.denom_metadata[0].denom_units[0].denom=$DENOM1 | .app_state.bank.denom_metadata[0].denom_units[0].exponent=$EXP1 | .app_state.bank.denom_metadata[0].denom_units[0].aliases[0]=$ALIAS1 | .app_state.bank.denom_metadata[0].denom_units[1].denom=$DENOM2 | .app_state.bank.denom_metadata[0].denom_units[1].exponent=$EXP2 | .app_state.bank.denom_metadata[0].denom_units[1].aliases[0]=$ALIAS2 | .app_state.bank.denom_metadata[0].denom_units[2].denom=$DENOM3 | .app_state.bank.denom_metadata[0].denom_units[2].exponent=$EXP3 | .app_state.bank.denom_metadata[0].denom_units[2].aliases[0]=$ALIAS3 | .app_state.bank.denom_metadata[0].denom_units[3].denom=$DENOM4 | .app_state.bank.denom_metadata[0].denom_units[3].exponent=$EXP4 | .app_state.bank.denom_metadata[0].denom_units[3].aliases[0]=$ALIAS4 | .app_state.bank.denom_metadata[0].denom_units[4].denom=$DENOM5 | .app_state.bank.denom_metadata[0].denom_units[4].exponent=$EXP5 | .app_state.bank.denom_metadata[0].denom_units[4].aliases[0]=$ALIAS5 | .app_state.bank.denom_metadata[0].denom_units[5].denom=$DENOM6 | .app_state.bank.denom_metadata[0].denom_units[5].exponent=$EXP6 | .app_state.bank.denom_metadata[0].denom_units[5].aliases[0]=$ALIAS6 | .app_state.bank.denom_metadata[0].denom_units[6].denom=$DENOM7 | .app_state.bank.denom_metadata[0].denom_units[6].exponent=$EXP7 | .app_state.bank.denom_metadata[0].base=$BASE | .app_state.bank.denom_metadata[0].display=$DISPLAY'  > "${CUDOS_HOME}/config/tmp_genesis.json" && mv "${CUDOS_HOME}/config/tmp_genesis.json" "${CUDOS_HOME}/config/genesis.json"

# add a new key entry from which to make validator
cudos-noded keys add root-validator-01 
ROOT_VALIDATOR_01_ADDRESS=$(  cudos-noded keys show root-validator-01 -a  )

cudos-noded keys add validator-02 
VALIDATOR_02_ADDRESS=$(  cudos-noded keys show validator-02 -a  )
  cudos-noded keys add validator-03
VALIDATOR_03_ADDRESS=$(  cudos-noded keys show validator-03 -a  )

  cudos-noded keys add private-sale-offer
PRIVATE_SALE_OFFER_ADDRESS=$(  cudos-noded keys show private-sale-offer -a  )

# create validators
cudos-noded add-genesis-account $ROOT_VALIDATOR_01_ADDRESS "100000000000000000000100000${BOND_DENOM},1cudosAdmin"
cudos-noded add-genesis-account $VALIDATOR_02_ADDRESS "100000000000000100000${BOND_DENOM},1cudosAdmin"
cudos-noded add-genesis-account $VALIDATOR_03_ADDRESS "100000000000000100000${BOND_DENOM},1cudosAdmin"
cudos-noded gentx root-validator-01 "1000000000000000000000${BOND_DENOM}" "0x9fdE6D55dDa637806DbF016a03B6970613630333" "cudos1g6wvv0fk2q5m2xh0yk2app0ewp6f6aum7v2j7y" --chain-id $CHAIN_ID


# test validator creation on genesis
cudos-noded keys add genesis-validator
GENESIS_VALIDATOR_ADDRESS=$(cudos-noded keys show genesis-validator -a)
cudos-noded add-genesis-account $GENESIS_VALIDATOR_ADDRESS "9152921504606846975000000${BOND_DENOM}"


# add faucet account
cudos-noded keys add faucet
FAUCET_ADDRESS=$(  cudos-noded keys show faucet -a  )
cudos-noded add-genesis-account $FAUCET_ADDRESS "1000000000000000000000000${BOND_DENOM}"

#cudos-noded gentx buidler-moniker '100000000000000000000acudos' --chain-id 'cudos-network' \
#  --commission-rate '0.10' \
#  --commission-max-rate '0.20' \
#  --commission-max-change-rate '0.02' \
#  --pubkey 'cudosvalconspub1zcjduepqmh9qn9n77z3ngc0n4lvyunuaj34pdrtcka8ywh6avjdg9x9lmvzskg0s47'


cudos-noded collect-gentxs

sed -i '' "s/pex = true/pex = false/" "${CUDOS_HOME}/config/config.toml"

MY_OWN_PEER_ID=$(cudos-noded tendermint show-node-id)
sed -i '' "s/private_peer_ids = \"\"/private_peer_ids = \"$MY_OWN_PEER_ID\"/g" "${CUDOS_HOME}/config/config.toml"

cudos-noded tendermint show-node-id "${CUDOS_HOME}/tendermint.nodeid"

chmod 755 "${CUDOS_HOME}/config"
cudos-noded start --home ./cudos-data/