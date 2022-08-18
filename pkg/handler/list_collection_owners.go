package handler

import (
	"net/http"

	apimodel "github.com/authgear/authgear-nft-indexer/pkg/api/model"
	"github.com/authgear/authgear-nft-indexer/pkg/model"
	"github.com/authgear/authgear-nft-indexer/pkg/query"
	urlutil "github.com/authgear/authgear-nft-indexer/pkg/util/url"
	authgearapi "github.com/authgear/authgear-server/pkg/api"
	"github.com/authgear/authgear-server/pkg/api/apierrors"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"github.com/authgear/authgear-server/pkg/util/log"
)

func ConfigureListCollectionOwnerRoute(route httproute.Route) httproute.Route {
	return route.
		WithMethods("GET").
		WithPathPattern("/collections/:blockchain/:network/owners/:contract_address")
}

type ListCollectionOwnerHandlerLogger struct{ *log.Logger }

func NewListCollectionOwnerHandlerLogger(lf *log.Factory) ListCollectionOwnerHandlerLogger {
	return ListCollectionOwnerHandlerLogger{lf.New("api-list-collection-owner")}
}

type ListCollectionOwnersAPIHandler struct {
	JSON          JSONResponseWriter
	Logger        ListCollectionOwnerHandlerLogger
	NFTOwnerQuery query.NFTOwnerQuery
}

func (h *ListCollectionOwnersAPIHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	blockchain := httproute.GetParam(req, "blockchain")
	network := httproute.GetParam(req, "network")

	if blockchain == "" || network == "" {
		h.Logger.Error("failed to list nft collections")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("invalid blockchain or network")})
		return
	}

	blockchainNetwork := model.BlockchainNetwork{
		Blockchain: blockchain,
		Network:    network,
	}
	contractAddress := httproute.GetParam(req, "contract_address")

	if contractAddress == "" {
		h.Logger.Error("invalid contract address")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("invalid contract address")})
		return
	}

	limit, offset, err := urlutil.ParsePaginationParams(req.URL.Query(), 10, 0)
	if err != nil {
		h.Logger.Error("invalid pagination")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: err})
		return
	}

	qb := h.NFTOwnerQuery.NewQueryBuilder()

	qb = qb.WithBlockchainNetwork(blockchainNetwork).WithContractAddress(contractAddress)

	owners, err := h.NFTOwnerQuery.ExecuteQuery(qb, limit, offset)
	if err != nil {
		h.Logger.Error("failed to list nft owners")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: err})
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
