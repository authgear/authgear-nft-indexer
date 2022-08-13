package handler

import (
	"fmt"
	"strconv"

	apimodel "github.com/authgear/authgear-nft-indexer/pkg/api/model"
	"github.com/authgear/authgear-nft-indexer/pkg/query"
	"github.com/authgear/authgear-server/pkg/api/apierrors"
	"github.com/gin-gonic/gin"
)

type ListOwnerNFTAPIHandler struct {
	Ctx           *gin.Context
	NFTOwnerQuery query.NFTOwnerQuery
}

func (h *ListOwnerNFTAPIHandler) Handle() {
	ownerAddress := h.Ctx.Param("owner_address")

	if ownerAddress == "" {
		fmt.Printf("invalid owner address: %s", ownerAddress)
		HandleError(h.Ctx, 400, apierrors.NewBadRequest("invalid owner address"))
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

	qb = qb.WithOwnerAddress(ownerAddress)

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
