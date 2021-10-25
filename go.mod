module cudos.org/cudos-node

go 1.15

replace github.com/althea-net/cosmos-gravity-bridge/module => ../CudosGravityBridge/module

replace github.com/CosmWasm/wasmd => github.com/provenance-io/wasmd v0.17.1-0.20210812214331-ce3a93a9268d

require (
	github.com/CosmWasm/wasmd v0.17.0
	github.com/althea-net/cosmos-gravity-bridge/module v0.0.0-00010101000000-000000000000
	// github.com/althea-net/cosmos-gravity-bridge/module v0.0.0-20210611212501-5de966823aa6
	github.com/cosmos/cosmos-sdk v0.44.3
	github.com/cosmos/ibc-go v1.0.0
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/regen-network/cosmos-proto v0.3.1
	github.com/spf13/cast v1.4.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/tendermint/tendermint v0.34.14
	github.com/tendermint/tm-db v0.6.4
	google.golang.org/genproto v0.0.0-20210602131652-f16073e35f0c
	google.golang.org/grpc v1.40.0

)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
