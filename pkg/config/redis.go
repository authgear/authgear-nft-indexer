package config

import (
	"net/url"
)

var _ = Schema.Add("RedisConfig", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"url": { "type": "string" },
		"max_open_connection": { "type": "integer" },
		"max_idle_connection": { "type": "integer" },
		"max_conn_lifetime": { "type": "integer" },
		"idle_conn_timeout": { "type": "integer" }
	},
	"required": ["url", "max_open_connection", "max_idle_connection", "max_conn_lifetime", "idle_conn_timeout"]
}
`)

type RedisConfig struct {
	URL                   string `json:"url"`
	MaxOpenConnection     int    `json:"max_open_connection"`
	MaxIdleConnection     int    `json:"max_idle_connection"`
	MaxConnectionLifeTime int    `json:"max_conn_lifetime"`
	IdleConnectionTimeout int    `json:"idle_conn_timeout"`
}

func (c *RedisConfig) RedisURL() (*url.URL, error) {
	return url.Parse(c.URL)
}
