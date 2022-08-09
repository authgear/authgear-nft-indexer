package config

var _ = Schema.Add("DatabaseConfig", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"url": { "type": "string" },
		"pool_size": { "type": "integer" },
		"verbose": { "type": "boolean" }
	},
	"required": ["url", "pool_size"]
}
`)

type DatabaseConfig struct {
	URL      string `json:"url"`
	PoolSize int    `json:"pool_size"`
	Verbose  bool   `json:"verbose"`
}
