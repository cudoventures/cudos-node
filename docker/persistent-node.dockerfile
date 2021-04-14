FROM golang:1.15

ENV PROTOC_ZIP=protoc-3.13.0-linux-x86_64.zip
RUN apt-get update && apt-get install -y unzip
RUN curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.13.0/$PROTOC_ZIP \
    && unzip -o $PROTOC_ZIP -d /usr/local bin/protoc \
    && unzip -o $PROTOC_ZIP -d /usr/local 'include/*' \ 
    && rm -f $PROTOC_ZIP

WORKDIR /usr/blockchain

COPY ./ ./

RUN make

RUN cudos-noded init cudos-node-network --chain-id=cudos-node-network

RUN cudos-noded keys add validator01 --keyring-backend test

RUN VALIDATOR_ADDRESS=$(cudos-noded keys show validator01 -a --keyring-backend test)&& cudos-noded add-genesis-account $VALIDATOR_ADDRESS 100000000000stake

RUN cudos-noded keys add account-vesting-01 --keyring-backend test

RUN VESTING_ADDRESS=$(cudos-noded keys show account-vesting-01 -a --keyring-backend test)&& cudos-noded add-genesis-account $VESTING_ADDRESS 1000stake --vesting-amount 500stake --vesting-end-time 1617613800

RUN cudos-noded gentx validator01 100000000stake --chain-id cudos-node-network --keyring-backend test

RUN cudos-noded collect-gentxs

CMD ["cudos-noded", "start"] 