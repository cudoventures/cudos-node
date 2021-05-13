# Prerequirements
    SSH key for accessing Google Cloud Infrastructure
    Install meteor - https://www.meteor.com/developers/install
    NodeJs >= 12

# Deploy testnet
    npm run deploy-testnet

# Deploy and initialize test net
    npm run deploy-and-init-testnet


# Deploy utils
1. Copy genesis time to settings.json in explorer
2. Copy node id to notion
3. Download faucet.wallet from cudos-root-node
4. Copy faucet mnemonic in /project-faucet-cli/run-testnet.sh
5. Download cudos-noded binary from cudos-root-node and copy it to /project-faucet-cli/bin/cudos-noded
6. Upload binary and genesis to notion

        npm run deploy-utils


