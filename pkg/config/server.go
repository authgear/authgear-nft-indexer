package config

var _ = Schema.Add("ServerConfig", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"listen_addr": { "type": "string" },
		"cache_ttl": { "type": "integer" },
		"max_nft_pages": { "type": "integer" }
	},
	"required": ["listen_addr", "cache_ttl", "max_nft_pages"]
}
`)

type ServerConfig struct {
	ListenAddr  string `json:"listen_addr"`
	CacheTTL    int    `json:"cache_ttl"`
	MaxNFTPages int    `json:"max_nft_pages"`
}
