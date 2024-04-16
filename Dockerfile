FROM golang:1.21-bullseye AS build
ARG VERSION
ARG COMMIT

WORKDIR /app

RUN export CGO_LDFLAGS=-Wl,-rpath,$$ORIGIN/../

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -mod=readonly \
    -ldflags "-X github.com/cosmos/cosmos-sdk/version.Name=cudos-node \
              -X github.com/cosmos/cosmos-sdk/version.AppName=cudos-noded \
              -X github.com/cosmos/cosmos-sdk/version.Version=${VERSION} \
              -X github.com/cosmos/cosmos-sdk/version.Commit=${COMMIT}" \
    -o ./build -tags "ledger" ./cmd/cudos-noded

FROM ubuntu:24.04
WORKDIR /app

RUN apt-get update && apt-get install -y jq moreutils yq

COPY --from=build /app/build/cudos-noded /app/
COPY ./tests/docker/libwasmvm.so /app/
COPY ./tests/docker/init.sh /app/

RUN ls

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
