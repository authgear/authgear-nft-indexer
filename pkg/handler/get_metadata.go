package handler

import (
	"math/big"
	"net/http"

	apimodel "github.com/authgear/authgear-nft-indexer/pkg/api/model"
	"github.com/authgear/authgear-nft-indexer/pkg/model"
	"github.com/authgear/authgear-nft-indexer/pkg/model/alchemy"
	dbmodel "github.com/authgear/authgear-nft-indexer/pkg/model/database"
	"github.com/authgear/authgear-nft-indexer/pkg/query"
	authgearapi "github.com/authgear/authgear-server/pkg/api"
	"github.com/authgear/authgear-server/pkg/api/apierrors"
	"github.com/authgear/authgear-server/pkg/lib/ratelimit"
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

type GetCollectionMetadataHandlerAlchemyAPI interface {
	GetContractMetadata(contractID authgearweb3.ContractID) (*alchemy.ContractMetadataResponse, error)
}

type GetCollectionMetadataRateLimiter interface {
	TakeToken(bucket ratelimit.Bucket) error
}

type GetCollectionMetadataNFTCollectionMutator interface {
	InsertNFTCollection(contractID authgearweb3.ContractID, contractName string, tokenType dbmodel.NFTCollectionType, totalSupply *big.Int) (*dbmodel.NFTCollection, error)
}

type GetCollectionMetadataAPIHandler struct {
	JSON                 JSONResponseWriter
	Logger               GetCollectionMetadataHandlerLogger
	AlchemyAPI           GetCollectionMetadataHandlerAlchemyAPI
	NFTCollectionQuery   query.NFTCollectionQuery
	NFTCollectionMutator GetCollectionMetadataNFTCollectionMutator
	RateLimiter          GetCollectionMetadataRateLimiter
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

	// Get existing collections
	qb := h.NFTCollectionQuery.NewQueryBuilder()
	qb = qb.WithContracts(contracts)
	collections, err := h.NFTCollectionQuery.ExecuteQuery(qb)
	if err != nil {
		h.Logger.WithError(err).Error("failed to get collections")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewInternalError("failed to get collections")})
		return
	}

	contractIDToCollectionMap := make(map[string]*dbmodel.NFTCollection)
	for i, collection := range collections {
		contractID, err := collection.ContractID()
		if err != nil {
			h.Logger.WithError(err).Error("failed to parse collection contract ID")
			h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewInternalError("failed to parse collection contract ID")})
			return
		}

		contractURL, err := contractID.URL()
		if err != nil {
			h.Logger.WithError(err).Error("failed to convert collection contract ID to URL")
			h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewInternalError("failed to convert collection contract ID to URL")})
			return
		}

		contractIDToCollectionMap[contractURL.String()] = &collections[i]
	}

	res := make([]apimodel.NFTCollection, 0, len(contracts))
	for _, contract := range contracts {
		contractURL, err := contract.URL()
		if err != nil {
			h.Logger.WithError(err).Error("failed to convert collection contract ID to URL")
			h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewInternalError("failed to convert collection contract ID to URL")})
			return
		}

		// If exists, append to result, otherwise get from alchemy
		collection := contractIDToCollectionMap[contractURL.String()]
		if collection != nil {
			res = append(res, collection.ToAPIModel())
			continue
		}

		err = h.RateLimiter.TakeToken(AntiSpamContractMetadataRequestBucket(appID))
		if err != nil {
			h.Logger.WithError(err).Error("unable to take token from rate limiter")
			h.JSON.WriteResponse(resp, &authgearapi.Response{
				Error: apierrors.TooManyRequest.WithReason(string(apierrors.TooManyRequest)).New("rate limited"),
			})
			return
		}

		contractMetadata, err := h.AlchemyAPI.GetContractMetadata(contract)
		if err != nil {
			h.Logger.WithError(err).Error("failed to get contract metadata")
			h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.BadRequest.WithReason(string(model.BadNFTCollectionError)).New("failed to get contract metadata")})
			return
		}

		tokenType, err := dbmodel.ParseNFTCollectionType(contractMetadata.ContractMetadata.TokenType)
		if err != nil {
			h.Logger.WithError(err).Error("failed to parse token type")
			h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.BadRequest.WithReason(string(model.BadNFTCollectionError)).New("failed to parse token type")})
			return
		}

		if contractMetadata.ContractMetadata.Name == "" {
			h.Logger.WithError(err).Error("missing contract metadata")
			h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.BadRequest.WithReason(string(model.BadNFTCollectionError)).New("missing contract metadata")})
			return
		}

		totalSupply := new(big.Int)
		if contractMetadata.ContractMetadata.TotalSupply != "" {
			if _, ok := totalSupply.SetString(contractMetadata.ContractMetadata.TotalSupply, 10); !ok {
				h.Logger.WithError(err).Error("failed to parse totalSupply")
				h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.BadRequest.WithReason(string(model.BadNFTCollectionError)).New("failed to parse totalSupply")})
				return
			}
		}

		newCollection, err := h.NFTCollectionMutator.InsertNFTCollection(
			contract,
			contractMetadata.ContractMetadata.Name,
			tokenType,
			totalSupply,
		)

		if err != nil {
			h.Logger.WithError(err).Error("failed to insert nft collection")
			h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewInternalError("failed to parse totalSupply")})
			return
		}

		res = append(res, newCollection.ToAPIModel())

	}

	h.JSON.WriteResponse(resp, &authgearapi.Response{
		Result: &apimodel.GetContractMetadataResponse{
			Collections: res,
		},
	})
}
