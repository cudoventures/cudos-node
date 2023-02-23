package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/CudoVentures/cudos-node/x/nft/types"
)

func registerTxRoutes(cliCtx client.Context, r *mux.Router) {
	// Issue a denom
	r.HandleFunc("/nft/nfts/denoms/issue", issueDenomHandlerFn(cliCtx)).Methods("POST")
	// Mint an NFT
	r.HandleFunc("/nft/nfts/mint", mintNFTHandlerFn(cliCtx)).Methods("POST")
	// Update an NFT
	r.HandleFunc(fmt.Sprintf("/nft/nfts/edit/{%s}/{%s}", RestParamDenomID, RestParamTokenID), editNFTHandlerFn(cliCtx)).Methods("PUT")
	// Transfer an NFT to an address
	r.HandleFunc(fmt.Sprintf("/nft/nfts/transfer/{%s}/{%s}", RestParamDenomID, RestParamTokenID), transferNFTHandlerFn(cliCtx)).Methods("POST")
	// Transfer an NFT Collection to an address
	r.HandleFunc(fmt.Sprintf("/nft/nfts/denoms/transfer/{%s}", RestParamDenomID), transferDenomHandlerFn(cliCtx)).Methods("POST")
	// Approve NFT transfers for address
	r.HandleFunc(fmt.Sprintf("/nft/nfts/approve/{%s}/{%s}", RestParamDenomID, RestParamTokenID), approveNFTHandlerFn(cliCtx)).Methods("POST")
	// Revoke NFT transfers for address
	r.HandleFunc(fmt.Sprintf("/nft/nfts/revoke/{%s}/{%s}", RestParamDenomID, RestParamTokenID), revokeNFTHandlerFn(cliCtx)).Methods("POST")
	// Burn an NFT
	r.HandleFunc(fmt.Sprintf("/nft/nfts/burn/{%s}/{%s}", RestParamDenomID, RestParamTokenID), burnNFTHandlerFn(cliCtx)).Methods("POST")
	// Approve All for address
	r.HandleFunc(fmt.Sprintf("/nft/nfts/approveAll"), approveAll(cliCtx)).Methods("POST")
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
		msg := types.NewMsgIssueDenom(req.ID, req.Name, req.Schema, req.BaseReq.From, "", req.Symbol, req.Traits, req.Minter, req.Description, req.Data)
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
			rest.WriteErrorResponse(w, http.StatusBadRequest, "no recipient specified")
			return
		}
		// create the message
		msg := types.NewMsgMintNFT(
			req.DenomID,
			req.Name,
			req.URI,
			req.Data,
			req.BaseReq.From,
			req.Recipient,
			"",
		)
		if err2 := msg.ValidateBasic(); err2 != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err2.Error())
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
			req.Data,
			req.BaseReq.From,
			"",
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
			vars[RestParamDenomID],
			vars[RestParamTokenID],
			req.From,
			req.To,
			req.BaseReq.From,
			"")
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func transferDenomHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req transferDenomReq
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
		msg := types.NewMsgTransferDenom(
			vars[RestParamDenomID],
			req.BaseReq.From,
			req.Recipient,
			"",
		)
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
		if _, err := sdk.AccAddressFromBech32(req.AddressToApprove); err != nil {
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
			req.BaseReq.From,
			req.AddressToApprove,
			"",
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
		if _, err := sdk.AccAddressFromBech32(req.AddressToRevoke); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		vars := mux.Vars(r)
		// create the message
		msg := types.NewMsgRevokeNft(
			req.AddressToRevoke,
			req.BaseReq.From,
			vars[RestParamDenomID],
			vars[RestParamTokenID],
			"",
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
			req.BaseReq.From,
			vars[RestParamTokenID],
			vars[RestParamDenomID],
			"",
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func approveAll(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req approveAllRequest
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the message
		msg := types.NewMsgApproveAllNft(
			req.ApprovedOperator,
			req.BaseReq.From,
			"",
			req.Approved,
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}
