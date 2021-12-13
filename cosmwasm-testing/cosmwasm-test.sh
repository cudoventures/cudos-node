# TODO: add url,port, chain-id, keyring pass, containerId, faucet address as vars

# DOWNLOAD
alias CUDOS_NODED='docker exec -i d2e6215eaf2f cudos-noded'
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
docker cp cw_escrow.wasm d2e6215eaf2f:/usr/cudos
## ADD USERS
echo "123123123" | CUDOS_NODED keys add wasm-poweruser8 --keyring-backend os  # handles prompt for keyring password
echo "123123123" | CUDOS_NODED keys add wasm-receiver8 --keyring-backend os # handles prompt for keyring password
echo "123123123" | CUDOS_NODED keys add wasm-thief8 --keyring-backend os # handles prompt for keyring password
# DEFINE ADDRESSES
wasmPowerUserAddress=$(echo "123123123" | CUDOS_NODED keys show -a wasm-poweruser8 --keyring-backend os)
wasmReceiverUserAddress=$(echo "123123123" | CUDOS_NODED keys show -a wasm-thief8 --keyring-backend os)
wasmThiefUserAddress=$(echo "123123123" | CUDOS_NODED keys show -a wasm-receiver8 --keyring-backend os)
faucetAddress='cudos1wwt7ugpnckrf44sv8yvzk42fdgefgr37a44dpa'
# FUND USERS
echo "123123123" | CUDOS_NODED tx bank send $faucetAddress "$wasmPowerUserAddress" 10000000000000acudos --keyring-backend os --chain-id=MyLocalCudosNetwok -y
echo "123123123" | CUDOS_NODED tx bank send $faucetAddress "$wasmReceiverUserAddress" 100000000000000acudos --keyring-backend os --chain-id=MyLocalCudosNetwok -y
echo "123123123" | CUDOS_NODED tx bank send $faucetAddress "$wasmThiefUserAddress" 100000000000000acudos --keyring-backend os --chain-id=MyLocalCudosNetwok -y
echo "$wasmPowerUserAddress"
contractBalanceBank=$(CUDOS_NODED query bank balances "$wasmPowerUserAddress")
echo "$contractBalanceBank"
# STORE IN THE NODE
RES=$(echo "123123123" |  CUDOS_NODED tx wasm store cw_escrow.wasm --from wasm-poweruser8 --gas auto --gas-adjustment 1.3 --keyring-backend os --chain-id MyLocalCudosNetwok -y)
echo "Store TX Result: $RES"
# ASSERT SUCCESSFULL DEPLOYMENT
CODE_ID=$(echo "$RES" | jq -r '.logs[0].events[-1].attributes[-1].value')
echo "CODE_ID: $CODE_ID" # CODE_ID value must be a positive integer

if [ "$CODE_ID" -lt 1 ]; then # CODE_ID value must be a positive integer
    printf '%s\n' "Failed to store the smart contract on the chain. Check result below for a detailed error info" >&2 # write error message to stderr
    printf '%s\n' "Failed Transaction: ""$RES"" " >&2 # write error message to stderr
fi

# instantiate contract and verify
INIT=$(jq -n --arg sender $(CUDOS_NODED keys show -a wasm-poweruser) --arg receiver $(CUDOS_NODED keys show -a wasm-receiver) '{"arbiter":$sender,"recipient":$receiver}')
CUDOS_NODED tx wasm instantiate "$CODE_ID" "$INIT" \
    --from wasm-poweruser --amount=5000acudos  --label "example escrow" --gas auto --gas-adjustment 1.3 --chain-id MyLocalCudosNetwok -y


CUDOS_NODED query wasm list-contract-by-code "$CODE_ID"
CONTRACT=$(CUDOS_NODED query wasm list-contract-by-code "$CODE_ID" --output json | jq -r '.contracts[0]')
echo "$CONTRACT"

if [ "$CONTRACT" == "" ];then
    printf '%s\n' "Failed to init the smart contract on the chain. Check the tx result above for a detailed error info" >&2 # write error message to stderr
    exit 1
fi

contractBalanceWasm=$(CUDOS_NODED query wasm contract "$CONTRACT")
contractBalanceBank=$(CUDOS_NODED query bank balances  "$CONTRACT")
# assert the contract has balance of 5000acudos
if [ "$contractBalanceWasm" != "5000acudos" ] || [ "$contractBalanceBank" != "5000acudos" ];then
    printf '%s\n' "Something went wrong, the balance of the contract should 5000, check if it was deployed and initialized successfully" >&2 # write error message to stderr
    exit 1
fi

# execute fails if wrong person
APPROVE='{"approve":{"quantity":[{"amount":"5000","denom":"acudos"}]}}'
CUDOS_NODED tx wasm execute $CONTRACT "$APPROVE" --from wasm-thief --gas auto --gas-adjustment 1.3 --chain-id MyLocalCudosNetwok -y

receiverAddress=$(CUDOS_NODED keys show wasm-receiver -a)
receiverAddressBalance=$(CUDOS_NODED query bank balances "$receiverAddress")
if [ "$receiverAddressBalance" != 0 ];then
    printf '%s\n' "Something went wrong, the balance of $receiverAddress should be 0, check above if the interaction was successful" >&2 # write error message to stderr
    exit 1
fi


# but it should succeed when wasm-poweruser tries
CUDOS_NODED tx wasm execute $CONTRACT "$APPROVE" --from wasm-poweruser --gas auto --gas-adjustment 1.3 --chain-id cudos-testnet-public --node tcp://35.232.27.92:26657 -y\

# the receiver user should have 5000acudos in its balance and the contract should have no money in its balance
receiverAddressBalance=$(CUDOS_NODED query bank balances "$receiverAddress")
if [ "$receiverAddressBalance" != 5000 ];then
    printf '%s\n' "Something went wrong, the balance of $receiverAddress should be 5000, check above if the interaction was successful" >&2 # write error message to stderr
    exit 1
fi

contractBalanceWasm=$(CUDOS_NODED query wasm contract "$CONTRACT")
contractBalanceBank=$(CUDOS_NODED query bank balances  "$CONTRACT")
# assert the contract has balance of 5000acudos
if [ "$contractBalanceWasm" != "0acudos" ] || [ "$contractBalanceBank" != "0acudos" ];then
    printf '%s\n' "Something went wrong, the balance of the contract should 0, check if it was deployed and initialized successfully" >&2 # write error message to stderr
    exit 1
fi