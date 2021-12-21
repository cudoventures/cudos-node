chain_id=$1 #MyLocalCudosNetwok
keyring_password=$2 #123123123
containerId=$3 #66268775226f
faucetAddress=$4 #cudos12kx3a8xva3hwhkqrfnw8n5m9ffmg2s9rtr272m b4f973a3c21d
initialDirectory=$PWD
#nodeUrl=$5
#port=$6
# ./simple_test.sh MyLocalCudosNetwok 123123123 66268775226f cudos12kx3a8xva3hwhkqrfnw8n5m9ffmg2s9rtr272m

# DOWNLOAD
alias CUDOS_NODED='docker exec -i $containerId cudos-noded'
alias CUDOS_NODED='docker exec -i $containerId cudos-noded'
git clone https://github.com/CudoVentures/cosmwasm-examples
cd cosmwasm-examples
git checkout
cd contracts/escrow/
## COMPILE
docker run --rm -v "$(pwd)":/code \
  --mount type=volume,source="$(basename "$(pwd)")_cache",target=/code/target \
  --mount type=volume,source=registry_cache,target=/usr/local/cargo/registry \
  cosmwasm/rust-optimizer:0.11.5
# COPY
cd artifacts/
docker cp cw_escrow.wasm $containerId:/usr/cudos
## ADD USERS
echo "$keyring_password" | CUDOS_NODED keys add wasm-poweruser --keyring-backend os  # handles prompt for keyring password
echo "$keyring_password" | CUDOS_NODED keys add wasm-receiver --keyring-backend os # handles prompt for keyring password
echo "$keyring_password" | CUDOS_NODED keys add wasm-thief --keyring-backend os # handles prompt for keyring password
# DEFINE ADDRESSES
wasmPowerUserAddress=$(echo "$keyring_password" | CUDOS_NODED keys show -a wasm-poweruser --keyring-backend os)
wasmReceiverUserAddress=$(echo "$keyring_password" | CUDOS_NODED keys show -a wasm-thief --keyring-backend os)
wasmThiefUserAddress=$(echo "$keyring_password" | CUDOS_NODED keys show -a wasm-receiver --keyring-backend os)
#faucetAddress='cudos1xqqscnu0pejkmkc36fc4ya5egpuwxevz868e4r'
# FUND USERS
echo "$keyring_password" | CUDOS_NODED tx bank send $faucetAddress "$wasmPowerUserAddress" 100000acudos --keyring-backend os --chain-id="$chain_id" -y
echo "$keyring_password" | CUDOS_NODED tx bank send $faucetAddress "$wasmReceiverUserAddress" 1111acudos --keyring-backend os --chain-id="$chain_id" -y
echo "$keyring_password" | CUDOS_NODED tx bank send $faucetAddress "$wasmThiefUserAddress" 999acudos --keyring-backend os --chain-id="$chain_id" -y
# STORE IN THE NODE
RES=$(echo "$keyring_password" |  CUDOS_NODED tx wasm store cw_escrow.wasm --from wasm-poweruser --gas auto --gas-adjustment 1.3 --keyring-backend os --chain-id="$chain_id" -y)
echo "Store TX Result: $RES"
# ASSERT SUCCESSFUL STORING
CODE_ID=$(echo "$RES" | jq -r '.logs[0].events[-1].attributes[-1].value')
echo "CODE_ID: $CODE_ID" # CODE_ID value must be a positive integer
if [ "$CODE_ID" -lt 1 ]; then # CODE_ID value must be a positive integer
    printf '%s\n' "Failed to store the smart contract on the chain. Check result below for a detailed error info" >&2 # write error message to stderr
    printf '%s\n' "Failed STORE Transaction: ""$RES"" " >&2 # write error message to stderr
fi

#INIT THE CONTRACT
INIT=$(jq -n --arg sender "$wasmPowerUserAddress" --arg receiver "$wasmReceiverUserAddress" '{"arbiter":$sender,"recipient":$receiver}')
RES=$(echo "$keyring_password" | CUDOS_NODED tx wasm instantiate "$CODE_ID" "$INIT" \
    --from wasm-poweruser --amount=5000acudos  --label "example escrow" --gas auto --gas-adjustment 1.3 --chain-id="$chain_id" -y --keyring-backend os)
