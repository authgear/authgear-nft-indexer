package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	apimodel "github.com/authgear/authgear-nft-indexer/pkg/api/model"
	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/authgear/authgear-nft-indexer/pkg/model"
	ethmodel "github.com/authgear/authgear-nft-indexer/pkg/model/eth"
	authgearapi "github.com/authgear/authgear-server/pkg/api"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"github.com/authgear/authgear-server/pkg/util/log"
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
	InsertNFTCollection(blockchain string, network string, name string, contractAddress string) (*ethmodel.NFTCollection, error)
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
		h.Logger.WithError(err).Error("Failed to decode request body")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: err})
		return
	}

	if body.ContractID == "" {
		h.Logger.Error("missing contract_id")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: err})
		return
	}

	contractID, err := model.ParseContractID(body.ContractID)
	if err != nil {
		h.Logger.WithError(err).Error("invalid contract_id")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: err})
		return
	}

	contractName := body.Name
	if contractName == "" {
		contractMetadata, err := h.AlchemyAPI.GetContractMetadata(contractID.Blockchain, contractID.Network, contractID.ContractAddress)
		if err != nil {
			h.Logger.WithError(err).Error("Failed to get contract metadata")
			h.JSON.WriteResponse(resp, &authgearapi.Response{Error: err})
			return
		}

		contractName = contractMetadata.ContractMetadata.Name
	}

	collection, err := h.NFTCollectionMutator.InsertNFTCollection(contractID.Blockchain, contractID.Network, contractName, contractID.ContractAddress)
	if err != nil {
		h.Logger.WithError(err).Error("Failed to insert nft collection")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: err})
		return
	}

	_, err = workers.Enqueue(h.Config.Worker.CollectionQueueName, "", nil)
	if err != nil {
		fmt.Printf("failed to enqueue collection: %+v", err)
	}

	h.JSON.WriteResponse(resp, &authgearapi.Response{
		Result: &apimodel.NFTCollection{
			ID:              collection.ID,
			Blockchain:      collection.Blockchain,
			Network:         collection.Network,
			Name:            collection.Name,
			ContractAddress: collection.ContractAddress,
		},
	})

}
