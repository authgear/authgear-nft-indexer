package eth

import (
	"fmt"
	"strings"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/extra/bunbig"

	"github.com/authgear/authgear-nft-indexer/pkg/model"
)

type NFTCollectionType string

const (
	NFTCollectionTypeERC721 NFTCollectionType = "erc721"
)

func ParseNFTCollectionType(t string) (NFTCollectionType, error) {
	tokenType := strings.ToLower(t)
	switch tokenType {
	case "erc721":
		return NFTCollectionTypeERC721, nil
	default:
		return "", fmt.Errorf("unknown nft collection type: %+v", tokenType)
	}
}

type NFTCollection struct {
	bun.BaseModel `bun:"table:eth_nft_collection"`
	model.BaseWithID

	Blockchain      string            `bun:"blockchain,notnull"`
	Network         string            `bun:"network,notnull"`
	ContractAddress string            `bun:"contract_address,notnull"`
	Name            string            `bun:"name,notnull"`
	FromBlockHeight *bunbig.Int       `bun:"from_block_height,notnull"`
	TotalSupply     *bunbig.Int       `bun:"total_supply,notnull"`
	Type            NFTCollectionType `bun:"type,notnull"`
}
