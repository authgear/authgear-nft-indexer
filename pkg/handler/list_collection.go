package handler

import (
	"net/http"

	apimodel "github.com/authgear/authgear-nft-indexer/pkg/api/model"
	"github.com/authgear/authgear-nft-indexer/pkg/query"
	authgearapi "github.com/authgear/authgear-server/pkg/api"
	"github.com/authgear/authgear-server/pkg/api/apierrors"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"github.com/authgear/authgear-server/pkg/util/log"
	authgearweb3 "github.com/authgear/authgear-server/pkg/util/web3"
)

func ConfigureListCollectionRoute(route httproute.Route) httproute.Route {
	return route.
		WithMethods("GET").
		WithPathPattern("/collections")
}

type ListCollectionHandlerLogger struct{ *log.Logger }

func NewListCollectionHandlerLogger(lf *log.Factory) ListCollectionHandlerLogger {
	return ListCollectionHandlerLogger{lf.New("api-list-collection")}
}

type ListCollectionAPIHandler struct {
	JSON               JSONResponseWriter
	Logger             ListCollectionHandlerLogger
	NFTCollectionQuery query.NFTCollectionQuery
}

func (h *ListCollectionAPIHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	urlValues := req.URL.Query()

	contracts := make([]authgearweb3.ContractID, 0)
	for _, url := range urlValues["contract_id"] {
		e, err := authgearweb3.ParseContractID(url)
		if err != nil {
			h.Logger.WithError(err).Error("failed to parse contract URL")
			h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("invalid contract URL")})
			return
		}

		contracts = append(contracts, *e)
	}

	qb := h.NFTCollectionQuery.NewQueryBuilder()

	qb = qb.WithContracts(contracts)

	collections, err := h.NFTCollectionQuery.ExecuteQuery(qb)
	if err != nil {
		h.Logger.WithError(err).Error("failed to list nft collections")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewInternalError("failed to list nft collections")})
		return
	}

	nftCollections := make([]apimodel.NFTCollection, 0, len(collections))
	for _, collection := range collections {
		nftCollections = append(nftCollections, collection.ToAPIModel())
	}

	h.JSON.WriteResponse(resp, &authgearapi.Response{
		Result: &apimodel.CollectionListResponse{
			Items: nftCollections,
		},
	})
}
