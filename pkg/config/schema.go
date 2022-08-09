package config

import (
	"github.com/authgear/authgear-server/pkg/util/validation"
)

var Schema = validation.NewMultipartSchema("Config")

func init() {
	Schema.Instantiate()
}
