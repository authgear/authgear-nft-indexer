package config

var _ = Schema.Add("ServerConfig", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"listen_addr": { "type": "string" },
		"collection_cache_ttl": { "type": "integer" },
		"ownership_cache_ttl": { "type": "integer" },
		"max_nft_pages": { "type": "integer" }
	},
	"required": ["listen_addr", "collection_cache_ttl", "ownership_cache_ttl", "max_nft_pages"]
}
`)

type ServerConfig struct {
	ListenAddr         string `json:"listen_addr"`
	OwnershipCacheTTL  int    `json:"ownership_cache_ttl"`
	CollectionCacheTTL int    `json:"collection_cache_ttl"`
	MaxNFTPages        int    `json:"max_nft_pages"`
}
