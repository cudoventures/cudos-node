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

RUN cudos-poc-01d init cudos-poc-01-network --chain-id=cudos-poc-01-network

RUN rm ~/.blog/config/genesis.json

ARG PERSISTENT_NODE_ID

RUN sed "s/persistent_peers = \"\"/persistent_peers = \"$PERSISTENT_NODE_ID\"/g" ~/.blog/config/config.toml > ./modified-config.toml&& cp ./modified-config.toml ~/.blog/config/config.toml&& rm ./modified-config.toml

COPY --from=cudos-network-persistent-node /root/.blog/config/genesis.json /root/.blog/config/genesis.json

CMD ["cudos-poc-01d", "start"]