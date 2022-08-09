package config

var _ = Schema.Add("WorkerConfig", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"collection_queue_name": { "type": "string" },
		"transfer_queue_name": { "type": "string" }
	},
	"required": ["collection_queue_name", "transfer_queue_name"]
}
`)

type WorkerConfig struct {
	CollectionQueueName string `json:"collection_queue_name"`
	TransferQueueName   string `json:"transfer_queue_name"`
}
