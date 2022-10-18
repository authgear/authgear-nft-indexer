package handler

import (
	"encoding/json"
	"net/http"

	apimodel "github.com/authgear/authgear-nft-indexer/pkg/api/model"
	authgearapi "github.com/authgear/authgear-server/pkg/api"
	"github.com/authgear/authgear-server/pkg/api/apierrors"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"github.com/authgear/authgear-server/pkg/util/log"
	authgearweb3 "github.com/authgear/authgear-server/pkg/util/web3"
)

func ConfigureProbeCollectionRoute(route httproute.Route) httproute.Route {
	return route.
		WithMethods("POST").
		WithPathPattern("/probe")
}

type ProbeCollectionHandlerLogger struct{ *log.Logger }

func NewProbeCollectionHandlerLogger(lf *log.Factory) ProbeCollectionHandlerLogger {
	return ProbeCollectionHandlerLogger{lf.New("api-probe-collection")}
}

type ProbeCollectionHandlerProbeService interface {
	ProbeCollection(contractID authgearweb3.ContractID) (bool, error)
}
type ProbeCollectionAPIHandler struct {
	JSON         JSONResponseWriter
	Logger       ProbeCollectionHandlerLogger
	ProbeService ProbeCollectionHandlerProbeService
}

func (h *ProbeCollectionAPIHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	var body apimodel.ProbeCollectionRequestData

	defer req.Body.Close()
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		h.Logger.WithError(err).Error("failed to decode request body")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("failed to decode request body")})
		return
	}

	contractID := body.ContractID
	probe, err := h.ProbeService.ProbeCollection(contractID)
	if err != nil {
		h.Logger.WithError(err).Error("failed to probe nft collection")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: err})
		return
	}

	h.JSON.WriteResponse(resp, &authgearapi.Response{
		Result: &apimodel.ProbeCollectionResponse{
			IsLargeCollection: probe,
		},
	})

}
