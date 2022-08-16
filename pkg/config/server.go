package config

var _ = Schema.Add("ServerConfig", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"listen_addr": { "type": "string" }
	},
	"required": ["listen_addr"]
}
`)

type ServerConfig struct {
	ListenAddr string `json:"listen_addr"`
}
