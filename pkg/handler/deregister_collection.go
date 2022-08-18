package handler

import (
	"encoding/json"
	"net/http"

	"github.com/authgear/authgear-nft-indexer/pkg/api/model"
	apimodel "github.com/authgear/authgear-nft-indexer/pkg/api/model"
	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/authgear/authgear-nft-indexer/pkg/model/eth"
	authgearapi "github.com/authgear/authgear-server/pkg/api"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"github.com/authgear/authgear-server/pkg/util/log"
)

func ConfigureDeregisterCollectionRoute(route httproute.Route) httproute.Route {
	return route.
		WithMethods("POST").
		WithPathPattern("/deregister")
}

type DeregisterCollectionHandlerLogger struct{ *log.Logger }

func NewDeregisterCollectionHandlerLogger(lf *log.Factory) DeregisterCollectionHandlerLogger {
	return DeregisterCollectionHandlerLogger{lf.New("api-deregister-collection")}
}

type DeregisterCollectionHandlerNFTCollectionMutator interface {
	DeleteNFTCollection(id string) (*eth.NFTCollection, error)
}

type DeregisterCollectionAPIHandler struct {
	JSON                 JSONResponseWriter
	Logger               RegisterCollectionHandlerLogger
	Config               config.Config
	NFTCollectionMutator DeregisterCollectionHandlerNFTCollectionMutator
}

func (h *DeregisterCollectionAPIHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	var body model.CollectionDeregistrationRequestData

	defer req.Body.Close()
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		h.Logger.WithError(err).Error("Failed to decode request body")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: err})
		return
	}

	collection, err := h.NFTCollectionMutator.DeleteNFTCollection(body.ID)
	if err != nil {
		h.Logger.WithError(err).Error("failed to delete nft collection")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: err})
		return
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
