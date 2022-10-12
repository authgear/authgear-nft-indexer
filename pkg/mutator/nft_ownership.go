package mutator

import (
	"context"
	"net/url"

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
	ids := make(map[string]bool)
	for _, ownership := range ownerships {
		ownerID, err := authgearweb3.NewContractID(
			ownership.Blockchain,
			ownership.Network,
			ownership.OwnerAddress,
			url.Values{},
		)
		if err != nil {
			return err
		}

		ownerURL, err := ownerID.URL()
		if err != nil {
			return err
		}

		if _, value := ids[ownerURL.String()]; !value {
			ids[ownerURL.String()] = true
			uniqueOwners = append(uniqueOwners, database.NFTOwner{
				Blockchain:   ownerID.Blockchain,
				Network:      ownerID.Network,
				Address:      ownerID.Address,
				LastSyncedAt: database.NewTimestamp(),
			})
		}

	}

	err := q.Session.RunInTx(q.Ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		// Inserting one by one to avoid conflict within the same insert statement
		for i := range ownerships {
			_, err := tx.NewInsert().
				Model(&ownerships[i]).
				On("CONFLICT (blockchain, network, contract_address, token_id) DO UPDATE").
				Set("balance = EXCLUDED.balance, block_number = EXCLUDED.block_number, owner_address = EXCLUDED.owner_address, txn_hash = EXCLUDED.txn_hash, txn_index = EXCLUDED.txn_index, block_timestamp = EXCLUDED.block_timestamp, updated_at = NOW()").
				Where("EXCLUDED.block_number > nft_ownership.block_number OR (EXCLUDED.block_number = nft_ownership.block_number AND EXCLUDED.txn_hash = nft_ownership.txn_hash AND EXCLUDED.txn_index > nft_ownership.txn_index)").
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
