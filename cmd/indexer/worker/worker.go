package worker

import (
	"context"
	"log"
	"time"

	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/authgear/authgear-nft-indexer/pkg/database"
	"github.com/authgear/authgear-nft-indexer/pkg/worker"
	"github.com/go-co-op/gocron"
	"github.com/jrallison/go-workers"
)

type Worker struct {
	ctx    context.Context
	config config.Config
}

func (w *Worker) Start() {
	worker.ConfigureWorkers(w.config.Redis)

	db := database.GetDatabase(w.config.Database)

	syncCollectionHandler := NewSyncETHNFTCollectionTaskHandler(w.ctx, w.config, db)
	workers.Process(w.config.Worker.CollectionQueueName, syncCollectionHandler.Handler, 1)
	syncTransferHandler := NewSyncETHNFTTransferTaskHandler(w.ctx, w.config, db)
	workers.Process(w.config.Worker.TransferQueueName, syncTransferHandler.Handler, 1)

	scheduler := gocron.NewScheduler(time.UTC)
	if _, err := scheduler.Every(5).Minutes().Do(func() {
		if _, err := workers.Enqueue(w.config.Worker.CollectionQueueName, "", nil); err != nil {
			panic(err)
		}
	}); err != nil {
		log.Fatalf("failed to schedule job: %+v", err)
	}

	scheduler.StartAsync()
	log.Printf("started scheduler")
	workers.Run()
}