echo "Init TX Result: $RES"


CUDOS_NODED query wasm list-contract-by-code "$CODE_ID"
CONTRACT=$(CUDOS_NODED query wasm list-contract-by-code "$CODE_ID" --output json | jq -r '.contracts[0]')
echo "$CONTRACT"

if [ "$CONTRACT" == "" ];then
    printf '%s\n' "Failed to init the smart contract on the chain. Check the tx result above for a detailed error info" >&2 # write error message to stderr
    printf '%s\n' "Failed INIT Transaction: ""$RES"" " >&2 # write error message to stderr

    exit 1
fi

contractBalanceBank=$(CUDOS_NODED query bank balances  "$CONTRACT")
# ASSERT the contract has balance of 5000acudos
if [[ ! $contractBalanceBank =~ "5000"  ]];then
    printf '%s\n' "Something went wrong, the balance of the contract should 5000, check if it was deployed and initialized successfully" >&2 # write error message to stderr
    exit 1
fi

# INTERACT
APPROVE='{"approve":{"quantity":[{"amount":"4700","denom":"acudos"}]}}'
RES=$(echo "$keyring_password" | CUDOS_NODED tx wasm execute "$CONTRACT" "$APPROVE" --from wasm-thief --gas auto --gas-adjustment 1.3 --chain-id="$chain_id" -y --keyring-backend os)

# RES SHOULD BE EMPTY - the TX should fail
receiverAddressBalance=$(CUDOS_NODED query bank balances "$wasmReceiverUserAddress")
if [[ ! $receiverAddressBalance =~ "1111"  ]];then
    printf '%s\n' "Something went wrong, the balance of $wasmReceiverUserAddress should be 1111, the transfer should have failed because it is coming from an unauthorized address" >&2 # write error message to stderr
    printf '%s\n' "Wrong Transaction: $RES" >&2
    exit 1
fi

# but it should succeed when wasm-poweruser tries
RES=$(echo "$keyring_password" | CUDOS_NODED tx wasm execute "$CONTRACT" "$APPROVE" --from wasm-poweruser --gas auto --gas-adjustment 1.3 --chain-id="$chain_id" -y --keyring-backend os)
echo "Successful INTERACT TX Result: $RES"

receiverAddressBalance=$(CUDOS_NODED query bank balances "$wasmReceiverUserAddress")
if [[ ! $receiverAddressBalance =~ "5811"  ]];then
    printf '%s\n' "Something went wrong, the balance of $wasmReceiverUserAddress should be 5811, check above if the interaction was successful" >&2 # write error message to stderr
    printf '%s\n' "Failed Transaction: $RES" >&2
    exit 1
fi

contractBalanceBank=$(CUDOS_NODED query bank balances  "$CONTRACT")
if [[ ! $contractBalanceBank =~ "300"  ]];then
    printf '%s\n' "Something went wrong, the balance of the contract should 300, check if it was deployed and initialized successfully" >&2 # write error message to stderr
    exit 1
fi

printf '%s\n' "Test successful! CosmWasm Smart contract is compatible with the cudos-node VM!"

### CLEAN UP: Remove USERS
printf '%s\n' 'Cleaning up users...' >&2 # write error message to stderr
echo "$keyring_password" |  CUDOS_NODED keys delete wasm-poweruser --keyring-backend os  -y # handles prompt for keyring password
echo "$keyring_password" | CUDOS_NODED keys delete wasm-receiver --keyring-backend os -y # handles prompt for keyring password
echo "$keyring_password" | CUDOS_NODED keys delete wasm-thief --keyring-backend os -y # handles prompt for keyring password

# Remove files
printf '%s\n' "Cleaning up files, deleting cosmwasm-examples at $initialDirectory" >&2 # write error message to stderr
cd $initialDirectory
rm -rf cosmwasm-examples
