package query

import (
	"context"

	"github.com/authgear/authgear-nft-indexer/pkg/model/eth"
	"github.com/uptrace/bun"
)

type NFTCollectionQuery struct {
	Ctx     context.Context
	Session *bun.DB
}

func (q *NFTCollectionQuery) QueryNFTCollections() ([]eth.NFTCollection, error) {
	nftCollections := make([]eth.NFTCollection, 0)

	err := q.Session.NewSelect().Model(&nftCollections).Scan(q.Ctx)
	if err != nil {
		return nil, err
	}

	return nftCollections, nil
}
