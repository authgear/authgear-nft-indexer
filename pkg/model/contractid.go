package model

import (
	"net/url"
	"strconv"

	web3util "github.com/authgear/authgear-server/pkg/util/web3"
)

type ContractID struct {
	Blockchain      string
	Network         string
	ContractAddress string
}

func ParseContractID(contractURL string) (*ContractID, error) {
	curl, err := url.Parse(contractURL)
	if err != nil {
		return nil, err
	}

	protocol := curl.Scheme

	switch protocol {
	case "ethereum":
		eip681, err := web3util.ParseEIP681(contractURL)
		if err != nil {
			return nil, err
		}

		return &ContractID{
			Blockchain:      "ethereum",
			Network:         strconv.Itoa(eip681.ChainID),
			ContractAddress: eip681.ContractAddress,
		}, nil
	default:
		panic("contract_id: unknown protocol")
	}
}
