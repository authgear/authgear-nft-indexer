-- +migrate Up

CREATE TABLE eth_nft_collection_probe
(
	blockchain text NOT NULL,
	network text NOT NULL,
	contract_address text NOT NULL,
	is_large_collection boolean NOT NULL
);

CREATE UNIQUE INDEX eth_nft_collection_probe_unq_collection_idx ON eth_nft_collection_probe (blockchain, network, contract_address);

CREATE TABLE eth_nft_collection
(
	id text PRIMARY KEY,
	blockchain text NOT NULL,
	network text NOT NULL,
	contract_address text NOT NULL,
	name text NOT NULL,
	total_supply bigint,
	type text NOT NULL,
	created_at timestamp without time zone NOT NULL,
	updated_at timestamp without time zone NOT NULL
);

CREATE UNIQUE INDEX eth_nft_collection_unq_collection_idx ON eth_nft_collection (blockchain, network, contract_address);

CREATE TABLE eth_nft_ownership
(
	blockchain text NOT NULL,
	network text NOT NULL,
	contract_address text NOT NULL,
	token_id text NOT NULL,
	balance text NOT NULL,
	block_number bigint NOT NULL,
	owner_address text NOT NULL,
	txn_hash text NOT NULL,
	txn_index integer NOT NULL,
	block_timestamp timestamp without time zone NOT NULL,
	created_at timestamp without time zone NOT NULL,
	updated_at timestamp without time zone NOT NULL
);

CREATE UNIQUE INDEX eth_nft_ownership_unq_owner_idx ON eth_nft_ownership (blockchain, network, contract_address, token_id);
CREATE INDEX eth_nft_ownership_owned_idx ON eth_nft_ownership (blockchain, network, owner_address);

CREATE TABLE eth_nft_owner
(
	blockchain text NOT NULL,
	network text NOT NULL,
	address text NOT NULL,
	last_synced_at timestamp without time zone NOT NULL
);

CREATE UNIQUE INDEX eth_nft_owner_unq_owner_idx ON eth_nft_owner (blockchain, network, address);

-- +migrate Down
DROP TABLE eth_nft_owner;
DROP TABLE eth_nft_ownership;
DROP TABLE eth_nft_collection;
DROP TABLE eth_nft_collection_probe;
