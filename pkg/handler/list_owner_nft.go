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

	ownerAddress := httproute.GetParam(req, "owner_address")

	if ownerAddress == "" {
		h.Logger.Error("invalid owner address")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("invalid owner address")})
		return
	}

	collections, err := h.NFTCollectionQuery.QueryAllNFTCollections()
	if err != nil {
		h.Logger.WithError(err).Error("failed to query nft collections")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewInternalError("failed to query nft collections")})
		return
	}

	networkToContractIDsMap := make(map[model.BlockchainNetwork][]model.ContractID)
	contractIDToCollectionMap := make(map[model.ContractID]ethmodel.NFTCollection)
	for _, collection := range collections {

		contractID := model.ContractID{
			Blockchain:      collection.Blockchain,
			Network:         collection.Network,
			ContractAddress: collection.ContractAddress,
		}

		contractIDToCollectionMap[contractID] = collection

		blockchainNetwork := model.BlockchainNetwork{
			Blockchain: collection.Blockchain,
			Network:    collection.Network,
		}

		if _, ok := networkToContractIDsMap[blockchainNetwork]; !ok {
			networkToContractIDsMap[blockchainNetwork] = []model.ContractID{contractID}
		} else {
			networkToContractIDsMap[blockchainNetwork] = append(networkToContractIDsMap[blockchainNetwork], contractID)
		}
	}

	// Start building query
	qb := h.NFTOwnerQuery.NewQueryBuilder()

	qb = qb.WithOwnerAddress(ownerAddress)

	if len(contracts) > 0 {
		qb = qb.WithContracts(contracts)
	}

	owners, err := h.NFTOwnerQuery.ExecuteQuery(qb)
	if err != nil {
		h.Logger.WithError(err).Error("failed to list nft owners")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewInternalError("failed to list nft owners")})
		return
	}

	contractIDToTokensMap := make(map[model.ContractID][]apimodel.Token, 0)
	for _, ownership := range owners {
		contractID := model.ContractID{
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
	for contractID, tokens := range contractIDToTokensMap {
		collection := contractIDToCollectionMap[contractID]
		nfts = append(nfts, apimodel.NFT{
			Contract: apimodel.Contract{
				Address: collection.ContractAddress,
				Name:    collection.Name,
			},
			Balance: len(tokens),
			Tokens:  tokens,
		})
	}

	ownerships := make([]apimodel.NFTOwnership, 0)
	for network, contractIDs := range networkToContractIDsMap {
		nfts := make([]apimodel.NFT, 0)
		for _, contractID := range contractIDs {
			collection := contractIDToCollectionMap[contractID]
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

		ownerships = append(ownerships, apimodel.NFTOwnership{
			AccountIdentifier: apimodel.AccountIdentifier{
				Address: ownerAddress,
			},
			NetworkIdentifier: apimodel.NetworkIdentifier{
				Blockchain: network.Blockchain,
				Network:    network.Network,
			},
			NFTs: nfts,
		})
	}

	h.JSON.WriteResponse(resp, &authgearapi.Response{
		Result: ownerships,
	})
}
