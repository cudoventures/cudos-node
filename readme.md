# Cudos Node

This project contains cudo-node source code

# Converting ethereum public keys to cosmos wallet address
Run the converter and pass a ethereum public key as argument.

    go run ./converter 0x03139bb3b92e99d034ee38674a0e29c4aad83dd09b3fa465a265da310f9948fbe6

<b>Example ethereum mnemonic:</b> battle erosion opinion city birth modify scale hood caught menu risk rather<br >
<b>Example ethereum public key (32 bytes, compressed form):</b> 0x03139bb3b92e99d034ee38674a0e29c4aad83dd09b3fa465a265da310f9948fbe6

This mnemonic could be imported into cudos blockchain in order to verify that resulting account access will be the same as generated from the converter.

    cudos-noded keys add ruser02 --recover --hd-path="m/44'/60'/0'/0/0"

# Usefull commands
## Send currency
    cudos-noded tx bank send $VALIDATOR_ADDRESS $RECIPIENT 51000000stake --chain-id=cudos-network --keyring-backend test

## Check balances
    cudos-noded query bank balances $RECIPIENT --chain-id=cudos-network

## Create validator
    cudos-noded tx staking create-validator --amount=2000000000000000000000000acudos \
    --from=val-2 \
    --pubkey=$(cudos-noded tendermint show-validator) \
    --moniker=cudos-node-02 \
    --chain-id=cudos-local-network \
    --commission-rate="0.10" \
    --commission-max-rate="0.20" \
    --commission-max-change-rate="0.01" \
    --min-self-delegation="2000000000000000000000000" \
    --gas="auto" \
    --gas-prices="5000000000000acudos" \
    --gas-adjustment="1.80" \
    --keyring-backend test

<br />
<br />
<br />

# Reset the blockchain
All data of the blockchain is store at ~/.cudos folder. By deleting it the entire blockchain is completely reset and it must be initialized again.

# Static compile
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-extldflags "-static"' ./cmd/cudos-noded/

export CGO_LDFLAGS="-lpthread -ldl"
go build -v -a -tags netgo,osusergo -ldflags='-lpthread -extldflags "-lpthread -static"' ./cmd/cudos-noded/

