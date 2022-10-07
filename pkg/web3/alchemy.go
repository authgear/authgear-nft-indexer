package web3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"

	apimodel "github.com/authgear/authgear-nft-indexer/pkg/api/model"
	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/authgear/authgear-server/pkg/util/hexstring"
	authgearweb3 "github.com/authgear/authgear-server/pkg/util/web3"
)

type AlchemyAPI struct {
	Config config.Config
}

func (a *AlchemyAPI) GetNFTTransfers(blockchain string, network string, contractAddresses []string, fromBlock string, toBlock string, pageKey string, maxCount int64) (*apimodel.AssetTransferResponse, error) {
	alchemyEndpoints, err := GetRequestEndpoints(a.Config.Alchemy, blockchain, network)
	if err != nil {
		return nil, err
	}

	maxCountHex, err := hexstring.NewFromInt64(maxCount)
	if err != nil {
		return nil, fmt.Errorf("invalid maxCount: %w", err)
	}

	params := &apimodel.AssetTransferRequestParams{
		ContractAddresses: contractAddresses,
		FromBlock:         fromBlock,
		ToBlock:           toBlock,
		PageKey:           pageKey,
		MaxCount:          maxCountHex.String(),
		Category:          []string{"erc721"},
		ExcludeZeroValue:  true,
		WithMetadata:      true,
	}

	body := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "alchemy_getAssetTransfers",
		"params":  params,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json: %w", err)
	}

	requestURL := alchemyEndpoints.TransferEndpoint

	log.Printf("Requesting NFT Transfers for contractAddresses: %s from network %s %s, fromBlock %s, toBlock %s", strings.Join(contractAddresses, ", "), blockchain, network, fromBlock, toBlock)
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

	log.Printf("Requesting contract metadata for contractAddress: %s from network %s %s", contractID.ContractAddress, contractID.Blockchain, contractID.Network)
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

	log.Printf("Requesting owners for contractAddress: %s from network %s %s", contractID.ContractAddress, contractID.Blockchain, contractID.Network)
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
