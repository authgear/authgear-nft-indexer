package handler

import (
	"net/http"

	"github.com/authgear/authgear-nft-indexer/pkg/api/model"
	"github.com/authgear/authgear-nft-indexer/pkg/model/eth"
	authgearapi "github.com/authgear/authgear-server/pkg/api"
	"github.com/authgear/authgear-server/pkg/api/apierrors"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"github.com/authgear/authgear-server/pkg/util/log"
)

func ConfigureListCollectionRoute(route httproute.Route) httproute.Route {
	return route.
		WithMethods("GET").
		WithPathPattern("/collections")
}

type ListCollectionHandlerLogger struct{ *log.Logger }

func NewListCollectionHandlerLogger(lf *log.Factory) ListCollectionHandlerLogger {
	return ListCollectionHandlerLogger{lf.New("api-list-collection")}
}

type ListCollectionHandlerCollectionsQuery interface {
	QueryNFTCollections() ([]eth.NFTCollection, error)
}
type ListCollectionAPIHandler struct {
	JSON               JSONResponseWriter
	Logger             ListCollectionHandlerLogger
	NFTCollectionQuery ListCollectionHandlerCollectionsQuery
}

func (h *ListCollectionAPIHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	collections, err := h.NFTCollectionQuery.QueryNFTCollections()
	if err != nil {
		h.Logger.WithError(err).Error("failed to list nft collections")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewInternalError("failed to list nft collections")})
		return
	}

	nftCollections := make([]model.NFTCollection, 0, len(collections))
	for _, collection := range collections {
		nftCollections = append(nftCollections, model.NFTCollection{
			ID:              collection.ID,
			Blockchain:      collection.Blockchain,
			Network:         collection.Network,
			Name:            collection.Name,
			ContractAddress: collection.ContractAddress,
			TotalSupply:     *collection.TotalSupply.ToMathBig(),
			Type:            string(collection.Type),
		})
	}

	h.JSON.WriteResponse(resp, &authgearapi.Response{
		Result: &model.CollectionListResponse{
			Items: nftCollections,
		},
	})
}
