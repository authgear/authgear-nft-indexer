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

func ConfigureGetCollectionRoute(route httproute.Route) httproute.Route {
	return route.
		WithMethods("GET").
		WithPathPattern("/collection/:contract_id")
}

type GetCollectionHandlerLogger struct{ *log.Logger }

func NewGetCollectionHandlerLogger(lf *log.Factory) GetCollectionHandlerLogger {
	return GetCollectionHandlerLogger{lf.New("api-get-collection")}
}

type GetCollectionAPIHandler struct {
	JSON               JSONResponseWriter
	Logger             ListCollectionHandlerLogger
	NFTCollectionQuery query.NFTCollectionQuery
}

func (h *GetCollectionAPIHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	contractIDStr := httproute.GetParam(req, "contract_id")

	contractID, err := model.ParseContractID(contractIDStr)
	if err != nil {
		h.Logger.WithError(err).Error("failed to parse contract URL")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("invalid contract URL")})
		return
	}

	qb := h.NFTCollectionQuery.NewQueryBuilder()

	qb.WithContracts([]model.ContractID{*contractID})

	collections, err := h.NFTCollectionQuery.ExecuteQuery(qb)
	if err != nil {
		h.Logger.WithError(err).Error("failed to list nft collections")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewInternalError("failed to list nft collections")})
		return
	}

	if len(collections) == 0 {
		h.Logger.Error("no contract found for the given contract ID")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewNotFound("no contract found for the given contract ID")})
		return
	}

	collection := collections[0]

	h.JSON.WriteResponse(resp, &authgearapi.Response{
		Result: &apimodel.NFTCollection{
			ID:              collection.ID,
			Blockchain:      collection.Blockchain,
			Network:         collection.Network,
			Name:            collection.Name,
			ContractAddress: collection.ContractAddress,
			BlockHeight:     *collection.FromBlockHeight.ToMathBig(),
			TotalSupply:     *collection.TotalSupply.ToMathBig(),
			Type:            string(collection.Type),
		},
	})
}
