export DENOMS="ucudos"
export PORT="5000"
export CAPTCHA_BACKEND="6Ldjs9EaAAAAAO6GmEKsDB1teUDi0AwRJNKfA7oF"
export ACCOUNT_NAME="faucet"
export CLI_NAME="cudos-noded"
export KEYRING_BACKEND="file"
export KEYRING_PASSWORD="123123123" #at least 8 characters
export CREDIT_AMOUNT="1000"
export MAX_CREDIT="100000000"
export NODE="http://localhost:26657"
export MNEMONIC="black select bar south name card input labor movie cluster try fantasy flip jazz north inherit chat paddle scatter fox chapter spot shield donkey"

#remove old data before launch
rm -R cudos-data

faucet