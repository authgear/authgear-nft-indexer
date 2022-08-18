package handler

import (
	"net/http"

	apimodel "github.com/authgear/authgear-nft-indexer/pkg/api/model"
	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/authgear/authgear-nft-indexer/pkg/query"
	urlutil "github.com/authgear/authgear-nft-indexer/pkg/util/url"
	authgearapi "github.com/authgear/authgear-server/pkg/api"
	"github.com/authgear/authgear-server/pkg/api/apierrors"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"github.com/authgear/authgear-server/pkg/util/log"
)

func ConfigureListOwnerNFTRoute(route httproute.Route) httproute.Route {
	return route.
		WithMethods("GET").
		WithPathPattern("/nfts/:owner_address")
}

type ListOwnerNFTHandlerLogger struct{ *log.Logger }

func NewListOwnerNFTHandlerLogger(lf *log.Factory) ListOwnerNFTHandlerLogger {
	return ListOwnerNFTHandlerLogger{lf.New("api-list-owner-nft")}
}

type ListOwnerNFTAPIHandler struct {
	JSON          JSONResponseWriter
	Logger        ListOwnerNFTHandlerLogger
	Config        config.Config
	NFTOwnerQuery query.NFTOwnerQuery
}

func (h *ListOwnerNFTAPIHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ownerAddress := httproute.GetParam(req, "owner_address")

	if ownerAddress == "" {
		h.Logger.Error("invalid owner address")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("invalid owner address")})
		return
	}

	limit, offset, err := urlutil.ParsePaginationParams(req.URL.Query(), 10, 0)
	if err != nil {
		h.Logger.Error("invalid pagination")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: err})
		return
	}

	qb := h.NFTOwnerQuery.NewQueryBuilder()

	qb = qb.WithOwnerAddress(ownerAddress)

	owners, err := h.NFTOwnerQuery.ExecuteQuery(qb, limit, offset)

	if err != nil {
		h.Logger.Error("failed to list nft owners")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("failed to list nft owners")})
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

	h.JSON.WriteResponse(resp, &authgearapi.Response{
		Result: &apimodel.CollectionOwnersResponse{
			Items:      nftOwners,
			TotalCount: owners.TotalCount,
		},
	})
}
