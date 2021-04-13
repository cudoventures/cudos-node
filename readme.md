# Cosmos upgradability PoC

This is a blockchain built using Cosmos SDK and Tendermint and created with [Starport](https://github.com/tendermint/starport) in order to test Cosmos treasury capabilitites. The `pocbasecosmos` module has a basic functionality of allowing users with cudosAdmin tokens to send funds from Community pool 

Links

* [Cosmos hub community pool](https://github.com/gavinly/CosmosCommunitySpend)
* [Cosmos hub 3 - community spend](https://github.com/cosmos/governance/tree/master/community-pool-spend)
* [Distribution module documentation](https://docs.cosmos.network/v0.42/modules/distribution/)

# Cosmos treasury approaches

* Community pool - Out of the box, Cosmos provides a community pool funded by a percent of all staking rewards generated (via block rewards & transaction fees). The concrete percent is defined via the param communitytax and is 2% by default. Out of the box, this pool can be operated only via governemnt proposal in the same way as the upgradability concept. This means that there is no single authority which can singlehandedly operate the funds wihout the validator approval. In order to facilitate a centralized governance over the community pool, we have created a poc module with admin capabilities. This module allows users to move funds from the community pool to an arbitrary address IFF they posses positive amount of "CudosAdmin" tokens. Our proposal is for these admin tokens to only be generated in the genesis. Optionally afterwards, we can make so that the owner of such tokens can burn them if they choose to as a way of gracefully deprecating the centralized fund approach.


Community tax percent param:
```
<blockchaind> q params subspace distribution communitytax                                                 
key: communitytax
subspace: distribution
value: '"0.020000000000000000"'
```

Cudos initialization via the config.yml file:
```
accounts:
  - name: alice
    coins: ["1000token", "10000000000stake", "1cudosAdmin"]
...
```
Note that we have also explicitly forbidden to trade this admin coins via:
```
genesis:
    bank:
      params:
        send_enabled:
          - denom: cudosAdmin
            enabled: false
```

Useage:
```
# When user possesing admin coins tries to spend community funds:
╰─$ <binaryd> tx poc-module admin-spend cosmos10z570tekyzemtzzgx4sx9mg2zml7l8rz0xgrg4 2stake --from alice -y
{"height":"312","txhash":"2EBC4EA67E1437D6E061BA376519CD90AC2F527CC834259A07E53FD98FF01791","codespace":"","code":0,"data":"0A190A1761646D696E5370656E64436F6D6D756E697479506F6F6C","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"adminSpendCommunityPool\"},{\"key\":\"sender\",\"value\":\"cosmos1jv65s3grqf6v6jl3dp4t6c9t9rk99cd88lyufl\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"cosmos10z570tekyzemtzzgx4sx9mg2zml7l8rz0xgrg4\"},{\"key\":\"sender\",\"value\":\"cosmos1jv65s3grqf6v6jl3dp4t6c9t9rk99cd88lyufl\"},{\"key\":\"amount\",\"value\":\"2stake\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"adminSpendCommunityPool"},{"key":"sender","value":"cosmos1jv65s3grqf6v6jl3dp4t6c9t9rk99cd88lyufl"}]},{"type":"transfer","attributes":[{"key":"recipient","value":"cosmos10z570tekyzemtzzgx4sx9mg2zml7l8rz0xgrg4"},{"key":"sender","value":"cosmos1jv65s3grqf6v6jl3dp4t6c9t9rk99cd88lyufl"},{"key":"amount","value":"2stake"}]}]}],"info":"","gas_wanted":"200000","gas_used":"60492","tx":null,"timestamp":""}

# When user without admin coins tires to spend community funds:
╰─$ <binaryd> tx poc-module admin-spend cosmos10z570tekyzemtzzgx4sx9mg2zml7l8rz0xgrg4 2stake --from bob -y
{"height":"307","txhash":"C58C8C42CBE2514190EE4D108D8161D9001440CF288C24A88B667672157E8349","codespace":"sdk","code":4,"data":"","raw_log":"failed to execute message; message index: 0: Insufficient permissions. Address 'cosmos1t0tzgruemlaahdjvhrg6r3d6fd2gl4nevyylsv' has no cudosAdmin tokens: unauthorized","logs":[],"info":"","gas_wanted":"200000","gas_used":"46205","tx":null,"timestamp":""}

```