RMDIR .\data /s /q

SET MONIKER="cudos-node"

cudos-noded init %MONIKER% --chain-id=cudos-node

cudos-noded keys add validator01 --keyring-backend test

cudos-noded keys show validator01 -a > .\address.txt

SET /p VALIDATOR_ADDRESS=< ./address.txt

cudos-noded add-genesis-account %VALIDATOR_ADDRESS% 100000000000stake

del ./address.txt

cudos-noded gentx validator01 100000000stake --chain-id cudos-node --keyring-backend test

cudos-noded collect-gentxs

cudos-noded start