export DENOMS="ucudos"
export PORT="5000"
export CAPTCHA_BACKEND="6Ldjs9EaAAAAAO6GmEKsDB1teUDi0AwRJNKfA7oF"
export ACCOUNT_NAME="faucet"
export CLI_NAME="cudos-noded"
export KEYRING_BACKEND="file"
export KEYRING_PASSWORD="123123123" #at least 8 characters
export CREDIT_AMOUNT="1000"
export MAX_CREDIT="100000000"
export NODE="http://10.128.0.2:26657"
export MNEMONIC="thumb web nasty tennis tenant leg exact exit found lawn wool swap tiger tobacco coach pelican under around delay call able attract desk million"

#remove old data before launch
rm -R cudos-data

faucet