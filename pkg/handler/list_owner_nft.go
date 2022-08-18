package handler

import (
	"net/http"
	"strings"

	apimodel "github.com/authgear/authgear-nft-indexer/pkg/api/model"
	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/authgear/authgear-nft-indexer/pkg/model"
	ethmodel "github.com/authgear/authgear-nft-indexer/pkg/model/eth"
	"github.com/authgear/authgear-nft-indexer/pkg/query"
	urlutil "github.com/authgear/authgear-nft-indexer/pkg/util/url"
	authgearapi "github.com/authgear/authgear-server/pkg/api"
	"github.com/authgear/authgear-server/pkg/api/apierrors"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"github.com/authgear/authgear-server/pkg/util/log"
)

type CollectionIdentifier struct {
	Blockchain      string
	Network         string
	ContractAddress string
}

func ConfigureListOwnerNFTRoute(route httproute.Route) httproute.Route {
	return route.
		WithMethods("GET").
		WithPathPattern("/nfts/:owner_address")
}

type ListOwnerNFTHandlerLogger struct{ *log.Logger }

func NewListOwnerNFTHandlerLogger(lf *log.Factory) ListOwnerNFTHandlerLogger {
	return ListOwnerNFTHandlerLogger{lf.New("api-list-owner-nft")}
}

type ListOwnerNFTHandlerNFTCollectionQuery interface {
	QueryNFTCollections() ([]ethmodel.NFTCollection, error)
}
type ListOwnerNFTAPIHandler struct {
	JSON               JSONResponseWriter
	Logger             ListOwnerNFTHandlerLogger
	Config             config.Config
	NFTOwnerQuery      query.NFTOwnerQuery
	NFTCollectionQuery ListOwnerNFTHandlerNFTCollectionQuery
}

func (h *ListOwnerNFTAPIHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ownerAddress := httproute.GetParam(req, "owner_address")

	if ownerAddress == "" {
		h.Logger.Error("invalid owner address")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("invalid owner address")})
		return
	}

	urlValues := req.URL.Query()

	limit, offset, err := urlutil.ParsePaginationParams(urlValues, 10, 0)
	if err != nil {
		h.Logger.Error("invalid pagination")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: err})
		return
	}

	contractAddresses := make([]string, 0)
	contractAddressesStr := urlValues.Get("contract_addresses")
	for _, contractAddress := range strings.Split(contractAddressesStr, ",") {
		if contractAddress != "" {
			contractAddresses = append(contractAddresses, contractAddress)
		}

	}

	blockchain := urlValues.Get("blockchain")
	network := urlValues.Get("network")

	collections, err := h.NFTCollectionQuery.QueryNFTCollections()
	if err != nil {
		h.Logger.WithError(err).Error("failed to query nft collections")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: err})
		return
	}

	collectionMap := make(map[CollectionIdentifier]ethmodel.NFTCollection)
	for _, collection := range collections {
		collectionMap[CollectionIdentifier{
			Blockchain:      collection.Blockchain,
			Network:         collection.Network,
			ContractAddress: collection.ContractAddress,
		}] = collection
	}

	// Start building query
	qb := h.NFTOwnerQuery.NewQueryBuilder()

	qb = qb.WithOwnerAddress(ownerAddress)

	if blockchain != "" && network != "" {
		qb = qb.WithBlockchainNetwork(model.BlockchainNetwork{
			Blockchain: blockchain,
			Network:    network,
		})
	}

	if len(contractAddresses) > 0 {
		qb = qb.WithContractAddresses(contractAddresses)
	}

	owners, err := h.NFTOwnerQuery.ExecuteQuery(qb, limit, offset)

	if err != nil {
		h.Logger.Error("failed to list nft owners")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("failed to list nft owners")})
		return
	}

	nftOwners := make([]apimodel.NFTOwner, 0, len(owners.Items))
	for _, owner := range owners.Items {
		collection := collectionMap[CollectionIdentifier{
			Blockchain:      owner.Blockchain,
			Network:         owner.Network,
			ContractAddress: owner.ContractAddress,
		}]

		nftOwners = append(nftOwners, apimodel.NFTOwner{
			AccountIdentifier: apimodel.AccountIdentifier{
				Address: owner.OwnerAddress,
			},
			NetworkIdentifier: apimodel.NetworkIdentifier{
				Blockchain: owner.Blockchain,
				Network:    owner.Network,
			},
			Contract: apimodel.Contract{
				Address: collection.ContractAddress,
				Name:    collection.Name,
			},
			TokenID: *owner.TokenID.ToMathBig(),
			TransactionIdentifier: apimodel.TransactionIdentifier{
				Hash: owner.TransactionHash,
			},
			BlockIdentifier: apimodel.BlockIdentifier{
				Index:     *owner.BlockNumber.ToMathBig(),
				Timestamp: owner.BlockTimestamp,
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
