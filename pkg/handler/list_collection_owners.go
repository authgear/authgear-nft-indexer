package handler

import (
	"fmt"
	"strconv"

	apimodel "github.com/authgear/authgear-nft-indexer/pkg/api/model"
	"github.com/authgear/authgear-nft-indexer/pkg/model"
	"github.com/authgear/authgear-nft-indexer/pkg/query"
	"github.com/authgear/authgear-server/pkg/api/apierrors"
	"github.com/gin-gonic/gin"
)

type ListCollectionOwnersAPIHandler struct {
	Ctx           *gin.Context
	NFTOwnerQuery query.NFTOwnerQuery
}

func (h *ListCollectionOwnersAPIHandler) Handle() {
	blockchain := h.Ctx.Param("blockchain")
	network := h.Ctx.Param("network")

	if blockchain == "" || network == "" {
		fmt.Printf("invalid blockchain or network: %s/%s", blockchain, network)
		HandleError(h.Ctx, 400, apierrors.NewBadRequest("invalid blockchain network"))
		return
	}

	blockchainNetwork := model.BlockchainNetwork{
		Blockchain: blockchain,
		Network:    network,
	}
	contractAddress := h.Ctx.Param("contract_address")

	if contractAddress == "" {
		fmt.Printf("invalid contract address: %s", contractAddress)
		HandleError(h.Ctx, 400, apierrors.NewBadRequest("invalid contract address"))
		return
	}

	limit := h.Ctx.DefaultQuery("limit", "100")
	limitNum, err := strconv.Atoi(limit)
	if err != nil {
		fmt.Printf("invalid limit: %s", limit)
		HandleError(h.Ctx, 400, apierrors.NewBadRequest("invalid limit"))
		return
	}

	offset := h.Ctx.DefaultQuery("offset", "0")
	offsetNum, err := strconv.Atoi(offset)
	if err != nil {
		fmt.Printf("invalid offset: %s", offset)
		HandleError(h.Ctx, 400, apierrors.NewBadRequest("invalid offset"))
		return
	}

	qb := h.NFTOwnerQuery.NewQueryBuilder()

	qb = qb.WithBlockchainNetwork(blockchainNetwork).WithContractAddress(contractAddress)

	owners, err := h.NFTOwnerQuery.ExecuteQuery(qb, limitNum, offsetNum)
	if err != nil {
		fmt.Printf("failed to list nft owners: %s", err)
		HandleError(h.Ctx, 500, apierrors.NewInternalError("failed to list nft owners"))
		return
	}

	nftOwners := make([]apimodel.NFTOwner, 0, len(owners.Items))
	for _, owner := range owners.Items {
		nftOwners = append(nftOwners, apimodel.NFTOwner{
			AccountIdentifier: apimodel.AccountIdentifier{
				Address: owner.OwnerAddress,
			},
			NetworkIdentifier: apimodel.NetworkIdentifier{
				Blockchain: owner.Blockchain,
				Network:    owner.Network,
			},
			Contract: apimodel.Contract{
				Address: owner.ContractAddress,
			},
			TokenID: *owner.TokenID.ToMathBig(),
			TransactionIdentifier: apimodel.TransactionIdentifier{
				Hash: owner.TransactionHash,
			},
			BlockIdentifier: apimodel.BlockIdentifier{
				Index: *owner.BlockNumber.ToMathBig(),
			},
		})
	}

	h.Ctx.JSON(200, &apimodel.CollectionOwnersResponse{
		Items:      nftOwners,
		TotalCount: owners.TotalCount,
	})
}
