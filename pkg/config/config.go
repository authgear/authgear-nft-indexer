package config

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/authgear/authgear-server/pkg/util/validation"
	"sigs.k8s.io/yaml"
)

var _ = Schema.Add("Config", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"redis": { "$ref": "#/$defs/RedisConfig" },
		"worker": { "$ref": "#/$defs/WorkerConfig" },
		"database": { "$ref": "#/$defs/DatabaseConfig" },
		"alchemy": {
			"type": "array",
			"items": { "$ref": "#/$defs/AlchemyConfig" }
		}
	},
	"required": ["redis", "worker", "database", "alchemy"]
}
`)

type Config struct {
	Redis    RedisConfig     `json:"redis"`
	Worker   WorkerConfig    `json:"worker"`
	Database DatabaseConfig  `json:"database"`
	Alchemy  []AlchemyConfig `json:"alchemy"`
}

func Parse(inputYAML []byte) (*Config, error) {
	const validationErrorMessage = "invalid configuration"

	jsonData, err := yaml.YAMLToJSON(inputYAML)
	if err != nil {
		return nil, err
	}

	err = Schema.Validator().ValidateWithMessage(bytes.NewReader(jsonData), validationErrorMessage)
	if err != nil {
		return nil, err
	}

	var config Config
	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	err = validation.ValidateValueWithMessage(&config, validationErrorMessage)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func NewConfig(configPath string) Config {
	file, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	config, err := Parse(file)
	if err != nil {
		panic(err)
	}

	return *config
}
