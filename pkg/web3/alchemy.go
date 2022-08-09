package web3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	apimodel "github.com/authgear/authgear-nft-indexer/pkg/api/model"
	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/authgear/authgear-nft-indexer/pkg/model"
	"github.com/authgear/authgear-nft-indexer/pkg/util/hexstring"
)

type AlchemyAPI struct {
	Config config.Config
}

func (a *AlchemyAPI) getNetworkConfig(blockchainNetwork model.BlockchainNetwork) *config.AlchemyConfig {
	alchemyNetworks := a.Config.Alchemy
	for _, alchemyNetwork := range alchemyNetworks {
		if alchemyNetwork.Blockchain == blockchainNetwork.Blockchain && alchemyNetwork.Network == blockchainNetwork.Network {
			return &alchemyNetwork
		}
	}

	return nil
}

func (a *AlchemyAPI) GetNFTTransfers(blockchainNetwork model.BlockchainNetwork, contractAddresses []string, fromBlock string, toBlock string, pageKey string, maxCount int64) (*apimodel.AssetTransferResponse, error) {
	alchemyNetwork := a.getNetworkConfig(blockchainNetwork)
	if alchemyNetwork == nil {
		return nil, fmt.Errorf("network %s %s not found", blockchainNetwork.Blockchain, blockchainNetwork.Network)
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
		return nil, fmt.Errorf("failed to marshal json: %s", err)
	}

	log.Printf("Requesting NFT Transfers for contractAddresses: %s from network %s, fromBlock %s, toBlock %s from endpoint %s", strings.Join(contractAddresses, ", "), alchemyNetwork.Network, fromBlock, toBlock, alchemyNetwork.Endpoint)
	res, err := http.Post(alchemyNetwork.Endpoint+alchemyNetwork.APIKey, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %s", err)
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err)
	}

	var response apimodel.AssetTransferResponse
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %s", err)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("failed to get NFT transfers: %s", response.Error.Message)
	}

	return &response, nil
}
