package mutator

import (
	"context"
	"database/sql"
	"math/big"

	"github.com/authgear/authgear-nft-indexer/pkg/model/database"
	authgearweb3 "github.com/authgear/authgear-server/pkg/util/web3"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/extra/bunbig"
)

type NFTCollectionMutator struct {
	Ctx     context.Context
	Session *bun.DB
}

func (q *NFTCollectionMutator) InsertNFTCollection(contractID authgearweb3.ContractID, contractName string, tokenType database.NFTCollectionType, totalSupply *big.Int) (*database.NFTCollection, error) {
	collection := &database.NFTCollection{
		Blockchain:      contractID.Blockchain,
		Network:         contractID.Network,
		ContractAddress: contractID.Address,
		Name:            contractName,
		TotalSupply:     bunbig.FromMathBig(totalSupply),
		Type:            tokenType,
	}

	err := q.Session.RunInTx(q.Ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		// Query record ID if exists, if not, generate new id and insert
		tx.NewSelect().
			Model(collection).
			Where("blockchain = ? AND network = ? AND contract_address = ?", collection.Blockchain, collection.Network, collection.ContractAddress).
			Limit(1).
			Scan(ctx) //nolint:errcheck

		_, err := tx.NewInsert().
			Model(collection).
			On("CONFLICT (blockchain, network, contract_address) DO UPDATE").
			Set("total_supply = EXCLUDED.total_supply, updated_at = NOW()").
			Returning("*").
			Exec(ctx)
		return err
	})

	if err != nil {
		return nil, err
	}

	return collection, nil
}

func (q *NFTCollectionMutator) DeleteNFTCollection(id string) (*database.NFTCollection, error) {

	collection := &database.NFTCollection{}

	err := q.Session.RunInTx(q.Ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		res, err := tx.NewDelete().
			Model(collection).
			Where("id = ?", id).
			Returning("*").
			Exec(ctx)

		if row, err := res.RowsAffected(); err != nil || row == 0 {
			return sql.ErrNoRows
		}

		return err
	})

	if err != nil {
		return nil, err
	}

	return collection, nil
}
