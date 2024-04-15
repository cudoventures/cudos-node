set -e -u -o pipefail

export CUDOS_NODED_KEYRING_BACKEND=test
export CUDOS_NODED_CHAIN_ID=$CHAIN_ID
export LD_LIBRARY_PATH=/app
export PATH=/app:$PATH

update_genesis() {
  cat /app/cudos-data/config/genesis.json | jq "$@" | sponge /app/cudos-data/config/genesis.json
}

update_app() {
  cat /app/cudos-data/config/app.toml | tomlq -t "$@" | sponge /app/cudos-data/config/app.toml
}

update_config() {
  cat /app/cudos-data/config/config.toml | tomlq -t "$@" | sponge /app/cudos-data/config/config.toml
}

if ! [ -d cudos-data ]; then
  echo Initialising new node

  cudos-noded init $NODE_MONIKER

  cudos-noded keys add validator

  ## Note that the validator MUST have acudos, regardless of what other tokens might be use,
  ## because the software specifically requires a self-delegation of 2 million CUDOS.
  cudos-noded add-genesis-account $(cudos-noded keys show validator -a) 2000000000000000000000000000acudos;

  ## ATTENTION DUELISTS!
  ## Change these to give yourself tokens if you need to. These are MY accounts.
  cudos-noded add-genesis-account cudos1p8ux3usg99fnu7e5ga72dnpskrf9jw7w6yxu7w 2000000000000000000000000000$DENOM;
  cudos-noded add-genesis-account cudos1genudzpvqe2t9k64xwueua35a8kfvl3fc6uc62 2000000000000000000000000000$DENOM;
  cudos-noded add-genesis-account cudos1m4kxqu2fhh0z0af5jlkhhy7e0qlanxkhvt074v 2000000000000000000000000000$DENOM;
  cudos-noded add-genesis-account cudos15yk64u7zc9g9k2yr2wmzeva5qgwxps6y8rxs8t 2000000000000000000000000000$DENOM;

  cudos-noded gentx validator \
    --min-self-delegation 2000000000000000000000000 \
    2000000000000000000000000acudos \
    0x0000000000000000000000000000000000000000 \
    $(cudos-noded keys show validator -a)

  update_genesis ".chain_id = \"$CHAIN_ID\""
  update_genesis '.app_state.staking.params.bond_denom = "acudos"'
  update_genesis --arg address $(cudos-noded keys show validator -a) '.app_state.gravity.static_val_cosmos_addrs = [$address]'
  update_genesis '.app_state.gravity.erc20_to_denoms = [{erc20: "0x28ea52f3ee46CaC5a72f72e8B3A387C0291d586d", denom: "acudos"}]'

  update_config '.consensus.timeout_commit = "1s"'
  update_config '.rpc.laddr = "tcp://0.0.0.0:26657"'
  update_config '.rpc.cors_allowed_origins = ["*"]'

  update_app '."minimum-gas-prices" = "1'$DENOM'"'
  update_app '.api.enable = true'
  update_app '.api."enabled-unsafe-cors" = true'

  cudos-noded collect-gentxs
fi

exec cudos-noded start
