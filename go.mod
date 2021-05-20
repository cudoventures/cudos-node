module cudos.org/cudos-node

go 1.15

require (
	github.com/CosmWasm/wasmd v0.16.0
	github.com/cosmos/cosmos-sdk v0.42.4
	github.com/cosmos/gravity-bridge v0.1.3-0.20210512001249-3c4c5d1dbe2e // indirect
	github.com/cosmos/gravity-bridge/module v0.0.0-20210512001249-3c4c5d1dbe2e
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/regen-network/cosmos-proto v0.3.1
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/tendermint/tendermint v0.34.10
	github.com/tendermint/tm-db v0.6.4
	google.golang.org/genproto v0.0.0-20210207032614-bba0dbe2a9ea
	google.golang.org/grpc v1.37.0

)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
