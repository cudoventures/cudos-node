module github.com/CudoVentures/cudos-node

go 1.15

require (
	github.com/CosmWasm/wasmd v0.25.0
	github.com/althea-net/cosmos-gravity-bridge/module v0.0.0-00010101000000-000000000000
	github.com/containerd/continuity v0.2.2 // indirect
	github.com/cosmos/cosmos-sdk v0.45.3
	github.com/cosmos/ibc-go/v2 v2.2.0
	github.com/dustin/go-humanize v1.0.1-0.20200219035652-afde56e7acac // indirect
	github.com/gogo/protobuf v1.3.3
	github.com/golang/glog v1.0.0 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/google/btree v1.0.1 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/mattn/go-runewidth v0.0.10 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/onsi/gomega v1.15.0 // indirect
	github.com/opencontainers/runc v1.1.0 // indirect
	github.com/regen-network/cosmos-proto v0.3.1
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/rs/zerolog v1.26.1 // indirect
	github.com/spf13/afero v1.8.0 // indirect
	github.com/spf13/cast v1.4.1
	github.com/spf13/cobra v1.4.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.8.0
	github.com/tendermint/tendermint v0.34.19
	github.com/tendermint/tm-db v0.6.7
	github.com/tidwall/gjson v1.6.7
	golang.org/x/crypto v0.0.0-20220112180741-5e0467b6c7ce // indirect
	golang.org/x/net v0.2.0 // indirect
	google.golang.org/genproto v0.0.0-20221114212237-e4508ebdbee1
	google.golang.org/grpc v1.50.1
	gopkg.in/ini.v1 v1.66.3 // indirect
	gopkg.in/yaml.v2 v2.4.0

)

// replace github.com/althea-net/cosmos-gravity-bridge/module => ../CudosGravityBridge/module
// replace github.com/cosmos/cosmos-sdk => ../cosmos-sdk
// replace github.com/tendermint/tendermint => ../tendermint
// replace github.com/cosmos/ibc-go/v2 => ../ibc-go

replace github.com/althea-net/cosmos-gravity-bridge/module => github.com/CudoVentures/cosmos-gravity-bridge/module v0.0.0-20230130131817-0381039012d6

replace github.com/cosmos/cosmos-sdk => github.com/CudoVentures/cosmos-sdk v0.0.0-20220914094638-cb70cee184ee

// replace globally the grpc version (https://docs.cosmos.network/v0.44/basics/app-anatomy.html#dependencies-and-makefile)
replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4

replace github.com/confio/ics23/go => github.com/cosmos/cosmos-sdk/ics23/go v0.8.0
