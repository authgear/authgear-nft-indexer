package config

var _ = Schema.Add("NetworkName", `
{
	"type": "string",
	"enum": ["eth_mainnet", "eth_goerli"]
}
`)

type NetworkName string

const (
	NetworkNameEthereumMainnet NetworkName = "eth_mainnet"
	NetworkNameEthereumGoerli  NetworkName = "eth_goerli"
)

var _ = Schema.Add("APIKeyConfig", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"name": { "$ref": "#/$defs/NetworkName" },
		"value": { "type": "string" }
	},
	"required": ["name", "value"]
}
`)

type APIKeyConfig struct {
	NetworkName NetworkName `json:"name"`
	APIKey      string      `json:"value"`
}

var _ = Schema.Add("AlchemyConfig", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"api_key": { "type": "array", "items": { "$ref":"#/$defs/APIKeyConfig" } }
	},
	"required": ["api_key"]
}
`)

type AlchemyConfig struct {
	APIKeys []APIKeyConfig `json:"api_key"`
}

func (a *AlchemyConfig) getAPIKeyByNetworkName(name NetworkName) string {
	for _, keyConfig := range a.APIKeys {
		if keyConfig.NetworkName == name {
			return keyConfig.APIKey
		}
	}

	return ""
}

func (a *AlchemyConfig) GetETHMainnetAPIKey() string {
	return a.getAPIKeyByNetworkName(NetworkNameEthereumMainnet)
}

func (a *AlchemyConfig) GetETHGoerliAPIKey() string {
	return a.getAPIKeyByNetworkName(NetworkNameEthereumGoerli)
}
