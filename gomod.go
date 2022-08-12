package gomod

// If we do not list "github.com/google/wire/cmd/wire" here
// There will be no "github.com/google/subcommands"
// "github.com/google/subcommands" is required to make `make generate` work.
import (
	_ "github.com/google/wire/cmd/wire"
)
