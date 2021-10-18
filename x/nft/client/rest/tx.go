package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"cudos.org/cudos-node/x/nft/types"
)

func registerTxRoutes(cliCtx client.Context, r *mux.Router, queryRoute string) {
	// Issue a denom
	r.HandleFunc("/nft/nfts/denoms/issue", issueDenomHandlerFn(cliCtx)).Methods("POST")
	// Mint an NFT
	r.HandleFunc("/nft/nfts/mint", mintNFTHandlerFn(cliCtx)).Methods("POST")
	// Update an NFT
	r.HandleFunc(fmt.Sprintf("/nft/nfts/{%s}/{%s}", RestParamDenomID, RestParamTokenID), editNFTHandlerFn(cliCtx)).Methods("PUT")
	// Transfer an NFT to an address
	r.HandleFunc(fmt.Sprintf("/nft/nfts/{%s}/{%s}/transfer", RestParamDenomID, RestParamTokenID), transferNFTHandlerFn(cliCtx)).Methods("POST")
	// Send an NFT to an address
	r.HandleFunc(fmt.Sprintf("/nft/nfts/{%s}/{%s}/send", RestParamDenomID, RestParamTokenID, RestParamMessage), sendNFTHandlerFn(cliCtx)).Methods("POST")
	// Approve NFT transfers for address
	r.HandleFunc(fmt.Sprintf("/nft/nfts/{%s}/{%s}/approve", RestParamDenomID, RestParamTokenID, RestParamMessage), approveNFTHandlerFn(cliCtx)).Methods("POST")
	// Revoke NFT transfers for address
	r.HandleFunc(fmt.Sprintf("/nft/nfts/{%s}/{%s}/revoke", RestParamDenomID, RestParamTokenID, RestParamMessage), revokeNFTHandlerFn(cliCtx)).Methods("POST")
	// Burn an NFT
	r.HandleFunc(fmt.Sprintf("/nft/nfts/{%s}/{%s}/burn", RestParamDenomID, RestParamTokenID), burnNFTHandlerFn(cliCtx)).Methods("POST")
}

func issueDenomHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req issueDenomReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the message
		msg := types.NewMsgIssueDenom(req.ID, req.Name, req.Schema, req.Owner)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func mintNFTHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req mintNFTReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		if req.Recipient == "" {
			req.Recipient = req.Owner
		}
		// create the message
		msg := types.NewMsgMintNFT(
			req.ID,
			req.DenomID,
			req.Name,
			req.URI,
			req.Data,
			req.Owner,
			req.Recipient,
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func editNFTHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req editNFTReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		vars := mux.Vars(r)
		// create the message
		msg := types.NewMsgEditNFT(
			vars[RestParamTokenID],
			vars[RestParamDenomID],
			req.Name,
			req.URI,
			req.Data, req.Owner,
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func transferNFTHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req transferNFTReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}
		if _, err := sdk.AccAddressFromBech32(req.To); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		vars := mux.Vars(r)
		// create the message
		msg := types.NewMsgTransferNft(
			vars[RestParamTokenID],
			vars[RestParamDenomID],
			req.From,
			req.To,
			req.BaseReq.From)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func approveNFTHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req approveNFTReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}
		if _, err := sdk.AccAddressFromBech32(req.ToAddress); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		vars := mux.Vars(r)
		tokenId := vars[RestParamTokenID]
		denomId := vars[RestParamDenomID]
		// create the message
		msg := types.NewMsgApproveNft(
			tokenId,
			denomId,
			req.Owner,
			req.ToAddress,
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func revokeNFTHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req revokeNFTReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}
		if _, err := sdk.AccAddressFromBech32(req.Recipient); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		vars := mux.Vars(r)
		// create the message
		msg := types.NewMsgRevokeNft(
			vars[RestParamTokenID],
			vars[RestParamDenomID],
			req.Owner,
			req.Recipient,
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func sendNFTHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req transferNFTReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}
		if _, err := sdk.AccAddressFromBech32(req.To); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		vars := mux.Vars(r)
		// create the message
		msg := types.NewMsgSendNft(
			vars[RestParamTokenID],
			vars[RestParamDenomID],
			req.From,
			req.To,
			vars[RestParamMessage],
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func burnNFTHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req burnNFTReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		vars := mux.Vars(r)

		// create the message
		msg := types.NewMsgBurnNFT(
			req.Owner,
			vars[RestParamTokenID],
			vars[RestParamDenomID],
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}
