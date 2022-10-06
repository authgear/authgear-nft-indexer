package handler

import (
	"encoding/json"
	"math/big"
	"net/http"

	apimodel "github.com/authgear/authgear-nft-indexer/pkg/api/model"
	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/authgear/authgear-nft-indexer/pkg/model/eth"
	ethmodel "github.com/authgear/authgear-nft-indexer/pkg/model/eth"
	authgearapi "github.com/authgear/authgear-server/pkg/api"
	"github.com/authgear/authgear-server/pkg/api/apierrors"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"github.com/authgear/authgear-server/pkg/util/log"
	authgearweb3 "github.com/authgear/authgear-server/pkg/util/web3"
	"github.com/jrallison/go-workers"
)

func ConfigureWatchCollectionRoute(route httproute.Route) httproute.Route {
	return route.
		WithMethods("POST").
		WithPathPattern("/watch")
}

type WatchCollectionHandlerLogger struct{ *log.Logger }

func NewWatchCollectionHandlerLogger(lf *log.Factory) WatchCollectionHandlerLogger {
	return WatchCollectionHandlerLogger{lf.New("api-watch-collection")}
}

type WatchCollectionHandlerAlchemyAPI interface {
	GetContractMetadata(blockchain string, network string, contractAddress string) (*apimodel.ContractMetadataResponse, error)
}

type WatchCollectionHandlerNFTCollectionMutator interface {
	InsertNFTCollection(blockchain string, network string, name string, contractAddress string, tokenType eth.NFTCollectionType, totalSupply *big.Int) (*ethmodel.NFTCollection, error)
}

type WatchCollectionAPIHandler struct {
	JSON                 JSONResponseWriter
	Logger               WatchCollectionHandlerLogger
	Config               config.Config
	AlchemyAPI           WatchCollectionHandlerAlchemyAPI
	NFTCollectionMutator WatchCollectionHandlerNFTCollectionMutator
}

func (h *WatchCollectionAPIHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	var body apimodel.WatchCollectionRequestData

	defer req.Body.Close()
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		h.Logger.WithError(err).Error("failed to decode request body")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("failed to decode request body")})
		return
	}

	if body.ContractID == "" {
		h.Logger.Error("missing contract_id")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("missing contract_id")})
		return
	}

	contractID, err := authgearweb3.ParseContractID(body.ContractID)
	if err != nil {
		h.Logger.WithError(err).Error("invalid contract_id")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("invalid contract_id")})
		return
	}

	contractMetadata, err := h.AlchemyAPI.GetContractMetadata(contractID.Blockchain, contractID.Network, contractID.ContractAddress)
	if err != nil {
		h.Logger.WithError(err).Error("failed to get contract metadata")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewInternalError("failed to get contract metadata")})
		return
	}

	contractName := body.Name
	if contractName == "" {
		contractName = contractMetadata.ContractMetadata.Name
	}

	tokenType, err := eth.ParseNFTCollectionType(contractMetadata.ContractMetadata.TokenType)
	if err != nil {
		h.Logger.WithError(err).Error("failed to parse token type")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewInternalError("failed to parse token type")})
		return
	}

	totalSupply := new(big.Int)
	totalSupply, ok := totalSupply.SetString(contractMetadata.ContractMetadata.TotalSupply, 10)
	if !ok {
		totalSupply = nil
	}

	collection, err := h.NFTCollectionMutator.InsertNFTCollection(
		contractID.Blockchain,
		contractID.Network,
		contractName,
		contractID.ContractAddress,
		tokenType,
		totalSupply,
	)
	if err != nil {
		h.Logger.WithError(err).Error("failed to insert nft collection")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewInternalError("failed to insert nft collection")})
		return
	}

	_, err = workers.Enqueue(h.Config.Worker.CollectionQueueName, "", nil)
	if err != nil {
		h.Logger.WithError(err).Error("failed to enqueue collection")
	}

	apiCollection := collection.ToAPIModel()

	h.JSON.WriteResponse(resp, &authgearapi.Response{
		Result: &apiCollection,
	})

}
