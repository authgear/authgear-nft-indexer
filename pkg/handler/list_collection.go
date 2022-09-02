package handler

import (
	"net/http"

	apimodel "github.com/authgear/authgear-nft-indexer/pkg/api/model"
	"github.com/authgear/authgear-nft-indexer/pkg/model"
	"github.com/authgear/authgear-nft-indexer/pkg/query"
	authgearapi "github.com/authgear/authgear-server/pkg/api"
	"github.com/authgear/authgear-server/pkg/api/apierrors"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"github.com/authgear/authgear-server/pkg/util/log"
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

	contracts := make([]model.ContractID, 0)
	for _, url := range urlValues["contract_id"] {
		e, err := model.ParseContractID(url)
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
		nftCollections = append(nftCollections, apimodel.NFTCollection{
			ID:              collection.ID,
			Blockchain:      collection.Blockchain,
			Network:         collection.Network,
			Name:            collection.Name,
			BlockHeight:     *collection.FromBlockHeight.ToMathBig(),
			ContractAddress: collection.ContractAddress,
			TotalSupply:     *collection.TotalSupply.ToMathBig(),
			Type:            string(collection.Type),
		})
	}

	h.JSON.WriteResponse(resp, &authgearapi.Response{
		Result: &apimodel.CollectionListResponse{
			Items: nftCollections,
		},
	})
}
