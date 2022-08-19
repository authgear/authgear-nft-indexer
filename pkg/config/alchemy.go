package config

var _ = Schema.Add("AlchemyConfig", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"api_key": { "type": "string" }
	},
	"required": ["api_key"]
}
`)

type AlchemyConfig struct {
	APIKey string `json:"api_key"`
}
