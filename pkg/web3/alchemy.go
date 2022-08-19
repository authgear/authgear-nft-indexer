package web3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"

	apimodel "github.com/authgear/authgear-nft-indexer/pkg/api/model"
	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/authgear/authgear-nft-indexer/pkg/model"
	"github.com/authgear/authgear-server/pkg/util/hexstring"
)

type AlchemyAPI struct {
	Config config.Config
}

type AlchemyEndpoint struct {
	TransferEndpoint string
	NFTEndpoint      string
}

func (a *AlchemyAPI) getRequestEndpoints(blockchainNetwork model.BlockchainNetwork) (*AlchemyEndpoint, error) {
	alchemyNetworks := a.Config.Alchemy
	for _, alchemyNetwork := range alchemyNetworks {
		if alchemyNetwork.Blockchain == blockchainNetwork.Blockchain && alchemyNetwork.Network == blockchainNetwork.Network {

			endpoint, err := url.Parse(alchemyNetwork.Endpoint)
			if err != nil {
				return nil, err
			}
			transferEndpoint := *endpoint
			transferEndpoint.Path = path.Join("v2", alchemyNetwork.APIKey)

			nftEndpoint := *endpoint
			nftEndpoint.Path = path.Join("nft", "v2", alchemyNetwork.APIKey)

			return &AlchemyEndpoint{
				TransferEndpoint: transferEndpoint.String(),
				NFTEndpoint:      nftEndpoint.String(),
			}, nil
		}
	}

	return nil, fmt.Errorf("request endpoint for %s %s not found", blockchainNetwork.Blockchain, blockchainNetwork.Network)
}

func (a *AlchemyAPI) GetNFTTransfers(blockchainNetwork model.BlockchainNetwork, contractAddresses []string, fromBlock string, toBlock string, pageKey string, maxCount int64) (*apimodel.AssetTransferResponse, error) {
	alchemyEndpoints, err := a.getRequestEndpoints(blockchainNetwork)
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

	log.Printf("Requesting NFT Transfers for contractAddresses: %s from network %s %s, fromBlock %s, toBlock %s", strings.Join(contractAddresses, ", "), blockchainNetwork.Blockchain, blockchainNetwork.Network, fromBlock, toBlock)
	res, err := http.Post(alchemyEndpoints.TransferEndpoint, "application/json", bytes.NewBuffer(jsonBody))
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

func (a *AlchemyAPI) GetContractMetadata(blockchainNetwork model.BlockchainNetwork, contractAddress string) (*apimodel.ContractMetadataResponse, error) {
	alchemyEndpoints, err := a.getRequestEndpoints(blockchainNetwork)
	if err != nil {
		return nil, err
	}

	if contractAddress == "" {
		return nil, fmt.Errorf("contractAddress is empty")
	}

	log.Printf("Requesting contract metadata for contractAddress: %s from network %s %s", contractAddress, blockchainNetwork.Blockchain, blockchainNetwork.Network)
	res, err := http.Get(fmt.Sprintf("%s/getContractMetadata?contractAddress=%s", alchemyEndpoints.NFTEndpoint, contractAddress))
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
