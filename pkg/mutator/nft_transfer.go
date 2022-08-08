package mutator

import (
	"context"

	"github.com/authgear/authgear-nft-indexer/pkg/model/eth"
	"github.com/uptrace/bun"
)

type NFTTransferMutator struct {
	Ctx     context.Context
	Session *bun.DB
}

func (q *NFTTransferMutator) InsertNFTTransfers(transfers []eth.NFTTransfer) error {
	err := q.Session.RunInTx(q.Ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().
			Model(&transfers).
			On("CONFLICT (blockchain, network, contract_address, token_id, from_address, to_address, txn_hash) DO NOTHING").
			Exec(ctx)
		return err
	})

	return err
}
