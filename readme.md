# Cosmos upgradability PoC

This is a blockchain built using Cosmos SDK and Tendermint and created with [Starport](https://github.com/tendermint/starport) in order to test Cosmos upgradability concepts.

## Links

 * [Upgrade module documentation](https://docs.cosmos.network/v0.42/modules/upgrade/)
 * [Governance module documentation](https://docs.cosmos.network/v0.42/modules/gov/)
 * [How cosmos goverance works?](https://lunie.io/guides/how-cosmos-governance-works/)
 * [Cosmovisor upgrade utility](https://github.com/cosmos/cosmos-sdk/tree/master/cosmovisor)
 * [IBC Chains Upgradability overview](https://docs.cosmos.network/master/ibc/upgrades/quick-guide.html)
 * [Chain upgrade guide](https://docs.cosmos.network/v0.42/migrations/chain-upgrade-guide-040.html)


## Cosmos upgrade procedure overview

* Upgrades are facilitated via the Cosmos's upgrade module. This module's keeper stores plans that have either upgrade height or upgrade time specified. Furthermore, this module has a hook in the `BeginBlocker` that that watches for such plans and execute specific upgrade handlers or halts the blockhain state machine if no upgrade handler has been defined.

* The upgrade module does not provide methods for creating upgrade plans directly (this would be an easy way to crash a public network), instead one must go through the government module's proposal. A brief overview on how the software upgrade proposal workflow can be found at the specified [lunie overview](https://lunie.io/guides/how-cosmos-governance-works/), but to reiterate:
1. [Submit Stage] Somebody submits a proposal (in this case software-upgrade proposal). Standart transaction fees apply.
2. [Deposit Stage] In order for this proposal to be eligible for voting, certain stake must be deposited. Generally, this stake is refunded in the end unless sufficient veto is applied. Default values for depositing:
 ```
       "deposit_params": {
        "max_deposit_period": "172800s",
        "min_deposit": [
          {
            "amount": "10000000",
            "denom": "stake"
          }
        ]
      }
``` 
3. [Voting Stage] If required deposit amount is reached in the deposit period, participants who have stake (i.e. bonded tokens) can vote for the proposal with `Yes`, `No`, `NoWithVeto` or `Abstain`.
   Default values for voting:
 ```
       "voting_params": {
        "voting_period": "172800s"
      }
```
4. [Tally Stage] Once the voting time ends, at the end of the next block the gov module, via the `EndBlcoker` hook will calculate the result. The vote counting is two-dimentional, i.e. it counts a quorum and threshold. It also uses the term voting power. Voting power, as per defined in the gov module source is
   - for delegates - (delegation shares * validator bonded tokens) / total shares
   - for validators - (validator's owned shares * validator bonded tokens) / total shares
  
   The quorum is the minimum percentage of voting power that needs to be casted on a proposal for the result to be valid. If less than 1/3 of the voting power holders have voted the voting becomes invalid and deposited stakes are refunded. (nb. it bears no importance how many votes are casted, even if only node votes however having > 1/3 of the total voting power, the quorum is fulfilled).
  Threshold is defined as the minimum proportion of Yes votes (excluding Abstain votes) for the proposal to be accepted. Initially, the threshold is set at 50% with a possibility to veto if more than 1/3rd of votes (excluding Abstain votes) are NoWithVeto votes. Important - Vote options, again, are calculated based on the voting power. The actual source calculates a vote options as in current_vote_option + voting_power_for_this_option. Interesting side effect of this is that if 1 node has > than the veto threshold (by default 1/3), it can by itself enact it, given it votes NoWithVeto.
  Default values for tallying:
 ```
       "tally_params": {
        "quorum": "0.334000000000000000",
        "threshold": "0.500000000000000000",
        "veto_threshold": "0.334000000000000000"
      }
 ```
5. [Upgrade Stage] If the tally results in a positive vote, a new plan is created in the upgrade module and executed as explained. Important: We need to be inline with the IBC integration requirements, such as that only upgrade height can be used as well as no ad-hoc upgrades should be performed.
   ! In the gov documentation there is a section about signalling and switching during the `SoftwareUpgradeProposal`. I cannot seem to find any evidence of such mechanisms. Upgrades should be handled by the upgrade module and done on the beggining of a deterministic block.

* Migration - There is a separate migration procedure as described in the [cosmos documentation](https://docs.cosmos.network/v0.42/migrations/chain-upgrade-guide-040.html) and [cosmos upgrade migrations plans](https://github.com/cosmos/gaia/blob/main/docs/migration/cosmoshub-3.md#upgrade-procedure). It is different from what is described here. While there is a proposal, the proposal is in text form and it is self imposed by modyfing the app.toml config, e.g:
 ```
 Make sure your chain halts at the right time and date: February 18, 2021 at 06:00 UTC is in UNIX seconds: 1613628000

perl -i -pe 's/^halt-time =.*/halt-time = 1613628000/' ~/.gaia/config/app.toml
```
This type of upgrate also allows for backwards incompatible changes. However, because one has to export the whole chain, it will take some time (30+ mins for Cosmos hub mainnet is reported).

## Upgrade example

---

Run (and then stop):
```
starport serve
```

`serve` command installs dependencies, builds, initialises and starts your blockchain in development. The generated genesis has voting period of `30s` in order to speed up the upgrade demonstration.

Afterwards install and configure the cosmovisor tools:
```
# Install cosmovisor
go get github.com/cosmos/cosmos-sdk/cosmovisor/cmd/cosmovisor
# Export necessary env variables
export DAEMON_NAME=pocbasecosmosd
export DAEMON_HOME=~/.pocbasecosmos/
export DAEMON_ALLOW_DOWNLOAD_BINARIES=true

# Create necesary cosmovisor folder structure
mkdir -p $DAEMON_HOME/cosmovisor/genesis/bin
cp $GOPATH/bin/poc-base-cosmosd $DAEMON_HOME/cosmovisor/genesis/bin

# Start cosmovisor
cosmovisor start
```
Once the blockchain is running, create a proposal and directly deposit the needed amount via
```
  poc-base-cosmosd tx gov submit-proposal software-upgrade test1 \
  --title "upgrade-demo" \
  --description "testing upgrade functionality"\
  --upgrade-info '{"binaries":{"linux/amd64":"https://github.com/rdpnd/cosmos-poc/raw/master/poc-base-cosmosd"}}'\
  --from alice\
  --upgrade-height 100\
  --deposit 10000000stake\
  --chain-id pocbasecosmos\
  -y
```
Afterwards vote on the proposal:
```
poc-base-cosmosd tx gov vote 1 yes --from alice
```
You can check the proposal status via:
```
poc-base-cosmosd q gov proposal 1
```
After reaching height 100, the blockchain should halt, cosmovisor will try to download the provided binary from the upgrade-info (if the user is operating on a linux/amd64 machine) and then symlink it as the current binary then exit as well. Upon issuing `cosmovisor start` again, the new version of the blockchain shall be running. The new version had the following code inserted in app/app.go's New function:
```go
	app.UpgradeKeeper.SetUpgradeHandler("test1", func(ctx sdk.Context, plan upgradetypes.Plan) {
		// Add some coins to a random account
		addr, err := sdk.AccAddressFromBech32("cosmos1kqse03ju8evqw0ed3f92dk5suf4fj24lfy8y0q")
		if err != nil {
			panic(err)

		}
		err = app.BankKeeper.AddCoins(ctx, addr, sdk.Coins{sdk.Coin{Denom: "stake", Amount: sdk.NewInt(345600000)}})
		if err != nil {
			panic(err)
		}
	})
```