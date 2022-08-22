package handler

import (
	"net/http"

	apimodel "github.com/authgear/authgear-nft-indexer/pkg/api/model"
	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/authgear/authgear-nft-indexer/pkg/model"
	ethmodel "github.com/authgear/authgear-nft-indexer/pkg/model/eth"
	"github.com/authgear/authgear-nft-indexer/pkg/query"
	authgearapi "github.com/authgear/authgear-server/pkg/api"
	"github.com/authgear/authgear-server/pkg/api/apierrors"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"github.com/authgear/authgear-server/pkg/util/log"
)

func ConfigureListOwnerNFTRoute(route httproute.Route) httproute.Route {
	return route.
		WithMethods("GET").
		WithPathPattern("/nfts")
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
	urlValues := req.URL.Query()

	contracts := make([]model.ContractID, 0)
	for _, url := range urlValues["contract_id"] {
		e, err := model.ParseContractID(url)
		if err != nil {
			h.Logger.WithError(err).Error("failed to parse contract URL")
			h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("invalid contract URL")})
			return
		}

		contracts = append(contracts, *e)
	}

	ownerAddresses := urlValues["owner_address"]

	collections, err := h.NFTCollectionQuery.QueryNFTCollections()
	if err != nil {
		h.Logger.WithError(err).Error("failed to query nft collections")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewInternalError("failed to query nft collections")})
		return
	}

	collectionMap := make(map[model.ContractID]ethmodel.NFTCollection)
	for _, collection := range collections {
		collectionMap[model.ContractID{
			Blockchain:      collection.Blockchain,
			Network:         collection.Network,
			ContractAddress: collection.ContractAddress,
		}] = collection
	}

	// Start building query
	qb := h.NFTOwnerQuery.NewQueryBuilder()

	if len(contracts) > 0 {
		qb = qb.WithContracts(contracts)
	}

	if len(ownerAddresses) > 0 {
		qb = qb.WithOwnerAddresses(ownerAddresses)
	}

	owners, err := h.NFTOwnerQuery.ExecuteQuery(qb)

	if err != nil {
		h.Logger.WithError(err).Error("failed to list nft owners")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewInternalError("failed to list nft owners")})
		return
	}

	nftOwners := make([]apimodel.NFTOwner, 0, len(owners.Items))
	for _, owner := range owners.Items {
		collection := collectionMap[model.ContractID{
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
