# faucet

A faucet that uses [cosmos-sdk](https://github.com/cosmos/cosmos-sdk) executable binaries only.

The main purpose of this `faucet` is to avoid using RPC or API endpoints, and use the CLI binary instead, more
specifically, the commands:

```bash
$ {app}d tx bank send
```

and:

```bash
$ {app}d query txs
```

Since the faucet only uses the CLI binary, it is compatible with practically any blockchain built with
[cosmos-sdk](https://github.com/cosmos/cosmos-sdk) even if different types of keys are used (such as in
[ethermint](https://github.com/cosmos/ethermint) for example).


## Installation

### Using cURL

```bash
$ curl https://get.starport.network/faucet! | bash 
```

### Use docker image

Use docker image `ghcr.io/tendermint/faucet`. You can use it in a Kubernetes pod with 
[shareProcessNamespace](https://kubernetes.io/docs/tasks/configure-pod-container/share-process-namespace/#configure-a-pod)
or mount the chain binary using docker:

```bash
$ docker run -it -v ~/go/bin/gaiad:/usr/local/bin/gaiad ghcr.io/tendermint/faucet
```

### From Source
You can build the faucet with:

```bash
$ make build
```

The executable binary will be avaialable in the `./build/` directory. To install it to `$GOPATH/bin`, use instead:

```bash
$ make install
```

## Usage

### Configuration

You can configure the faucet either using command line flags or environment variables. The following table
shows the available configuration options and respective defaults:

| flag             	| env              	| description                                                   | default                      	|
|------------------	|------------------	|--------------------------------------------------------------	|------------------------------	|
| port             	| PORT             	| tcp port where faucet will be listening for requests 	        | 8000                         	|
| account-name    	| ACCOUNT_NAME     	| name of the account to be used by the faucet                  | faucet                       	|
| mnemonic         	| MNEMONIC         	| mnemonic for restoring an account         	                |                              	|
| keyring-password 	| KEYRING_PASSWORD 	| password for accessing keyring                      	        |                              	|
| cli-name         	| CLI_NAME         	| name of the cli executable                                   	| gaiad 	                    |
| denom            	| DENOM            	| denomination of the coins sent by default (comma separated)  	| uatom                        	|
| credit-amount    	| CREDIT_AMOUNT    	| amount to credit in each request                             	| 10000000                     	|
| max-credit       	| MAX_CREDIT       	| maximum credit per account                                   	| 100000000                    	|
| sdk-version      	| SDK_VERSION      	| version of sdk (launchpad or stargate)                        | stargate                    	|
| node            	| NODE            	| address of tendermint RPC endpoint for this chain             |                    	        |
| keyring-backend   | KEYRING_BACKEND   | keyring backend to be used                                    |                               |
| legacy-send       | LEGACY_SEND       | whether to use legacy send command                            | false                         |
|                   |                   |                                                               |                               |

### [gaia](https://github.com/cosmos/gaia) example

This faucet options default to work with [gaia](https://github.com/cosmos/gaia). So you can start the faucet with just:

```bash
$ faucet --keyring-password 12345678
INFO[0000] listening on :8000
```

or, with environment variables:

```bash
$ export KEYRING_PASSWORD=12345678
$ faucet
INFO[0000] listening on :8000
```

### [ethermint](https://github.com/cosmos/ethermint) example

Start the faucet with:

```bash
$ faucet --cli-name ethermintcli --denoms ueth --keyring-password 12345678 --sdk-version launchpad
INFO[0000] listening on :8000
```

or, with environment variables:

```bash
$ export CLI_NAME=ethermintcli
$ export SDK_VERSION=launchpad
$ export DENOMS=ueth
$ export KEYRING_PASSWORD=12345678
$ faucet
INFO[0000] listening on :8000
```

### [wasmd](https://github.com/CosmWasm/wasmd) example

Start the faucet with:

```bash
$ faucet --cli-name wasmcli --denoms ucosm --keyring-password 12345678
INFO[0000] listening on :8000
```

or, with environment variables:

```bash
$ export CLI_NAME=wasmcli
$ export DENOMS=ucosm
$ export KEYRING_PASSWORD=12345678
$ faucet
INFO[0000] listening on :8000
```

### Request tokens

You can request tokens by sending a `POST` request to the faucet, with a key address in a `JSON`:

```bash
$ curl -X POST -d '{"address": "cosmos1kd63kkhtswlh5vcx5nd26fjmr9av74yd4sf8ve"}' http://localhost:8000
{"transfers":[{"coin":"10000000uatom","status":"ok"}]}
```

For requesting specific coins, use:

```bash
$ curl -X POST -d '{"address": "cosmos1kd63kkhtswlh5vcx5nd26fjmr9av74yd4sf8ve", "coins": ["10uatom", "20ueth"]}' http://localhost:8000
{"transfers":[{"coin":"10uatom","status":"ok"}, {"coin":"20ueth","status":"ok"}]}
```