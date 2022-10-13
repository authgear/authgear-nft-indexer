package handler

import (
	"encoding/json"
	"net/http"
	"net/url"

	apimodel "github.com/authgear/authgear-nft-indexer/pkg/api/model"
	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/authgear/authgear-nft-indexer/pkg/model/database"
	authgearapi "github.com/authgear/authgear-server/pkg/api"
	"github.com/authgear/authgear-server/pkg/api/apierrors"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"github.com/authgear/authgear-server/pkg/util/log"
	authgearweb3 "github.com/authgear/authgear-server/pkg/util/web3"
)

func ConfigureListOwnerNFTRoute(route httproute.Route) httproute.Route {
	return route.
		WithMethods("POST").
		WithPathPattern("/nfts")
}

type ListOwnerNFTHandlerLogger struct{ *log.Logger }

func NewListOwnerNFTHandlerLogger(lf *log.Factory) ListOwnerNFTHandlerLogger {
	return ListOwnerNFTHandlerLogger{lf.New("api-list-owner-nft")}
}

type ListOwnerNFTHandlerOwnershipService interface {
	GetOwnerships(ownerID authgearweb3.ContractID, contracts []authgearweb3.ContractID) ([]database.NFTOwnership, error)
}

type ListOwnerNFTHandlerMetadataService interface {
	GetContractMetadata(appID string, contracts []authgearweb3.ContractID) ([]database.NFTCollection, error)
}

type ListOwnerNFTAPIHandler struct {
	JSON             JSONResponseWriter
	Logger           ListOwnerNFTHandlerLogger
	Config           config.Config
	OwnershipService ListOwnerNFTHandlerOwnershipService
	MetadataService  ListOwnerNFTHandlerMetadataService
}

func (h *ListOwnerNFTAPIHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	var body apimodel.ListOwnerNFTRequestData
	defer req.Body.Close()
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		h.Logger.WithError(err).Error("failed to decode request body")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("failed to decode request body")})
		return
	}

	if body.AppID == "" {
		h.Logger.Error("missing app id")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("missing app id")})
		return
	}

	if body.OwnerAddress == "" {
		h.Logger.Error("missing owner address")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("missing owner address")})
		return
	}

	if len(body.ContractIDs) == 0 {
		h.Logger.Error("missing contract id")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("missing contract_id")})
		return
	}

	ownerID, err := authgearweb3.ParseContractID(body.OwnerAddress)
	if err != nil {
		h.Logger.WithError(err).Error("failed to parse owner address")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("invalid owner address")})
		return
	}

	contracts := make([]authgearweb3.ContractID, 0)
	for _, url := range body.ContractIDs {
		e, err := authgearweb3.ParseContractID(url)
		if err != nil {
			h.Logger.WithError(err).Error("failed to parse contract ID")
			h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("invalid contract ID")})
			return
		}

		// Filter out contracts that are not in owner's network
		if e.Blockchain == ownerID.Blockchain && e.Network == ownerID.Network {
			contracts = append(contracts, *e)
		}
	}

	// Ensure there are at least one valid contract ID
	if len(contracts) == 0 {
		ownership := apimodel.NewNFTOwnership(*ownerID, []apimodel.NFT{})

		h.JSON.WriteResponse(resp, &authgearapi.Response{
			Result: &ownership,
		})
		return
	}

	collections, err := h.MetadataService.GetContractMetadata(body.AppID, contracts)
	if err != nil {
		h.Logger.WithError(err).Error("failed to get nft collections")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewInternalError("failed to get nft collections")})
		return
	}

	// Check if the input contract IDs have token ids if they are erc1155
	contractIDToCollection := make(map[string]database.NFTCollection)
	for _, collection := range collections {
		contractID, err := collection.ContractID()
		if err != nil {
			h.Logger.WithError(err).Error("failed to contract ID of collection")
			h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewInternalError("failed to contract ID of collection")})
			return
		}

		contractURL, err := contractID.URL()
		if err != nil {
			h.Logger.WithError(err).Error("failed to contract ID of collection")
			h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewInternalError("failed to contract ID of collection")})
			return
		}

		contractIDToCollection[contractURL.String()] = collection
	}

	for _, contract := range contracts {
		tokenIDs := contract.Query["token_ids"]
		contractID, err := authgearweb3.NewContractID(contract.Blockchain, contract.Network, contract.Address, url.Values{})
		if err != nil {
			h.Logger.WithError(err).Error("failed to contract ID of collection")
			h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewInternalError("failed to contract ID of collection")})
			return
		}

		contractURL, err := contractID.URL()
		if err != nil {
			h.Logger.WithError(err).Error("failed to contract ID of collection")
			h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewInternalError("failed to contract ID of collection")})
			return
		}

		collection := contractIDToCollection[contractURL.String()]

		if collection.Type == database.NFTCollectionTypeERC1155 && len(tokenIDs) == 0 {
			h.Logger.Error("erc1155 contract address is specified but token ids are not provided")
			h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewBadRequest("erc1155 contract address is specified but token ids are not provided")})
			return
		}
	}

	ownerships, err := h.OwnershipService.GetOwnerships(*ownerID, contracts)
	if err != nil {
		h.Logger.WithError(err).Error("failed to get nft ownerships")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewInternalError("failed to get nft ownerships")})
		return
	}

	// Build response
	nfts := make([]apimodel.NFT, 0)
	for _, collection := range collections {
		apiNFT := collection.ToAPINFT(ownerships)
		if apiNFT != nil {
			nfts = append(nfts, *apiNFT)
		}

	}

	ownership := apimodel.NewNFTOwnership(*ownerID, nfts)

	h.JSON.WriteResponse(resp, &authgearapi.Response{
		Result: &ownership,
	})
}
