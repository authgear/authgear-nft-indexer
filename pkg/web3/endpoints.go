package web3

import (
	"net/url"
	"path"
	"strconv"
)

const (
	// Alchemy
	EthereumAlchemyEndpoint = "https://eth-mainnet.alchemyapi.io/v2/"
)

type AlchemyEndpoint struct {
	TransferEndpoint *url.URL
	NFTEndpoint      *url.URL
}

func GetAlchemyEndpoint(blockchain string, network string) string {
	switch blockchain {
	case "ethereum":
		chainID, err := strconv.Atoi(network)
		if err != nil || chainID < 0 {
			panic("Ethereum network should have a positive numeric chain ID")
		}

		switch chainID {
		case 1:
			return EthereumAlchemyEndpoint
		default:
			panic("unsupported chain ID")
		}
	}

	panic("unsupported blockchain")
}

func GetRequestEndpoints(blockchain string, network string) (*AlchemyEndpoint, error) {
	endpoint := GetAlchemyEndpoint(blockchain, network)

	url, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	transferEndpoint := *url
	transferEndpoint.Path = path.Join("v2")

	nftEndpoint := *url
	nftEndpoint.Path = path.Join("nft", "v2")

	return &AlchemyEndpoint{
		TransferEndpoint: &transferEndpoint,
		NFTEndpoint:      &nftEndpoint,
	}, nil
}
