package handler

import (
	"net/http"
	"time"

	apimodel "github.com/authgear/authgear-nft-indexer/pkg/api/model"
	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/authgear/authgear-nft-indexer/pkg/model/database"
	dbmodel "github.com/authgear/authgear-nft-indexer/pkg/model/database"
	"github.com/authgear/authgear-nft-indexer/pkg/query"
	"github.com/authgear/authgear-nft-indexer/pkg/web3"
	authgearapi "github.com/authgear/authgear-server/pkg/api"
	"github.com/authgear/authgear-server/pkg/api/apierrors"
	"github.com/authgear/authgear-server/pkg/util/hexstring"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"github.com/authgear/authgear-server/pkg/util/log"
	authgearweb3 "github.com/authgear/authgear-server/pkg/util/web3"
	"github.com/uptrace/bun/extra/bunbig"
)

type ContractTokenID struct {
	authgearweb3.ContractID
	TokenID string
}

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
	GetOwnerNFTs(ownerAddress string, contractIDs []authgearweb3.ContractID, pageKey string) (*apimodel.GetNFTsResponse, error)
	GetAssetTransfers(params web3.GetAssetTransferParams) (*apimodel.AssetTransferResponse, error)
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
	ownerships := make([]database.NFTOwnership, 0)
	pageKey := ""
	nftFetchCount := 0
	contractIDsToEnquire := make([]authgearweb3.ContractID, 0)
	contractTokenIDToBalance := make(map[ContractTokenID]string)
	// Fetch user nfts until no extra page or has reached the page limit
	for ok := true; ok; ok = pageKey != "" && nftFetchCount <= h.Config.Server.MaxNFTPages {
		nfts, err := h.AlchemyAPI.GetOwnerNFTs(ownerID.ContractAddress, contracts, pageKey)
		if err != nil {
			return err
		}

		for _, ownedNFT := range nfts.OwnedNFTs {

			contractID := authgearweb3.ContractID{
				Blockchain:      ownerID.Blockchain,
				Network:         ownerID.Network,
				ContractAddress: ownedNFT.Contract.Address,
			}
			contractIDsToEnquire = append(contractIDsToEnquire, contractID)

			contractTokenIDToBalance[ContractTokenID{
				ContractID: contractID,
				TokenID:    ownedNFT.ID.TokenID,
			}] = ownedNFT.Balance
		}

		if nfts.PageKey != nil {
			pageKey = *nfts.PageKey
		}
		nftFetchCount++
	}

	pageKey = ""
	transferFetchCount := 0
	// Fetch transfers until no extra page or has reached the page limit
	for ok := true; ok; ok = pageKey != "" && transferFetchCount <= 5 {
		transfers, err := h.AlchemyAPI.GetAssetTransfers(web3.GetAssetTransferParams{
			ContractIDs: contractIDsToEnquire,
			ToAddress:   ownerID.ContractAddress,
			FromBlock:   "0x0",
			ToBlock:     "latest",
			PageKey:     pageKey,
			MaxCount:    1000,
			Order:       "desc",
		})
		if err != nil {
			return err
		}

		for _, transfer := range transfers.Result.Transfers {

			blockNum := hexstring.MustParse(transfer.BlockNum)

			blockTime, err := time.Parse(time.RFC3339, transfer.Metadata.BlockTimestamp)
			if err != nil {
				return err
			}

			contractID := authgearweb3.ContractID{
				Blockchain:      ownerID.Blockchain,
				Network:         ownerID.Network,
				ContractAddress: transfer.RawContract.Address,
			}

			// Handle ERC-1155
			if transfer.ERC1155Metadata != nil {
				for _, erc1155 := range *transfer.ERC1155Metadata {

					balance := contractTokenIDToBalance[ContractTokenID{
						ContractID: contractID,
						TokenID:    erc1155.TokenID,
					}]

					ownerships = append(ownerships, dbmodel.NFTOwnership{
						Blockchain:      contractID.Blockchain,
						Network:         contractID.Network,
						ContractAddress: contractID.ContractAddress,
						TokenID:         erc1155.TokenID,
						Balance:         balance,
						BlockNumber:     bunbig.FromMathBig(blockNum.ToBigInt()),
						OwnerAddress:    ownerID.ContractAddress,
						TransactionHash: transfer.Hash,
						BlockTimestamp:  blockTime,
					})
				}
				continue
			}

			// Handle ERC-721
			ownerships = append(ownerships, dbmodel.NFTOwnership{
				Blockchain:      contractID.Blockchain,
				Network:         contractID.Network,
				ContractAddress: transfer.RawContract.Address,
				TokenID:         transfer.TokenID,
				Balance:         "1",
				BlockNumber:     bunbig.FromMathBig(blockNum.ToBigInt()),
				OwnerAddress:    ownerID.ContractAddress,
				TransactionHash: transfer.Hash,
				BlockTimestamp:  blockTime,
			})

		}
		transferFetchCount++
	}

	// Insert ownerships
	err := h.NFTOwnershipMutator.InsertNFTOwnerships(ownerships)
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

	tokenIDs := urlValues["token_id"]

	// Ensure there are at least one valid contract ID
	if len(contracts) == 0 {
		ownership := apimodel.NFTOwnership{
			AccountIdentifier: apimodel.AccountIdentifier{
				Address: ownerID.ContractAddress,
			},
			NetworkIdentifier: apimodel.NetworkIdentifier{
				Blockchain: ownerID.Blockchain,
				Network:    ownerID.Network,
			},
			NFTs: []apimodel.NFT{},
		}

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
	if len(tokenIDs) != 0 {
		ownershipQb = ownershipQb.WithTokenIDs(tokenIDs)
	}
	ownerships, err := h.NFTOwnershipQuery.ExecuteQuery(ownershipQb)
	if err != nil {
		h.Logger.WithError(err).Error("failed to get user nft ownerships")
		h.JSON.WriteResponse(resp, &authgearapi.Response{Error: apierrors.NewInternalError("failed to get user nft ownerships")})
		return
	}

	availableContractIDs := make([]authgearweb3.ContractID, 0)
	for _, ownership := range ownerships {
		availableContractIDs = append(availableContractIDs, ownership.ContractID())
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
	contractIDToCollectionMap := make(map[authgearweb3.ContractID]dbmodel.NFTCollection)
	for _, collection := range uniqueCollections {
		contractID := collection.ContractID()

		contractIDToCollectionMap[contractID] = collection
	}

	contractIDToTokensMap := make(map[authgearweb3.ContractID][]apimodel.Token, 0)
	for _, ownership := range ownerships {
		contractID := ownership.ContractID()

		token := apimodel.Token{
			TokenID: ownership.TokenID,
			TransactionIdentifier: apimodel.TransactionIdentifier{
				Hash: ownership.TransactionHash,
			},
			BlockIdentifier: apimodel.BlockIdentifier{
				Index:     *ownership.BlockNumber.ToMathBig(),
				Timestamp: ownership.BlockTimestamp,
			},
			Balance: ownership.Balance,
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
				Type:    string(collection.Type),
			},
			Tokens: tokens,
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
		Result: &ownership,
	})
}
