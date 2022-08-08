package model

import "github.com/jrallison/go-workers"

type Task interface {
	Handler(message *workers.Msg)
}
