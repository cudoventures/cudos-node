module cudos.org/cudos-node

go 1.15

replace github.com/althea-net/cosmos-gravity-bridge/module => ../CudosGravityBridge/module

replace github.com/CosmWasm/wasmd => github.com/provenance-io/wasmd v0.17.1-0.20210812214331-ce3a93a9268d

require (
	github.com/CosmWasm/wasmd v0.17.0
	github.com/althea-net/cosmos-gravity-bridge/module v0.0.0-00010101000000-000000000000
	// github.com/althea-net/cosmos-gravity-bridge/module v0.0.0-20210611212501-5de966823aa6
	github.com/cosmos/cosmos-sdk v0.43.0
	github.com/cosmos/ibc-go v1.0.0
	github.com/gogo/protobuf v1.3.3
	github.com/golang/glog v0.0.0-20210429001901-424d2337a529 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/regen-network/cosmos-proto v0.3.1
	github.com/spf13/cast v1.4.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.11
	github.com/tendermint/tm-db v0.6.4
	google.golang.org/genproto v0.0.0-20210617175327-b9e0b3197ced
	google.golang.org/grpc v1.38.0

)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
