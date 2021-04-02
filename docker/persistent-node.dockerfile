FROM golang:1.15

WORKDIR /usr/blockchain

COPY ./ ./

RUN make

RUN cudos-poc-01d init cudos-poc-01-network --chain-id=cudos-poc-01-network

RUN cudos-poc-01d keys add validator01 --keyring-backend test

RUN cudos-poc-01d keys show validator01 -a --keyring-backend test

RUN MY_VALIDATOR_ADDRESS=$(cudos-poc-01d keys show validator01 -a --keyring-backend test)&& cudos-poc-01d add-genesis-account $MY_VALIDATOR_ADDRESS 100000000000stake

RUN cudos-poc-01d gentx validator01 100000000stake --chain-id cudos-poc-01-network --keyring-backend test

RUN cudos-poc-01d collect-gentxs

CMD ["cudos-poc-01d", "start"] 