package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

func main() {
	fmt.Println(os.Args)
	if len(os.Args) != 2 {
		panic(fmt.Sprintf("Expected exactly 1 argument but recieved %d", len(os.Args)-1))
	}

	ethPubKeyString := os.Args[1]
	if ethPubKeyString[:2] == "0x" {
		ethPubKeyString = ethPubKeyString[2:]
	}

	fmt.Println(ethPubKeyString)
	ethPubKeyBytes, err := hex.DecodeString(ethPubKeyString)
	if err != nil {
		panic(err)
	}

	ethPubKey := secp256k1.PubKey{
		Key: ethPubKeyBytes,
	}
	accAddrString, err := sdk.AccAddressFromHex(ethPubKey.Address().String())
	if err != nil {
		panic(err)
	}
	addrString, err := bech32.ConvertAndEncode("cosmos", accAddrString)
	if err != nil {
		panic(err)
	}
	fmt.Println(addrString)
}
