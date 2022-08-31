package handler

import (
	"net/http"

	apimodel "github.com/authgear/authgear-nft-indexer/pkg/api/model"
	"github.com/authgear/authgear-nft-indexer/pkg/model"
	ethmodel "github.com/authgear/authgear-nft-indexer/pkg/model/eth"
	authgearapi "github.com/authgear/authgear-server/pkg/api"
	"github.com/authgear/authgear-server/pkg/api/apierrors"
	"github.com/authgear/authgear-server/pkg/lib/ratelimit"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"github.com/authgear/authgear-server/pkg/util/log"
)

func ConfigureGetCollectionMetadataRoute(route httproute.Route) httproute.Route {
	return route.
		WithMethods("GET").
		WithPathPattern("/metadata/:contract_id")
}

type GetCollectionMetadataHandlerLogger struct{ *log.Logger }

func NewGetCollectionMetadataHandlerLogger(lf *log.Factory) GetCollectionMetadataHandlerLogger {
	return GetCollectionMetadataHandlerLogger{lf.New("api-get-collection-metadata")}
}

type GetCollectionMetadataHandlerAlchemyAPI interface {
	GetContractMetadata(blockchain string, network string, contractAddress string) (*apimodel.ContractMetadataResponse, error)
}

type GetCollectionMeatadataRateLimiter interface {
	TakeToken(bucket ratelimit.Bucket) error
}

type GetCollectionMetadataAPIHandler struct {
	JSON        JSONResponseWriter
	Logger      GetCollectionMetadataHandlerLogger
	AlchemyAPI  GetCollectionMetadataHandlerAlchemyAPI
	RateLimiter GetCollectionMeatadataRateLimiter
}

func (h *GetCollectionMetadataAPIHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	contractIDStr := httproute.GetParam(req, "contract_id")

	contractID, err := model.ParseContractID(contractIDStr)
	if err != nil {
		h.Logger.WithError(err).Error("failed to parse contract URL")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("invalid contract URL")})
		return
	}

	err = h.RateLimiter.TakeToken(AntiSpamContractMetadataRequestBucket())
	if err != nil {
		h.Logger.WithError(err).Error("unable to take token from rate limiter")
		h.JSON.WriteResponse(resp, &authgearapi.Response{
			Error: apierrors.TooManyRequest.WithReason(string(apierrors.TooManyRequest)).New("rate limited"),
		})
		return
	}

	contractMetadata, err := h.AlchemyAPI.GetContractMetadata(contractID.Blockchain, contractID.Network, contractID.ContractAddress)
	if err != nil {
		h.Logger.WithError(err).Error("failed to get contract metadata")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewInternalError("failed to get contract metadata")})
		return
	}

	tokenType, err := ethmodel.ParseNFTCollectionType(contractMetadata.ContractMetadata.TokenType)
	if err != nil {
		h.Logger.WithError(err).Error("failed to parse token type")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewInternalError("failed to parse token type")})
		return
	}

	h.JSON.WriteResponse(resp, &authgearapi.Response{
		Result: &apimodel.GetContractMetadataResponse{
			Address: contractMetadata.Address,
			ContractMetadata: apimodel.GetContractMetadataContractMetadata{
				Name:        contractMetadata.ContractMetadata.Name,
				Symbol:      contractMetadata.ContractMetadata.Symbol,
				TotalSupply: contractMetadata.ContractMetadata.TotalSupply,
				TokenType:   string(tokenType),
			},
		},
	})
}
