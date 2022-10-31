package web3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/authgear/authgear-nft-indexer/pkg/model/alchemy"
	"github.com/authgear/authgear-server/pkg/util/hexstring"
	authgearweb3 "github.com/authgear/authgear-server/pkg/util/web3"
)

var alchemyClient = &http.Client{
	Timeout: 5 * time.Second,
}

func wrapAlchemyTimeout(err error) error {
	if os.IsTimeout(err) {
		return ErrAlchemyProtocol.Wrap(err, "timeout")
	}

	return err
}

func decodeAlchemyJSON[T any](res *http.Response, tag string, t *T) error {
	var buf bytes.Buffer
	reader := io.TeeReader(res.Body, &buf)

	err := json.NewDecoder(reader).Decode(t)
	if err != nil {
		return ErrAlchemyProtocol.Wrap(err, fmt.Sprintf("%v: %v", tag, buf.String()))
	}

	return nil
}

type AlchemyAPI struct {
	Config config.Config
}

type GetAssetTransferParams struct {
	ContractIDs []authgearweb3.ContractID
	FromAddress authgearweb3.EIP55
	ToAddress   authgearweb3.EIP55
	FromBlock   string
	ToBlock     string
	PageKey     string
	MaxCount    int64
	Order       string
}

func (p *GetAssetTransferParams) ToRequestParams() (*alchemy.AssetTransferRequestParams, error) {
	contractAddresses := make([]authgearweb3.EIP55, 0, len(p.ContractIDs))
	for _, contractID := range p.ContractIDs {
		contractAddresses = append(contractAddresses, contractID.Address)
	}

	maxCountHex, err := hexstring.NewFromInt64(p.MaxCount)
	if err != nil {
		return nil, fmt.Errorf("invalid maxCount: %w", err)
	}

	return &alchemy.AssetTransferRequestParams{
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

func (a *AlchemyAPI) GetOwnerNFTs(ownerAddress string, contractIDs []authgearweb3.ContractID, pageKey string) (*alchemy.GetNFTsResponse, error) {
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

		contractAddresses = append(contractAddresses, contractID.Address.String())
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
	requestQuery[`contractAddresses[]`] = contractAddresses

	if pageKey != "" {
		requestQuery.Set("pageKey", pageKey)
	}

	requestURL.RawQuery = requestQuery.Encode()

	res, err := alchemyClient.Get(requestURL.String())
	if err != nil {
		return nil, wrapAlchemyTimeout(err)
	}
	defer res.Body.Close()

	var response alchemy.GetNFTsResponse
	err = decodeAlchemyJSON(res, "getNFTs", &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (a *AlchemyAPI) GetAssetTransfers(params GetAssetTransferParams) (*alchemy.AssetTransferResult, error) {
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

	res, err := alchemyClient.Post(requestURL.String(), "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, wrapAlchemyTimeout(err)
	}
	defer res.Body.Close()

	var response alchemy.AssetTransferResponse
	err = decodeAlchemyJSON(res, "alchemy_getAssetTransfers", &response)
	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, ErrAlchemyProtocol.New(fmt.Sprintf(
			"alchemy_getAssetTransfers: %v %v",
			response.Error.Code,
			response.Error.Message,
		))
	}

	return &response.Result, nil
}

func (a *AlchemyAPI) GetContractMetadata(contractID authgearweb3.ContractID) (*alchemy.ContractMetadataResponse, error) {
	alchemyEndpoints, err := GetRequestEndpoints(a.Config.Alchemy, contractID.Blockchain, contractID.Network)

	if err != nil {
		return nil, err
	}

	if contractID.Address == "" {
		return nil, fmt.Errorf("contractAddress is empty")
	}

	requestURL := alchemyEndpoints.TransferEndpoint
	requestURL.Path = path.Join(requestURL.Path, "getContractMetadata")

	requestQuery := requestURL.Query()
	requestQuery.Set("contractAddress", contractID.Address.String())

	requestURL.RawQuery = requestQuery.Encode()

	res, err := alchemyClient.Get(requestURL.String())
	if err != nil {
		return nil, wrapAlchemyTimeout(err)
	}
	defer res.Body.Close()

	var response alchemy.ContractMetadataResponse
	err = decodeAlchemyJSON(res, "GetContractMetadata", &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (a *AlchemyAPI) GetOwnersForCollection(contractID authgearweb3.ContractID) (*alchemy.GetOwnersForCollectionResponse, error) {
	alchemyEndpoints, err := GetRequestEndpoints(a.Config.Alchemy, contractID.Blockchain, contractID.Network)
	if err != nil {
		return nil, err
	}

	if contractID.Address == "" {
		return nil, fmt.Errorf("contractAddress is empty")
	}

	requestURL := alchemyEndpoints.NFTEndpoint
	requestURL.Path = path.Join(requestURL.Path, "getOwnersForCollection")

	requestQuery := requestURL.Query()
	requestQuery.Set("contractAddress", contractID.Address.String())

	requestURL.RawQuery = requestQuery.Encode()

	res, err := alchemyClient.Get(requestURL.String())
	if err != nil {
		return nil, wrapAlchemyTimeout(err)
	}
	defer res.Body.Close()

	var response alchemy.GetOwnersForCollectionResponse
	err = decodeAlchemyJSON(res, "getOwnersForCollection", &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
