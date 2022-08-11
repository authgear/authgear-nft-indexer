package handler

import (
	"fmt"

	apimodel "github.com/authgear/authgear-nft-indexer/pkg/api/model"
	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/authgear/authgear-nft-indexer/pkg/model"
	ethmodel "github.com/authgear/authgear-nft-indexer/pkg/model/eth"
	"github.com/authgear/authgear-server/pkg/api/apierrors"
	"github.com/gin-gonic/gin"
	"github.com/jrallison/go-workers"
)

type RegisterCollectionHandlerAlchemyAPI interface {
	GetContractMetadata(blockchainNetwork model.BlockchainNetwork, contractAddress string) (*apimodel.ContractMetadataResponse, error)
}

type RegisterCollectionHandlerNFTCollectionMutator interface {
	InsertNFTCollection(blockchainNetwork model.BlockchainNetwork, name string, contractAddress string) (*ethmodel.NFTCollection, error)
}

type RegisterCollectionAPIHandler struct {
	Ctx                  *gin.Context
	Config               config.Config
	AlchemyAPI           RegisterCollectionHandlerAlchemyAPI
	NFTCollectionMutator RegisterCollectionHandlerNFTCollectionMutator
}

func (h *RegisterCollectionAPIHandler) Handle() {
	var body apimodel.CollectionRegistrationRequestData
	if err := h.Ctx.BindJSON(&body); err != nil {
		fmt.Printf("invalid request body: %s", err)
		HandleError(h.Ctx, 400, apierrors.NewInternalError("invalid request body"))
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
			fmt.Printf("failed to get contract metadata: %s", err)
			HandleError(h.Ctx, 400, apierrors.NewInternalError("failed to get contract metadata"))
			return
		}

		contractName = contractMetadata.ContractMetadata.Name
	}

	collection, err := h.NFTCollectionMutator.InsertNFTCollection(blockchainNetwork, contractName, body.ContractAddress)
	if err != nil {
		fmt.Printf("Failed to insert nft collection: %s", err)
		HandleError(h.Ctx, 500, apierrors.NewInternalError("failed to insert nft collection"))
		return
	}

	_, err = workers.Enqueue(h.Config.Worker.CollectionQueueName, "", nil)
	if err != nil {
		fmt.Printf("failed to enqueue collection: %s", err)
	}

	h.Ctx.JSON(201, &apimodel.NFTCollection{
		ID:              collection.ID,
		Blockchain:      collection.Blockchain,
		Network:         collection.Network,
		Name:            collection.Name,
		ContractAddress: collection.ContractAddress,
	})
}
