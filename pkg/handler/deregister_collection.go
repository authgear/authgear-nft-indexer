package handler

import (
	"fmt"

	"github.com/authgear/authgear-nft-indexer/pkg/api/model"
	"github.com/authgear/authgear-nft-indexer/pkg/model/eth"
	"github.com/authgear/authgear-server/pkg/api/apierrors"
	"github.com/gin-gonic/gin"
)

type DeregisterCollectionHandlerNFTCollectionMutator interface {
	DeleteNFTCollection(id string) (*eth.NFTCollection, error)
}

type DeregisterCollectionAPIHandler struct {
	Ctx                  *gin.Context
	NFTCollectionMutator DeregisterCollectionHandlerNFTCollectionMutator
}

func (h *DeregisterCollectionAPIHandler) Handle() {
	var body model.CollectionDeregistrationRequestData
	if err := h.Ctx.BindJSON(&body); err != nil {
		fmt.Printf("invalid request body: %s", err)
		HandleError(h.Ctx, 400, apierrors.NewBadRequest("invalid request body"))
		return
	}

	collection, err := h.NFTCollectionMutator.DeleteNFTCollection(body.ID)
	if err != nil {
		fmt.Printf("failed to delete nft collection: %s", err)
		HandleError(h.Ctx, 500, apierrors.NewInternalError("failed to delete nft collection"))
		return
	}

	h.Ctx.JSON(200, &model.NFTCollection{
		ID:              collection.ID,
		Network:         collection.Network,
		Name:            collection.Name,
		ContractAddress: collection.ContractAddress,
	})
}
