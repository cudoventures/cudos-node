module github.com/CudoVentures/cudos-node

go 1.15

// replace github.com/althea-net/cosmos-gravity-bridge/module => ../CudosGravityBridge/module

replace github.com/althea-net/cosmos-gravity-bridge/module => github.com/CudoVentures/cosmos-gravity-bridge/module v0.0.0-20220411110512-445f2a71933e

replace github.com/CosmWasm/wasmd => github.com/provenance-io/wasmd v0.17.1-0.20210812214331-ce3a93a9268d

// replace github.com/tendermint/tendermint => github.com/osmosis-labs/tendermint v0.34.14-starving-waiters

// replace github.com/cosmos/cosmos-sdk => ../cosmos-sdk

replace github.com/cosmos/cosmos-sdk => github.com/CudoVentures/cosmos-sdk v0.0.0-20220411110303-34c676cfb016

// replace globally the grpc version (https://docs.cosmos.network/v0.44/basics/app-anatomy.html#dependencies-and-makefile)
replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

require (
	github.com/CosmWasm/wasmd v0.17.0
	github.com/althea-net/cosmos-gravity-bridge/module v0.0.0-00010101000000-000000000000
	github.com/cosmos/cosmos-sdk v0.44.5
	github.com/cosmos/ibc-go v1.2.3
	github.com/gogo/protobuf v1.3.3
	github.com/golang/glog v1.0.0 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/regen-network/cosmos-proto v0.3.1
	github.com/spf13/cast v1.4.1
	github.com/spf13/cobra v1.4.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.1
	github.com/tendermint/tendermint v0.34.19
	github.com/tendermint/tm-db v0.6.6
	github.com/tidwall/gjson v1.6.7
	google.golang.org/genproto v0.0.0-20220118154757-00ab72f36ad5
	google.golang.org/grpc v1.45.0

)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
