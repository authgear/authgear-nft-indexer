package config

var _ = Schema.Add("RedisConfig", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"server": { "type": "string" },
		"database": { "type": "string" },
		"pool_size": { "type": "integer" }
	},
	"required": ["server", "database", "pool_size"]
}
`)

type RedisConfig struct {
	Server   string `json:"server"`
	Database string `json:"database"`
	PoolSize int    `json:"pool_size"`
}
