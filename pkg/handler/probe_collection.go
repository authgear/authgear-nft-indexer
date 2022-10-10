package handler

import (
	"encoding/json"
	"net/http"

	apimodel "github.com/authgear/authgear-nft-indexer/pkg/api/model"
	"github.com/authgear/authgear-nft-indexer/pkg/model/alchemy"
	authgearapi "github.com/authgear/authgear-server/pkg/api"
	"github.com/authgear/authgear-server/pkg/api/apierrors"
	"github.com/authgear/authgear-server/pkg/lib/ratelimit"
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

type ProbeCollectionHandlerAlchemyAPI interface {
	GetOwnersForCollection(contractID authgearweb3.ContractID) (*alchemy.GetOwnersForCollectionResponse, error)
}

type ProbeCollectionHandlerRateLimiter interface {
	TakeToken(bucket ratelimit.Bucket) error
}

type ProbeCollectionAPIHandler struct {
	JSON        JSONResponseWriter
	Logger      ProbeCollectionHandlerLogger
	AlchemyAPI  ProbeCollectionHandlerAlchemyAPI
	RateLimiter ProbeCollectionHandlerRateLimiter
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

	if body.AppID == "" {
		h.Logger.Error("missing app id")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("missing app id")})
		return
	}

	err = h.RateLimiter.TakeToken(AntiSpamProbeCollectionRequestBucket(body.AppID))
	if err != nil {
		h.Logger.WithError(err).Error("unable to take token from rate limiter")
		h.JSON.WriteResponse(resp, &authgearapi.Response{
			Error: apierrors.TooManyRequest.WithReason(string(apierrors.TooManyRequest)).New("rate limited"),
		})
		return
	}

	if body.ContractID == "" {
		h.Logger.Error("missing contract_id")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("missing contract_id")})
		return
	}

	contractID, err := authgearweb3.ParseContractID(body.ContractID)
	if err != nil {
		h.Logger.WithError(err).Error("failed to parse contract URL")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("invalid contract URL")})
		return
	}

	res, err := h.AlchemyAPI.GetOwnersForCollection(*contractID)
	if err != nil {
		h.Logger.WithError(err).Error("failed to prob collection")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("failed to probe collection")})
		return
	}

	h.JSON.WriteResponse(resp, &authgearapi.Response{
		Result: &apimodel.ProbeCollectionResponse{
			IsLargeCollection: res.PageKey != nil,
		},
	})

}
