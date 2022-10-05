package handler

import (
	"fmt"
	"net/http"

	apimodel "github.com/authgear/authgear-nft-indexer/pkg/api/model"
	"github.com/authgear/authgear-nft-indexer/pkg/config"
	ethmodel "github.com/authgear/authgear-nft-indexer/pkg/model/eth"
	"github.com/authgear/authgear-nft-indexer/pkg/query"
	authgearapi "github.com/authgear/authgear-server/pkg/api"
	"github.com/authgear/authgear-server/pkg/api/apierrors"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"github.com/authgear/authgear-server/pkg/util/log"
	authgearweb3 "github.com/authgear/authgear-server/pkg/util/web3"
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

type ListOwnerNFTHandlerNFTCollectionQuery interface {
	QueryAllNFTCollections() ([]ethmodel.NFTCollection, error)
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

	ownerAddress := httproute.GetParam(req, "owner_address")
	if ownerAddress == "" {
		h.Logger.Error("invalid owner address")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("invalid owner address")})
		return
	}

	ownerID, err := authgearweb3.ParseContractID(ownerAddress)
	if err != nil {
		h.Logger.WithError(err).Error("failed to parse owner address")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("invalid owner address")})
		return
	}

	collections, err := h.NFTCollectionQuery.QueryAllNFTCollections()
	if err != nil {
		h.Logger.WithError(err).Error("failed to query nft collections")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewInternalError("failed to query nft collections")})
		return
	}

	contractIDToCollectionMap := make(map[authgearweb3.ContractID]ethmodel.NFTCollection)
	for _, collection := range collections {

		contractID := authgearweb3.ContractID{
			Blockchain:      collection.Blockchain,
			Network:         collection.Network,
			ContractAddress: collection.ContractAddress,
		}

		contractIDToCollectionMap[contractID] = collection
	}

	filteredContracts := make([]authgearweb3.ContractID, 0)
	for _, url := range urlValues["contract_id"] {
		e, err := authgearweb3.ParseContractID(url)
		if err != nil {
			h.Logger.WithError(err).Error("failed to parse contract ID")
			h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("invalid contract ID")})
			return
		}

		// Check whether the contract is being watched or not
		if _, ok := contractIDToCollectionMap[*e]; !ok {
			h.Logger.Error("one or more contract ID(s) are not being watched")
			h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest(fmt.Sprintf("contract %s is not being watched", url))})
			return
		}

		// Filter out contracts that are not in owner's network
		if e.Blockchain == ownerID.Blockchain && e.Network == ownerID.Network {
			filteredContracts = append(filteredContracts, *e)
		}
	}

	// Ensure there are at least one valid contract ID
	if len(filteredContracts) == 0 {
		h.Logger.Error("invalid contract ID")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("missing contract ID")})
		return
	}

	// Start building query
	qb := h.NFTOwnerQuery.NewQueryBuilder()

	qb = qb.WithOwner(ownerID).WithContracts(filteredContracts)

	owners, err := h.NFTOwnerQuery.ExecuteQuery(qb)
	if err != nil {
		h.Logger.WithError(err).Error("failed to list nft owners")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewInternalError("failed to list nft owners")})
		return
	}

	contractIDToTokensMap := make(map[authgearweb3.ContractID][]apimodel.Token, 0)
	for _, ownership := range owners {
		contractID := authgearweb3.ContractID{
			Blockchain:      ownership.Blockchain,
			Network:         ownership.Network,
			ContractAddress: ownership.ContractAddress,
		}

		token := apimodel.Token{
			TokenID: *ownership.TokenID.ToMathBig(),
			TransactionIdentifier: apimodel.TransactionIdentifier{
				Hash: ownership.TransactionHash,
			},
			BlockIdentifier: apimodel.BlockIdentifier{
				Index:     *ownership.BlockNumber.ToMathBig(),
				Timestamp: ownership.BlockTimestamp,
			},
		}

		if _, ok := contractIDToTokensMap[contractID]; !ok {
			contractIDToTokensMap[contractID] = []apimodel.Token{token}
		} else {
			contractIDToTokensMap[contractID] = append(contractIDToTokensMap[contractID], token)
		}
	}

	nfts := make([]apimodel.NFT, 0)
	for contractID, collection := range contractIDToCollectionMap {
		tokens := contractIDToTokensMap[contractID]
		if len(tokens) == 0 {
			continue
		}

		nfts = append(nfts, apimodel.NFT{
			Contract: apimodel.Contract{
				Address: collection.ContractAddress,
				Name:    collection.Name,
			},
			Balance: len(tokens),
			Tokens:  tokens,
		})
	}

	ownership := apimodel.NFTOwnership{
		AccountIdentifier: apimodel.AccountIdentifier{
			Address: ownerID.ContractAddress,
		},
		NetworkIdentifier: apimodel.NetworkIdentifier{
			Blockchain: ownerID.Blockchain,
			Network:    ownerID.Network,
		},
		NFTs: nfts,
	}

	h.JSON.WriteResponse(resp, &authgearapi.Response{
		Result: ownership,
	})
}
