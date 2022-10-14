package handler

import (
	"encoding/json"
	"net/http"

	apimodel "github.com/authgear/authgear-nft-indexer/pkg/api/model"
	"github.com/authgear/authgear-nft-indexer/pkg/model/database"
	authgearapi "github.com/authgear/authgear-server/pkg/api"
	"github.com/authgear/authgear-server/pkg/api/apierrors"
	"github.com/authgear/authgear-server/pkg/lib/ratelimit"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"github.com/authgear/authgear-server/pkg/util/log"
	authgearweb3 "github.com/authgear/authgear-server/pkg/util/web3"
)

func ConfigureGetCollectionMetadataRoute(route httproute.Route) httproute.Route {
	return route.
		WithMethods("POST").
		WithPathPattern("/metadata")
}

type GetCollectionMetadataHandlerLogger struct{ *log.Logger }

func NewGetCollectionMetadataHandlerLogger(lf *log.Factory) GetCollectionMetadataHandlerLogger {
	return GetCollectionMetadataHandlerLogger{lf.New("api-get-collection-metadata")}
}

type GetCollectionMetadataHandlerRateLimiter interface {
	TakeToken(bucket ratelimit.Bucket) error
}

type GetCollectionMetadataHandlerMetadataService interface {
	GetContractMetadata(contracts []authgearweb3.ContractID) ([]database.NFTCollection, error)
}

type GetCollectionMetadataAPIHandler struct {
	JSON            JSONResponseWriter
	Logger          GetCollectionMetadataHandlerLogger
	MetadataService GetCollectionMetadataHandlerMetadataService
	RateLimiter     GetCollectionMetadataHandlerRateLimiter
}

func (h *GetCollectionMetadataAPIHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	var body apimodel.GetContractMetadataRequestData

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

	err = h.RateLimiter.TakeToken(AntiSpamContractMetadataRequestBucket(body.AppID))
	if err != nil {
		h.Logger.WithError(err).Error("unable to take token from rate limiter")
		h.JSON.WriteResponse(resp, &authgearapi.Response{
			Error: apierrors.TooManyRequest.WithReason(string(apierrors.TooManyRequest)).New("rate limited"),
		})
		return
	}

	contracts := body.ContractIDs
	if len(contracts) == 0 {
		h.Logger.Error("invalid contract ID")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("missing contract ID")})
		return
	}

	metadatas, err := h.MetadataService.GetContractMetadata(contracts)
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
