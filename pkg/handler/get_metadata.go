package handler

import (
	"net/http"

	apimodel "github.com/authgear/authgear-nft-indexer/pkg/api/model"
	"github.com/authgear/authgear-nft-indexer/pkg/model/database"
	authgearapi "github.com/authgear/authgear-server/pkg/api"
	"github.com/authgear/authgear-server/pkg/api/apierrors"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"github.com/authgear/authgear-server/pkg/util/log"
	authgearweb3 "github.com/authgear/authgear-server/pkg/util/web3"
)

func ConfigureGetCollectionMetadataRoute(route httproute.Route) httproute.Route {
	return route.
		WithMethods("GET").
		WithPathPattern("/metadata")
}

type GetCollectionMetadataHandlerLogger struct{ *log.Logger }

func NewGetCollectionMetadataHandlerLogger(lf *log.Factory) GetCollectionMetadataHandlerLogger {
	return GetCollectionMetadataHandlerLogger{lf.New("api-get-collection-metadata")}
}

type GetCollectionMetadataHandlerMetadataService interface {
	GetContractMetadata(appID string, contracts []authgearweb3.ContractID) ([]database.NFTCollection, error)
}

type GetCollectionMetadataAPIHandler struct {
	JSON            JSONResponseWriter
	Logger          GetCollectionMetadataHandlerLogger
	MetadataService GetCollectionMetadataHandlerMetadataService
}

func (h *GetCollectionMetadataAPIHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	appID := query.Get("app_id")
	if appID == "" {
		h.Logger.Error("missing app id")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("missing app id")})
		return
	}
	urlValues := req.URL.Query()

	contracts := make([]authgearweb3.ContractID, 0)
	for _, url := range urlValues["contract_id"] {
		e, err := authgearweb3.ParseContractID(url)
		if err != nil {
			h.Logger.WithError(err).Error("failed to parse contract ID")
			h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("invalid contract ID")})
			return
		}

		contracts = append(contracts, *e)
	}

	if len(contracts) == 0 {
		h.Logger.Error("invalid contract ID")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("missing contract ID")})
		return
	}

	metadatas, err := h.MetadataService.GetContractMetadata(appID, contracts)
	if err != nil {
		h.Logger.WithError(err).Error("failed to get contract metadata")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewInternalError("failed to get contract metadata")})
		return
	}

	res := make([]apimodel.NFTCollection, 0, len(contracts))
	for _, metadata := range metadatas {
		res = append(res, metadata.ToAPIModel())
	}

	h.JSON.WriteResponse(resp, &authgearapi.Response{
		Result: &apimodel.GetContractMetadataResponse{
			Collections: res,
		},
	})
}
