package web3

import (
	"net/url"
	"path"
	"strconv"

	"github.com/authgear/authgear-nft-indexer/pkg/config"
)

const (
	EthereumMainnetAlchemyEndpoint = "https://eth-mainnet.alchemyapi.io/v2/"
	EthereumGoerliAlchemyEndpoint  = "https://eth-goerli.g.alchemy.com/v2/"
)

type AlchemyEndpoint struct {
	TransferEndpoint *url.URL
	NFTEndpoint      *url.URL
}

func GetAlchemyEndpoint(alchemyConfig config.AlchemyConfig, blockchain string, network string) (endpoint string, apiKey string) {
	switch blockchain {
	case "ethereum":
		chainID, err := strconv.Atoi(network)
		if err != nil || chainID < 0 {
			panic("Ethereum network should have a positive numeric chain ID")
		}

		switch chainID {
		case 1:
			return EthereumMainnetAlchemyEndpoint, alchemyConfig.GetETHMainnetAPIKey()
		case 5:
			return EthereumGoerliAlchemyEndpoint, alchemyConfig.GetETHGoerliAPIKey()
		default:
			panic("unsupported chain ID")
		}
	}

	panic("unsupported blockchain")
}

func GetRequestEndpoints(config config.AlchemyConfig, blockchain string, network string) (*AlchemyEndpoint, error) {
	endpoint, apiKey := GetAlchemyEndpoint(config, blockchain, network)

	url, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	transferEndpoint := *url
	transferEndpoint.Path = path.Join("v2", apiKey)

	nftEndpoint := *url
	nftEndpoint.Path = path.Join("nft", "v2", apiKey)

	return &AlchemyEndpoint{
		TransferEndpoint: &transferEndpoint,
		NFTEndpoint:      &nftEndpoint,
	}, nil
}
