package handler

import (
	"fmt"

	"github.com/authgear/authgear-nft-indexer/pkg/api/model"
	"github.com/authgear/authgear-nft-indexer/pkg/model/eth"
	"github.com/authgear/authgear-server/pkg/api/apierrors"
	"github.com/gin-gonic/gin"
)

type ListCollectionHandlerCollectionsQuery interface {
	QueryNFTCollections() ([]eth.NFTCollection, error)
}
type ListCollectionAPIHandler struct {
	Ctx                *gin.Context
	NFTCollectionQuery ListCollectionHandlerCollectionsQuery
}

func (h *ListCollectionAPIHandler) Handle() {
	collections, err := h.NFTCollectionQuery.QueryNFTCollections()
	if err != nil {
		fmt.Printf("failed to list nft collections: %s", err)
		HandleError(h.Ctx, 500, apierrors.NewInternalError("failed to list nft collections"))
		return
	}

	nftCollections := make([]model.NFTCollection, 0, len(collections))
	for _, collection := range collections {
		nftCollections = append(nftCollections, model.NFTCollection{
			ID:              collection.ID,
			Network:         collection.Network,
			Name:            collection.Name,
			ContractAddress: collection.ContractAddress,
		})
	}

	h.Ctx.JSON(200, &model.CollectionListResponse{
		Items: nftCollections,
	})
}
