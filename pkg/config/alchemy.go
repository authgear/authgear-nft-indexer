package config

var _ = Schema.Add("AlchemyConfig", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"blockchain": { "type": "string" },
		"network": { "type": "string" },
		"endpoint": { "type": "string" },
		"api_key": { "type": "string" }
	},
	"required": ["api_key"]
}
`)

type AlchemyConfig struct {
	APIKey     string `json:"api_key"`
	Blockchain string `json:"blockchain"`
	Network    string `json:"network"`
	EndPoint   string `json:"endpoint"`
}
