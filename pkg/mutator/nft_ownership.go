package mutator

import (
	"context"

	"github.com/authgear/authgear-nft-indexer/pkg/model/database"
	authgearweb3 "github.com/authgear/authgear-server/pkg/util/web3"
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

	uniqueOwners := make([]database.NFTOwner, 0)
	ids := make(map[authgearweb3.ContractID]bool)
	for _, ownership := range ownerships {
		contractID := authgearweb3.ContractID{
			Blockchain:      ownership.Blockchain,
			Network:         ownership.Network,
			ContractAddress: ownership.OwnerAddress,
		}

		if _, value := ids[contractID]; !value {
			ids[contractID] = true
			uniqueOwners = append(uniqueOwners, database.NFTOwner{
				Blockchain:   contractID.Blockchain,
				Network:      contractID.Network,
				Address:      contractID.ContractAddress,
				LastSyncedAt: database.NewTimestamp(),
			})
		}

	}

	err := q.Session.RunInTx(q.Ctx, nil, func(ctx context.Context, tx bun.Tx) error {

		// Insert individually to prevent conflict within the values
		for i := range ownerships {
			_, err := tx.NewInsert().
				Model(&ownerships[i]).
				On("CONFLICT (blockchain, network, contract_address, token_id) DO UPDATE").
				Set("txn_hash = EXCLUDED.txn_hash, block_number = EXCLUDED.block_number, block_timestamp = EXCLUDED.block_timestamp").
				Returning("*").
				Exec(ctx)
			if err != nil {
				return err
			}
		}
		_, err := tx.NewInsert().
			Model(&uniqueOwners).
			On("CONFLICT (blockchain, network, address) DO UPDATE").
			Set("last_synced_at = EXCLUDED.last_synced_at").
			Returning("*").
			Exec(ctx)

		return err
	})

	return err
}
