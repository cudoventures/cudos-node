package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"cudos.org/cudos-node/x/nft/types"
)

const (
	QueryGlobalRoutePrefix  = "cudosnode.cudosnode.nft.Query"
	QuerySupplyRoute        = "Supply"
	QueryOwnerRoute         = "Owner"
	QueryCollectionRoute    = "Collection"
	QueryDenomsRoute        = "Denoms"
	QueryDenomRoute         = "Denom"
	QueryDenomByNameRoute   = "DenomByName"
	QueryDenomBySymbolRoute = "DenomBySymbol"
	QueryNFTRoute           = "NFT"
	QueryApprovalsNFTRoute  = "GetApprovalsNFT"
	QueryIsApprovedForAll   = "QueryApprovalsIsApprovedForAll"
)

func registerQueryRoutes(cliCtx client.Context, r *mux.Router) {
	// Query the denom
	r.HandleFunc(fmt.Sprintf("/%s/denoms/{%s}", types.ModuleName, RestParamDenomID), queryDenom(cliCtx)).Methods("GET")

	// Query the denom by name
	r.HandleFunc(fmt.Sprintf("/%s/denoms/name/{%s}", types.ModuleName, RestParamDenomName), queryDenomByName(cliCtx)).Methods("GET")

	// Query the denom by symbol
	r.HandleFunc(fmt.Sprintf("/%s/denoms/symbol/{%s}", types.ModuleName, RestParamDenomSymbol), queryDenoBySymbol(cliCtx)).Methods("GET")

	// Query all denoms
	r.HandleFunc(fmt.Sprintf("/%s/denoms", types.ModuleName), queryDenoms(cliCtx)).Methods("POST")

	// Get all the NFTs from a given collection
	r.HandleFunc(fmt.Sprintf("/%s/collections", types.ModuleName), queryCollection(cliCtx)).Methods("POST")

	// Get the total supply of a collection or owner
	r.HandleFunc(fmt.Sprintf("/%s/collections/supply/{%s}", types.ModuleName, RestParamDenomID), querySupply(cliCtx)).Methods("GET")

	// Get the collections of NFTs owned by an address
	r.HandleFunc(fmt.Sprintf("/%s/owners", types.ModuleName), queryOwner(cliCtx)).Methods("POST")

	// Query a single NFT
	r.HandleFunc(fmt.Sprintf("/%s/nfts/{%s}/{%s}", types.ModuleName, RestParamDenomID, RestParamTokenID), queryNFT(cliCtx)).Methods("GET")

	// Query approvals for NFT
	r.HandleFunc(fmt.Sprintf("/%s/approvals/{%s}/{%s}", types.ModuleName, RestParamDenomID, RestParamTokenID), queryApprovalsNFT(cliCtx)).Methods("GET")

	// Query is approved for all
	r.HandleFunc(fmt.Sprintf("/%s/isApprovedForAll", types.ModuleName), queryIsApprovedForAll(cliCtx)).Methods("POST")
}

func querySupply(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		denomID := mux.Vars(r)[RestParamDenomID]
		err := types.ValidateDenomID(denomID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}

		var _ sdk.AccAddress
		ownerStr := r.FormValue(RestParamOwner)
		if len(ownerStr) > 0 {
			_, err = sdk.AccAddressFromBech32(ownerStr)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		request := types.QuerySupplyRequest{
			DenomId: denomID,
			Owner:   ownerStr,
		}

		bz, err := request.Marshal()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		queryPath := fmt.Sprintf("/%s/%s", QueryGlobalRoutePrefix, QuerySupplyRoute)
		res, height, err := cliCtx.QueryWithData(
			queryPath, bz,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var querySupplyResponse types.QuerySupplyResponse
		cliCtx.Codec.MustUnmarshal(res, &querySupplyResponse)

		// nolint: govet
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, querySupplyResponse)
	}
}

func queryOwner(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req queryOwnerRequest
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		if err := types.ValidateDenomID(req.DenomId); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}

		var _ sdk.AccAddress
		if len(req.OwnerAddress) > 0 {
			_, err := sdk.AccAddressFromBech32(req.OwnerAddress)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}
		request := types.QueryOwnerRequest{
			DenomId:    req.DenomId,
			Owner:      req.OwnerAddress,
			Pagination: &req.Pagination,
		}

		bz, err := request.Marshal()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		queryPath := fmt.Sprintf("/%s/%s", QueryGlobalRoutePrefix, QueryOwnerRoute)
		res, height, err := cliCtx.QueryWithData(
			queryPath, bz,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var ownerResponse types.QueryOwnerResponse
		cliCtx.Codec.MustUnmarshal(res, &ownerResponse)

		// nolint: govet
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, ownerResponse)
	}
}

func queryCollection(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req queryCollectionRequest
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		if err := types.ValidateDenomID(req.DenomId); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}

		request := types.QueryCollectionRequest{
			DenomId:    req.DenomId,
			Pagination: &req.Pagination,
		}

		bz, err := request.Marshal()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// nolint: govet
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		queryPath := fmt.Sprintf("/%s/%s", QueryGlobalRoutePrefix, QueryCollectionRoute)
		res, height, err := cliCtx.QueryWithData(
			queryPath, bz,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		var collectionResponse types.QueryCollectionResponse
		cliCtx.Codec.MustUnmarshal(res, &collectionResponse)

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, collectionResponse)
	}
}

