package mutator

import (
	"context"

	"github.com/authgear/authgear-nft-indexer/pkg/model/database"
	"github.com/uptrace/bun"
)

type NFTOwnershipMutator struct {
	Ctx     context.Context
	Session *bun.DB
}

func (q *NFTOwnershipMutator) InsertNFTOwnerships(ownerships []database.NFTOwnership) error {

	if len(ownerships) == 0 {
		return nil
	}

	err := q.Session.RunInTx(q.Ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().
			Model(&ownerships).
			Returning("*").
			Exec(ctx)

		return err
	})

	return err
}
