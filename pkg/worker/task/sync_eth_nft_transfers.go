package task

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	apimodel "github.com/authgear/authgear-nft-indexer/pkg/api/model"
	"github.com/authgear/authgear-nft-indexer/pkg/config"
	ethmodel "github.com/authgear/authgear-nft-indexer/pkg/model/eth"
	"github.com/authgear/authgear-server/pkg/util/hexstring"
	authgearweb3 "github.com/authgear/authgear-server/pkg/util/web3"
	"github.com/jrallison/go-workers"
	"github.com/uptrace/bun/extra/bunbig"
)

const TransferPageSize = 1000

type SycnNFTTransferTransferMutator interface {
	InsertNFTTransfers(transfers []ethmodel.NFTTransfer) error
}

type SycnNFTTransferTaskCollectionQuery interface {
	QueryAllNFTCollections() ([]ethmodel.NFTCollection, error)
}

type SycnNFTTransferAlchemyAPI interface {
	GetNFTTransfers(blockchain string, network string, contractAddresses []string, fromBlock string, toBlock string, pageKey string, maxCount int64) (*apimodel.AssetTransferResponse, error)
}

type SyncETHNFTTransfersMessageArgs struct {
	Blockchain        string   `json:"blockchain"`
	Network           string   `json:"network"`
	ContractAddresses []string `json:"contract_address"`
	FromBlock         *big.Int `json:"from_block"`
	PageKey           string   `json:"page_key"`
}

type SyncETHNFTTransferTaskHandler struct {
	AlchemyAPI         SycnNFTTransferAlchemyAPI
	Config             config.Config
	NftTransferMutator SycnNFTTransferTransferMutator
	NftCollectionQuery SycnNFTTransferTaskCollectionQuery
}

func (h *SyncETHNFTTransferTaskHandler) Handler(message *workers.Msg) {
	args := message.Args()
	if args == nil {
		panic("SyncNFTTranfers: missing args")
	}

	argsJSON, err := args.Json.Encode()
	if err != nil {
		panic(fmt.Errorf("SyncNFTTranfers: failed to serialize args: %w", err))
	}

	var castedArgs SyncETHNFTTransfersMessageArgs
	err = json.Unmarshal(argsJSON, &castedArgs)
	if err != nil {
		panic(fmt.Errorf("SyncNFTTranfers: failed to unmarshal args: %w", err))
	}

	collections, err := h.NftCollectionQuery.QueryAllNFTCollections()
	if err != nil {
		panic(fmt.Errorf("SyncNFTTranfers: failed to get NFT collections: %w", err))
	}

	contractIDToCollectionMap := make(map[authgearweb3.ContractID]ethmodel.NFTCollection)
	for _, collection := range collections {

		contractID := authgearweb3.ContractID{
			Blockchain:      collection.Blockchain,
			Network:         collection.Network,
			ContractAddress: collection.ContractAddress,
		}

		contractIDToCollectionMap[contractID] = collection
	}

	validAddresses := make([]string, 0)
	for _, contractAddress := range castedArgs.ContractAddresses {
		contractID := authgearweb3.ContractID{
			Blockchain:      castedArgs.Blockchain,
			Network:         castedArgs.Network,
			ContractAddress: contractAddress,
		}

		collection, ok := contractIDToCollectionMap[contractID]
		// collection is not being watched
		if !ok {
			continue
		}

		collectionBlockHeight := collection.FromBlockHeight.ToMathBig()

		// Check whether the collection is synced or being synced
		// Check whether the collection height is larger than the request height
		// Check whether the task is an on-going task, i.e not from page 1
		if collectionBlockHeight.Cmp(big.NewInt(0)) > 0 && collectionBlockHeight.Cmp(castedArgs.FromBlock) > 0 && castedArgs.PageKey == "" {
			continue
		}

		validAddresses = append(validAddresses, contractAddress)

	}

	// no-op if no valid address
	if len(validAddresses) == 0 {
		fmt.Println("SyncNFTCollections: No valid addresses, quiting task")
		return
	}

	fromBlockHex, err := hexstring.NewFromBigInt(castedArgs.FromBlock)
	if err != nil {
		panic(fmt.Errorf("SyncNFTCollections: failed to convert smallest block to hex string: %w", err))
	}

	res, err := h.AlchemyAPI.GetNFTTransfers(castedArgs.Blockchain, castedArgs.Network, validAddresses, fromBlockHex.String(), "latest", castedArgs.PageKey, TransferPageSize)
	if err != nil {
		panic(fmt.Errorf("SyncNFTTranfers: failed to get NFT transfers: %w", err))
	}

	nftTransfers := make([]ethmodel.NFTTransfer, 0, len(res.Result.Transfers))
	for _, transfer := range res.Result.Transfers {

		tokenID := hexstring.MustParse(transfer.TokenID)
		blockNum := hexstring.MustParse(transfer.BlockNum)

		blockTime, err := time.Parse(time.RFC3339, transfer.Metadata.BlockTimestamp)
		if err != nil {
			panic(fmt.Errorf("SyncNFTTranfers: failed to parse block time: %w", err))
		}

		nftTransfers = append(nftTransfers, ethmodel.NFTTransfer{
			Blockchain:      castedArgs.Blockchain,
			Network:         castedArgs.Network,
			ContractAddress: transfer.RawContract.Address,
			TokenID:         bunbig.FromMathBig(tokenID.ToBigInt()),
			BlockNumber:     bunbig.FromMathBig(blockNum.ToBigInt()),
			FromAddress:     transfer.From,
			ToAddress:       transfer.To,
			TransactionHash: transfer.Hash,
			BlockTimestamp:  blockTime,
		})
	}

	// Rico: NFT Owners are automatically updated when we insert new records to the NFT transfers table.
	// While we know the sequential ordering for the NFT Transfers based on the block number,
	// we can't be sure that a single NFT token wouldn't be transferred multiple times in a single block.
	// Given the low chance of it happening and the fact that it will be automatically resolve if the token is transferred again,
	// we will skip it for now.
	err = h.NftTransferMutator.InsertNFTTransfers(nftTransfers)
	if err != nil {
		panic(fmt.Errorf("SyncNFTTranfers: failed to insert NFT transfers: %w", err))
	}

	if res.Result.PageKey != "" {
		_, err := workers.Enqueue(h.Config.Worker.TransferQueueName, "Sync", SyncETHNFTTransfersMessageArgs{
			Blockchain:        castedArgs.Blockchain,
			Network:           castedArgs.Network,
			ContractAddresses: castedArgs.ContractAddresses,
			FromBlock:         castedArgs.FromBlock,
			PageKey:           res.Result.PageKey,
		})

		if err != nil {
			panic(fmt.Errorf("SyncNFTTranfers: failed to enqueue NFT transfers: %w", err))
		}
	}
}
