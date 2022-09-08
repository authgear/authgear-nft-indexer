package config

var _ = Schema.Add("AlchemyConfig", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"blockchain": { "type": "string" },
		"network": { "type": "string" },
		"api_key": { "type": "string" }
	},
	"required": ["blockchain", "network", "api_key"]
}
`)

type AlchemyConfig struct {
	Blockchain string `json:"blockchain"`
	Network    string `json:"network"`
	APIKey     string `json:"api_key"`
}
