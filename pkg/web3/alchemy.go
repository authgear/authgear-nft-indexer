package web3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	apimodel "github.com/authgear/authgear-nft-indexer/pkg/api/model"
	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/authgear/authgear-server/pkg/util/hexstring"
	authgearweb3 "github.com/authgear/authgear-server/pkg/util/web3"
)

type AlchemyAPI struct {
	Config config.Config
}

type GetAssetTransferParams struct {
	ContractIDs []authgearweb3.ContractID
	FromAddress string
	ToAddress   string
	FromBlock   string
	ToBlock     string
	PageKey     string
	MaxCount    int64
	Order       string
}

func (p *GetAssetTransferParams) ToRequestParams() (*apimodel.AssetTransferRequestParams, error) {
	contractAddresses := make([]string, 0, len(p.ContractIDs))
	for _, contractID := range p.ContractIDs {
		contractAddresses = append(contractAddresses, contractID.ContractAddress)
	}

	maxCountHex, err := hexstring.NewFromInt64(p.MaxCount)
	if err != nil {
		return nil, fmt.Errorf("invalid maxCount: %w", err)
	}

	return &apimodel.AssetTransferRequestParams{
		ContractAddresses: contractAddresses,
		FromBlock:         p.FromBlock,
		ToBlock:           p.ToBlock,
		FromAddress:       p.FromAddress,
		ToAddress:         p.ToAddress,
		PageKey:           p.PageKey,
		Order:             p.Order,
		MaxCount:          maxCountHex.String(),
		Category:          []string{"erc1155", "erc721"},
		ExcludeZeroValue:  true,
		WithMetadata:      true,
	}, nil
}

func (a *AlchemyAPI) GetOwnerNFTs(ownerAddress string, contractIDs []authgearweb3.ContractID, pageKey string) (*apimodel.GetNFTsResponse, error) {
	blockchain := ""
	network := ""
	contractAddresses := make([]string, 0, len(contractIDs))

	for _, contractID := range contractIDs {
		if blockchain == "" && network == "" {
			blockchain = contractID.Blockchain
			network = contractID.Network
		} else if blockchain != contractID.Blockchain || network != contractID.Network {
			return nil, fmt.Errorf("Invalid contract IDs, blockchain networks are not the same")
		}

		contractAddresses = append(contractAddresses, contractID.ContractAddress)
	}

	alchemyEndpoints, err := GetRequestEndpoints(a.Config.Alchemy, blockchain, network)
	if err != nil {
		return nil, err
	}

	requestURL := alchemyEndpoints.NFTEndpoint
	requestURL.Path = path.Join(requestURL.Path, "getNFTs")

	requestQuery := requestURL.Query()
	requestQuery.Set("owner", ownerAddress)
	requestQuery.Set("withMetadata", "true")
	requestQuery["contractAddresses"] = contractAddresses

	if pageKey != "" {
		requestQuery.Set("pageKey", pageKey)
	}

	requestURL.RawQuery = requestQuery.Encode()

	res, err := http.Get(requestURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	var response apimodel.GetNFTsResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &response, nil
}

func (a *AlchemyAPI) GetAssetTransfers(params GetAssetTransferParams) (*apimodel.AssetTransferResponse, error) {
	blockchain := ""
	network := ""
	for _, contractID := range params.ContractIDs {
		if blockchain == "" && network == "" {
			blockchain = contractID.Blockchain
			network = contractID.Network
		} else if blockchain != contractID.Blockchain || network != contractID.Network {
			return nil, fmt.Errorf("Invalid contract IDs, blockchain networks are not the same")
		}
	}

	alchemyEndpoints, err := GetRequestEndpoints(a.Config.Alchemy, blockchain, network)
	if err != nil {
		return nil, err
	}

	requestParams, err := params.ToRequestParams()
	if err != nil {
		return nil, err
	}

	body := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "alchemy_getAssetTransfers",
		"params":  requestParams,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json: %w", err)
	}

	requestURL := alchemyEndpoints.TransferEndpoint

	res, err := http.Post(requestURL.String(), "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	var response apimodel.AssetTransferResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("failed to get NFT transfers: %s", response.Error.Message)
	}

	return &response, nil
}

func (a *AlchemyAPI) GetContractMetadata(contractID authgearweb3.ContractID) (*apimodel.ContractMetadataResponse, error) {
	alchemyEndpoints, err := GetRequestEndpoints(a.Config.Alchemy, contractID.Blockchain, contractID.Network)

	if err != nil {
		return nil, err
	}

	if contractID.ContractAddress == "" {
		return nil, fmt.Errorf("contractAddress is empty")
	}

	requestURL := alchemyEndpoints.TransferEndpoint
	requestURL.Path = path.Join(requestURL.Path, "getContractMetadata")

	requestQuery := requestURL.Query()
	requestQuery.Set("contractAddress", contractID.ContractAddress)

	requestURL.RawQuery = requestQuery.Encode()

	res, err := http.Get(requestURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	var response apimodel.ContractMetadataResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &response, nil
}

func (a *AlchemyAPI) GetOwnersForCollection(contractID authgearweb3.ContractID) (*apimodel.GetOwnersForCollectionResponse, error) {
	alchemyEndpoints, err := GetRequestEndpoints(a.Config.Alchemy, contractID.Blockchain, contractID.Network)
	if err != nil {
		return nil, err
	}

	if contractID.ContractAddress == "" {
		return nil, fmt.Errorf("contractAddress is empty")
	}

	requestURL := alchemyEndpoints.NFTEndpoint
	requestURL.Path = path.Join(requestURL.Path, "getOwnersForCollection")

	requestQuery := requestURL.Query()
	requestQuery.Set("contractAddress", contractID.ContractAddress)

	requestURL.RawQuery = requestQuery.Encode()

	res, err := http.Get(requestURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	var response apimodel.GetOwnersForCollectionResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &response, nil

}
