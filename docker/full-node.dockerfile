FROM golang:1.15

WORKDIR /usr/blockchain

COPY ./ ./

RUN make

RUN cudos-poc-01d init cudos-poc-01-network --chain-id=cudos-poc-01-network

RUN rm ~/.blog/config/genesis.json

ARG PERSISTENT_NODE_ID

RUN sed "s/persistent_peers = \"\"/persistent_peers = \"$PERSISTENT_NODE_ID\"/g" ~/.blog/config/config.toml > ./modified-config.toml&& cp ./modified-config.toml ~/.blog/config/config.toml&& rm ./modified-config.toml

COPY --from=cudos-network-persistent-node /root/.blog/config/genesis.json /root/.blog/config/genesis.json

CMD ["cudos-poc-01d", "start"]