// nolint: dupl
func queryDenom(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// nolint: govet
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		denomID := mux.Vars(r)[RestParamDenomID]
		if err := types.ValidateDenomID(denomID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		request := types.QueryDenomRequest{DenomId: denomID}
		bz, err := request.Marshal()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		queryPath := fmt.Sprintf("/%s/%s", QueryGlobalRoutePrefix, QueryDenomRoute)
		res, height, err := cliCtx.QueryWithData(
			queryPath, bz,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var denomResponse types.QueryDenomResponse
		cliCtx.Codec.MustUnmarshal(res, &denomResponse)

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, denomResponse)
	}
}

// nolint: dupl
func queryDenomByName(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// nolint: govet
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		denomName := mux.Vars(r)[RestParamDenomName]
		if err := types.ValidateDenomName(denomName); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}

		request := types.QueryDenomByNameRequest{DenomName: denomName}
		bz, err := request.Marshal()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		queryPath := fmt.Sprintf("/%s/%s", QueryGlobalRoutePrefix, QueryDenomByNameRoute)
		res, height, err := cliCtx.QueryWithData(
			queryPath, bz,
		)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var denomByNameResponse types.QueryDenomByNameResponse
		cliCtx.Codec.MustUnmarshal(res, &denomByNameResponse)

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, denomByNameResponse)
	}
}

// nolint: dupl
func queryDenoBySymbol(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// nolint: govet
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		denomSymbol := mux.Vars(r)[RestParamDenomSymbol]
		if err := types.ValidateDenomSymbol(denomSymbol); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}

		request := types.QueryDenomBySymbolRequest{Symbol: denomSymbol}
		bz, err := request.Marshal()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		queryPath := fmt.Sprintf("/%s/%s", QueryGlobalRoutePrefix, QueryDenomBySymbolRoute)
		res, height, err := cliCtx.QueryWithData(
			queryPath, bz,
		)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var denomBySymbolResponse types.QueryDenomBySymbolResponse
		cliCtx.Codec.MustUnmarshal(res, &denomBySymbolResponse)

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, denomBySymbolResponse)
	}
}

func queryDenoms(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// nolint: govet
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		var req queryDenomsRequest
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		request := types.QueryDenomsRequest{
			Pagination: &req.Pagination,
		}

		bz, err := request.Marshal()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		queryPath := fmt.Sprintf("/%s/%s", QueryGlobalRoutePrefix, QueryDenomsRoute)
		res, height, err := cliCtx.QueryWithData(
			queryPath, bz,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var denomsResponse types.QueryDenomsResponse
		cliCtx.Codec.MustUnmarshal(res, &denomsResponse)

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, denomsResponse)
	}
}

func queryNFT(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		denomID := vars[RestParamDenomID]
		if err := types.ValidateDenomID(denomID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tokenID := vars[RestParamTokenID]
		if err := types.ValidateTokenID(tokenID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		request := types.QueryNFTRequest{
			DenomId: denomID,
			TokenId: tokenID,
		}
		bz, err := request.Marshal()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// nolint: govet
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		queryPath := fmt.Sprintf("/%s/%s", QueryGlobalRoutePrefix, QueryNFTRoute)
		res, height, err := cliCtx.QueryWithData(
			queryPath, bz,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var nftResponse types.QueryNFTResponse
		cliCtx.Codec.MustUnmarshal(res, &nftResponse)

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, nftResponse)
	}
}

func queryApprovalsNFT(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		denomID := vars[RestParamDenomID]
		if err := types.ValidateDenomID(denomID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}

		tokenID := vars[RestParamTokenID]
		if err := types.ValidateTokenID(tokenID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}

		request := types.QueryApprovalsNFTRequest{
			DenomId: denomID,
			TokenId: tokenID,
		}

		bz, err := request.Marshal()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// nolint: govet
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		queryPath := fmt.Sprintf("/%s/%s", QueryGlobalRoutePrefix, QueryApprovalsNFTRoute)
		res, height, err := cliCtx.QueryWithData(
			queryPath, bz,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var approvalsNFTResponse types.QueryApprovalsNFTResponse
		cliCtx.Codec.MustUnmarshal(res, &approvalsNFTResponse)

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, approvalsNFTResponse)
	}
}

func queryIsApprovedForAll(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req queryIsApprovedForAllRequest
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		request := types.QueryApprovalsIsApprovedForAllRequest{
			Owner:    req.Owner,
			Operator: req.Operator,
		}
		bz, err := request.Marshal()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// nolint: govet
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		queryPath := fmt.Sprintf("/%s/%s", QueryGlobalRoutePrefix, QueryIsApprovedForAll)
		res, height, err := cliCtx.QueryWithData(
			queryPath, bz,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var isApprovedForAll types.QueryApprovalsIsApprovedForAllResponse
		cliCtx.Codec.MustUnmarshal(res, &isApprovedForAll)

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, isApprovedForAll)
	}
}
