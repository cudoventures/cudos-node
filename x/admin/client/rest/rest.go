package rest

import (
	"github.com/CudoVentures/cudos-node/x/admin/types"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	// this line is used by starport scaffolding # 1
)

// RegisterRoutes registers admin-related REST handlers to a router
func RegisterRoutes(clientCtx client.Context, r *mux.Router) {
	registerTxHandlers(clientCtx, r)
	registerQueryRoutes(clientCtx, r)
}

func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
}

func registerTxHandlers(clientCtx client.Context, r *mux.Router) {
	r.HandleFunc("/admin/spend", adminSpendCommunityPool(clientCtx)).Methods("POST")
}

func adminSpendCommunityPool(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req adminSpendCommunityPoolRequest
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(baseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request - wrong from address")
			return
		}
		toAddr, err := sdk.AccAddressFromBech32(req.ToAddress)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request - wrong from to address")
			return
		}
		coins, err := sdk.ParseCoinsNormalized(req.Amount)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request - wrong coins")
			return
		}

		// create the message
		msg := types.NewMsgAdminSpendCommunityPool(
			fromAddr,
			toAddr,
			coins,
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

type adminSpendCommunityPoolRequest struct {
	BaseReq   rest.BaseReq `json:"base_req"`
	ToAddress string       `json:"to_address"`
	Amount    string       `json:"amount"`
}
