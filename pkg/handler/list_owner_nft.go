package handler

import (
	"net/http"
	"net/url"
	"time"

	apimodel "github.com/authgear/authgear-nft-indexer/pkg/api/model"
	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/authgear/authgear-nft-indexer/pkg/model/alchemy"
	dbmodel "github.com/authgear/authgear-nft-indexer/pkg/model/database"
	"github.com/authgear/authgear-nft-indexer/pkg/query"
	"github.com/authgear/authgear-nft-indexer/pkg/web3"
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

type ListOwnerNFTHandlerNFTOwnerQuery interface {
	QueryOwner(ownerID authgearweb3.ContractID) (*dbmodel.NFTOwner, error)
}

type ListOwnerNFTHandlerNFTOwnershipMutator interface {
	InsertNFTOwnerships(ownerships []dbmodel.NFTOwnership) error
}

type ListOwnerNFTHandlerAlchemyAPI interface {
	GetOwnerNFTs(ownerAddress string, contractIDs []authgearweb3.ContractID, pageKey string) (*alchemy.GetNFTsResponse, error)
	GetAssetTransfers(params web3.GetAssetTransferParams) (*alchemy.AssetTransferResult, error)
}

func NewListOwnerNFTHandlerLogger(lf *log.Factory) ListOwnerNFTHandlerLogger {
	return ListOwnerNFTHandlerLogger{lf.New("api-list-owner-nft")}
}

type ListOwnerNFTAPIHandler struct {
	JSON                JSONResponseWriter
	Logger              ListOwnerNFTHandlerLogger
	Config              config.Config
	AlchemyAPI          ListOwnerNFTHandlerAlchemyAPI
	NFTOwnerQuery       ListOwnerNFTHandlerNFTOwnerQuery
	NFTCollectionQuery  query.NFTCollectionQuery
	NFTOwnershipQuery   query.NFTOwnershipQuery
	NFTOwnershipMutator ListOwnerNFTHandlerNFTOwnershipMutator
}

func (h *ListOwnerNFTAPIHandler) FetchAndInsertNFTTransfers(ownerID authgearweb3.ContractID, contracts []authgearweb3.ContractID) error {
	pageKey := ""
	nftFetchCount := 0
	ownedNFTs := make([]alchemy.OwnedNFT, 0)
	contractIDsToEnquire := make([]authgearweb3.ContractID, 0)
	// Fetch user nfts until no extra page or has reached the page limit
	for ok := true; ok; ok = pageKey != "" && nftFetchCount <= h.Config.Server.MaxNFTPages {
		nfts, err := h.AlchemyAPI.GetOwnerNFTs(ownerID.Address, contracts, pageKey)
		if err != nil {
			return err
		}

		for _, ownedNFT := range nfts.OwnedNFTs {
			contractID, err := authgearweb3.NewContractID(ownerID.Blockchain, ownerID.Network, ownedNFT.Contract.Address, url.Values{})
			if err != nil {
				return err
			}

			contractIDsToEnquire = append(contractIDsToEnquire, *contractID)
		}

		if nfts.PageKey != nil {
			pageKey = *nfts.PageKey
		}

		ownedNFTs = append(ownedNFTs, nfts.OwnedNFTs...)
		nftFetchCount++
	}

	if len(ownedNFTs) == 0 {
		return nil
	}

	pageKey = ""
	transferFetchCount := 0
	nftTransfers := make([]alchemy.TokenTransfer, 0)
	// Fetch transfers until no extra page or has reached the page limit
	for ok := true; ok; ok = pageKey != "" && transferFetchCount <= 5 {
		transfers, err := h.AlchemyAPI.GetAssetTransfers(web3.GetAssetTransferParams{
			ContractIDs: contractIDsToEnquire,
			ToAddress:   ownerID.Address,
			FromBlock:   "0x0",
			ToBlock:     "latest",
			PageKey:     pageKey,
			MaxCount:    1000,
			Order:       "desc",
		})
		if err != nil {
			return err
		}
		nftTransfers = append(nftTransfers, transfers.Transfers...)
		transferFetchCount++
	}

	ownerships, err := alchemy.MakeNFTOwnerships(ownerID.Blockchain, ownerID.Network, nftTransfers, ownedNFTs)
	if err != nil {
		return err
	}

	// Insert ownerships
	err = h.NFTOwnershipMutator.InsertNFTOwnerships(ownerships)
	if err != nil {
		return err
	}
	return nil
}

func (h *ListOwnerNFTAPIHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	now := time.Now()
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

	contracts := make([]authgearweb3.ContractID, 0)
	for _, url := range urlValues["contract_id"] {
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

	owner, err := h.NFTOwnerQuery.QueryOwner(*ownerID)
	// Check if owner exists and not expired, otherwise fetch from alchemy
	if owner == nil || err != nil || owner.LastSyncedAt.Add(time.Second*time.Duration(h.Config.Server.CacheTTL)).Before(now) {
		err := h.FetchAndInsertNFTTransfers(*ownerID, contracts)
		if err != nil {
			h.Logger.WithError(err).Error("failed to fetch user nfts")
			h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewInternalError("failed to fetch user nfts")})
			return
		}
	}

	// Query ownership from database
	ownershipQb := h.NFTOwnershipQuery.NewQueryBuilder()
	ownershipQb = ownershipQb.WithContracts(contracts).WithOwner(ownerID)
	ownerships, err := h.NFTOwnershipQuery.ExecuteQuery(ownershipQb)
	if err != nil {
		h.Logger.WithError(err).Error("failed to get user nft ownerships")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewInternalError("failed to get user nft ownerships")})
		return
	}

	availableContractIDs := make([]authgearweb3.ContractID, 0)
	for _, ownership := range ownerships {
		contractID, err := ownership.ContractID()
		if err != nil {
			h.Logger.WithError(err).Error("failed to parse ownership contract ID")
			h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewInternalError("failed to parse ownership contract ID")})
			return
		}
		availableContractIDs = append(availableContractIDs, *contractID)
	}

	// Fetch available collections
	collectionQb := h.NFTCollectionQuery.NewQueryBuilder()
	collectionQb = collectionQb.WithContracts(availableContractIDs)
	uniqueCollections, err := h.NFTCollectionQuery.ExecuteQuery(collectionQb)
	if err != nil {
		h.Logger.WithError(err).Error("failed to get nft collections")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewInternalError("failed to get nft collections")})
		return
	}

	// Build response
	nfts := make([]apimodel.NFT, 0)
	for _, collection := range uniqueCollections {
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
