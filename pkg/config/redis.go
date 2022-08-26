package config

import "net/url"

var _ = Schema.Add("RedisConfig", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"server": { "type": "string" },
		"database": { "type": "string" },
		"max_open_connection": { "type": "integer" },
		"max_idle_connection": { "type": "integer" },
		"max_conn_lifetime": { "type": "integer" },
		"idle_conn_timeout": { "type": "integer" }
	},
	"required": ["server", "database", "max_open_connection", "max_idle_connection", "max_conn_lifetime", "idle_conn_timeout"]
}
`)

type RedisConfig struct {
	Server                string `json:"server"`
	Database              string `json:"database"`
	MaxOpenConnection     int    `json:"max_open_connection"`
	MaxIdleConnection     int    `json:"max_idle_connection"`
	MaxConnectionLifeTime int    `json:"max_conn_lifetime"`
	IdleConnectionTimeout int    `json:"idle_conn_timeout"`
}

func (rc *RedisConfig) RedisURL() string {
	url := &url.URL{}
	url.Scheme = "redis"
	url.Host = rc.Server
	url.Path = rc.Database

	return url.String()
}
