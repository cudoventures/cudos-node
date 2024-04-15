FROM ubuntu:24.04

WORKDIR /app

RUN apt-get update
RUN apt-get install -y jq moreutils yq

# TODO: build binary file here
COPY ./build/cudos-noded ./tests/cosmos/libwasmvm.so ./tests/cosmos/init.sh ./

ENV NODE_MONIKER=node
ENV CHAIN_ID=cudos-1
ENV DENOM=acudos

CMD bash init.sh

# RPC
EXPOSE 26657

# LCD
EXPOSE 1317

# Peer
EXPOSE 26656
