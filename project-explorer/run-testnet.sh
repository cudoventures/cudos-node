export MONGO_URL=mongodb://root:cudos-root-db-pass@localhost:27017
export ROOT_URL=http://35.238.210.147
export PORT=3000
export METEOR_SETTINGS='{"public":{"chainName": "CudosNetwork","chainId": "cudos-network","gtm": "","slashingWindow": 10000,"uptimeWindow": 250,"initialPageSize": 30,"secp256k1": true,"bech32PrefixAccAddr": "cudos","bech32PrefixAccPub": "cudospub","bech32PrefixValAddr": "cudosvaloper","bech32PrefixValPub": "cudosvaloperpub","bech32PrefixConsAddr": "cudosvalcons","bech32PrefixConsPub": "cudosvalconspub","bondDenom": "cudos","powerReduction": 1000000,"genesisTime": "2021-04-65T09:15:0.268797802Z","faucetUrl": "http://35.238.210.147:5000","coins": [{"denom": "cudos","displayName": "CUDOS","fraction": 1000000} ],"ledger":{"coinType": 118,"appName": "Cudos","appVersion": "0.0.1","gasPrice": 0.02},"modules": {"bank": true,"supply": true,"minting": true,"gov": true,"distribution": true},"coingeckoId": "cudos"},"remote":{"rpc":"http://localhost:26657","api":"http://localhost:1317"},"debug": {"startTimer": true},"params":{"startHeight": 0,"defaultBlockTime": 5000,"validatorUpdateWindow": 300,"blockInterval": 15000,"transactionsInterval": 18000,"keybaseFetchingInterval": 18000000,"consensusInterval": 1000,"statusInterval":7500,"signingInfoInterval": 1800000,"proposalInterval": 5000,"missedBlocksInterval": 60000,"delegationInterval": 900000}}'

node main.js