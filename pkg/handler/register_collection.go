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

func ConfigureRegisterCollectionRoute(route httproute.Route) httproute.Route {
	return route.
		WithMethods("POST").
		WithPathPattern("/register")
}

type RegisterCollectionHandlerLogger struct{ *log.Logger }

func NewRegisterCollectionHandlerLogger(lf *log.Factory) RegisterCollectionHandlerLogger {
	return RegisterCollectionHandlerLogger{lf.New("api-register-collection")}
}

type RegisterCollectionHandlerAlchemyAPI interface {
	GetContractMetadata(blockchainNetwork model.BlockchainNetwork, contractAddress string) (*apimodel.ContractMetadataResponse, error)
}

type RegisterCollectionHandlerNFTCollectionMutator interface {
	InsertNFTCollection(blockchainNetwork model.BlockchainNetwork, name string, contractAddress string) (*ethmodel.NFTCollection, error)
}

type RegisterCollectionAPIHandler struct {
	JSON                 JSONResponseWriter
	Logger               RegisterCollectionHandlerLogger
	Config               config.Config
	AlchemyAPI           RegisterCollectionHandlerAlchemyAPI
	NFTCollectionMutator RegisterCollectionHandlerNFTCollectionMutator
}

func (h *RegisterCollectionAPIHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	var body apimodel.CollectionRegistrationRequestData

	defer req.Body.Close()
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		h.Logger.WithError(err).Error("Failed to decode request body")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: err})
		return
	}

	blockchainNetwork := model.BlockchainNetwork{
		Blockchain: body.Blockchain,
		Network:    body.Network,
	}

	contractName := body.Name
	if contractName == "" {
		contractMetadata, err := h.AlchemyAPI.GetContractMetadata(blockchainNetwork, body.ContractAddress)
		if err != nil {
			h.Logger.WithError(err).Error("Failed to get contract metadata")
			h.JSON.WriteResponse(resp, &authgearapi.Response{Error: err})
			return
		}

		contractName = contractMetadata.ContractMetadata.Name
	}

	collection, err := h.NFTCollectionMutator.InsertNFTCollection(blockchainNetwork, contractName, body.ContractAddress)
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
