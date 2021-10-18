package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

// RegisterHandlers registers the NFT REST routes.
func RegisterHandlers(cliCtx client.Context, r *mux.Router, queryRoute string) {
	registerQueryRoutes(cliCtx, r, queryRoute)
	registerTxRoutes(cliCtx, r, queryRoute)
}

const (
	RestParamDenomID   = "denom-id"
	RestParamDenomName = "denom-name"
	RestParamTokenID   = "token-id"
	RestParamOwner     = "owner"
	RestParamMessage   = "msg"
)

type issueDenomReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Owner   string       `json:"owner"`
	ID      string       `json:"id"`
	Name    string       `json:"name"`
	Schema  string       `json:"schema"`
}

type mintNFTReq struct {
	BaseReq   rest.BaseReq `json:"base_req"`
	Owner     string       `json:"owner"`
	Recipient string       `json:"recipient"`
	DenomID   string       `json:"denom_id"`
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	URI       string       `json:"uri"`
	Data      string       `json:"data"`
}

type editNFTReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Owner   string       `json:"owner"`
	Name    string       `json:"name"`
	URI     string       `json:"uri"`
	Data    string       `json:"data"`
}

type transferNFTReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	From    string       `json:"from"`
	To      string       `json:"to"`
	Name    string       `json:"name"`
	URI     string       `json:"uri"`
	Data    string       `json:"data"`
}

type sendNFTReq struct {
	BaseReq   rest.BaseReq `json:"base_req"`
	Owner     string       `json:"owner"`
	Recipient string       `json:"recipient"`
	Name      string       `json:"name"`
	URI       string       `json:"uri"`
	Data      string       `json:"data"`
	Messsage  string       `json:"msg"`
}

type approveNFTReq struct {
	BaseReq   rest.BaseReq `json:"base_req"`
	Owner     string       `json:"owner"`
	ToAddress string       `json:"recipient"`
	Expires   string       `json:"expires"`
}

type revokeNFTReq struct {
	BaseReq   rest.BaseReq `json:"base_req"`
	Owner     string       `json:"owner"`
	Recipient string       `json:"recipient"`
}

type burnNFTReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Owner   string       `json:"owner"`
}
