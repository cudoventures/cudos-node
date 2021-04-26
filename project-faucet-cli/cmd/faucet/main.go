package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/tendermint/starport/starport/pkg/chaincmd"
	chaincmdrunner "github.com/tendermint/starport/starport/pkg/chaincmd/runner"
	"github.com/tendermint/starport/starport/pkg/cosmoscoin"
	"github.com/tendermint/starport/starport/pkg/cosmosfaucet"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	"github.com/tendermint/starport/starport/pkg/xhttp"
)

func main() {
	flag.Parse()

	configKeyringBackend, err := chaincmd.KeyringBackendFromString(keyringBackend)
	if err != nil {
		log.Fatal(err)
	}

	ccoptions := []chaincmd.Option{
		chaincmd.WithKeyringPassword(keyringPassword),
		chaincmd.WithKeyringBackend(configKeyringBackend),
		chaincmd.WithAutoChainIDDetection(),
		chaincmd.WithNodeAddress(nodeAddress),
	}

	if legacySendCmd {
		ccoptions = append(ccoptions, chaincmd.WithLegacySendCommand())
	}

	if sdkVersion == string(cosmosver.Stargate) {
		ccoptions = append(ccoptions,
			chaincmd.WithVersion(cosmosver.StargateZeroFourtyAndAbove),
		)
	} else {
		ccoptions = append(ccoptions,
			chaincmd.WithVersion(cosmosver.LaunchpadAny),
			chaincmd.WithLaunchpadCLI(appCli),
		)
	}

	cr, err := chaincmdrunner.New(context.Background(), chaincmd.New(appCli, ccoptions...))
	if err != nil {
		log.Fatal(err)
	}

	coins := strings.Split(defaultDenoms, denomSeparator)

	faucetOptions := make([]cosmosfaucet.Option, len(coins))
	for i, coin := range coins {
		faucetOptions[i] = cosmosfaucet.Coin(creditAmount, maxCredit, coin)
	}

	faucetOptions = append(faucetOptions, cosmosfaucet.Account(keyName, keyMnemonic))

	faucet, err := cosmosfaucet.New(context.Background(), cr, faucetOptions...)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		originHeader := r.Header.Get("Origin")
		if originHeader == "http://localhost:3000" || originHeader == "http://35.238.210.147:3000" {
			w.Header().Set("Access-Control-Allow-Origin", originHeader)
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			if (*r).Method == "OPTIONS" {
				return
			}
		}

		buf, _ := ioutil.ReadAll(r.Body)
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))

		var req cosmosfaucet.TransferRequest
		err := json.NewDecoder(rdr1).Decode(&req)
		if err == nil {
			coin := req.Coins[0]
			amount, denom, err := cosmoscoin.Parse(coin)
			if err == nil {
				if amount > maxCredit {
					var transfers []cosmosfaucet.Transfer
					t := cosmosfaucet.Transfer{
						Coin:   denom,
						Status: "error",
						Error:  fmt.Sprintf("max credit (%d)", maxCredit),
					}

					transfers = append(transfers, t)

					xhttp.ResponseJSON(w, http.StatusOK, cosmosfaucet.TransferResponse{
						Transfers: transfers,
					})
					return
				}
			}
		}

		r.Body = rdr2
		faucet.ServeHTTP(w, r)
	})
	log.Infof("listening on :%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
