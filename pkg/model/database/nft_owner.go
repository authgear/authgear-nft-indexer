package database

import (
	"time"

	"github.com/uptrace/bun"
)

type NFTOwner struct {
	bun.BaseModel `bun:"table:nft_owner"`

	Blockchain   string    `bun:"blockchain,notnull"`
	Network      string    `bun:"network,notnull"`
	Address      string    `bun:"address,notnull"`
	LastSyncedAt time.Time `bun:"last_synced_at,notnull"`
}
