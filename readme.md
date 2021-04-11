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
