//go:build tools

package gomod

// As of go1.21, we need to build flag to enable this trick.
// See https://github.com/golang/go/issues/48429
// and https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
//
// If we do not list "github.com/google/wire/cmd/wire" here
// There will be no "github.com/google/subcommands"
// "github.com/google/subcommands" is required to make `make generate` work.
import (
	_ "github.com/google/wire/cmd/wire"
)
