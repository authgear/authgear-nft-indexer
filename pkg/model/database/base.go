package database

import (
	"context"
	"crypto/rand"
	"strings"
	"time"

	ulid "github.com/oklog/ulid/v2"
	"github.com/uptrace/bun"
)

type Base struct {
	CreatedAt time.Time `bun:"created_at,notnull"`
	UpdatedAt time.Time `bun:"updated_at,notnull"`
}

type BaseWithID struct {
	ID string `bun:"id,notnull"`
	Base
}

func NewTimestamp() time.Time {
	return time.Now().UTC().Truncate(1000 * time.Nanosecond)
}

func NewID() (string, error) {
	newID, err := ulid.New(ulid.Now(), rand.Reader)
	if err != nil {
		return "", err
	}
	return strings.ToLower(newID.String()), nil
}

var _ bun.BeforeAppendModelHook = (*Base)(nil)
var _ bun.BeforeAppendModelHook = (*BaseWithID)(nil)

func (base *Base) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	now := NewTimestamp()
	switch query.(type) {
	case *bun.InsertQuery:
		base.CreatedAt = now
		base.UpdatedAt = now
	case *bun.UpdateQuery:
		base.UpdatedAt = now
	}
	return nil
}

func (base *BaseWithID) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	if len(base.ID) == 0 {
		newID, err := NewID()
		if err != nil {
			return err
		}
		base.ID = newID
	}

	now := NewTimestamp()
	switch query.(type) {
	case *bun.InsertQuery:
		base.CreatedAt = now
		base.UpdatedAt = now
	case *bun.UpdateQuery:
		base.UpdatedAt = now
	}
	return nil
}
