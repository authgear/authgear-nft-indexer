-- +migrate Up

CREATE TABLE eth_nft_collection
(
	id text PRIMARY KEY,
	contract_address text NOT NULL,
	name text NOT NULL,
	blockchain text NOT NULL,
	network text NOT NULL,
	from_block_height bigint NOT NULL,
	created_at    timestamp without time zone NOT NULL,
    updated_at    timestamp without time zone NOT NULL
);

CREATE UNIQUE INDEX eth_nft_collection_unq_collection_idx ON eth_nft_collection (blockchain, network, contract_address);

CREATE TABLE eth_nft_transfer
(
	blockchain text NOT NULL,
	network text NOT NULL,
	contract_address text NOT NULL,
	token_id bigint NOT NULL,
	block_number bigint NOT NULL,
	from_address text NOT NULL,
	to_address text NOT NULL,
	txn_hash text NOT NULL,
	block_timestamp timestamp without time zone NOT NULL,
	created_at    timestamp without time zone NOT NULL,
	updated_at    timestamp without time zone NOT NULL
);

CREATE UNIQUE INDEX eth_nft_transfer_unq_transfer_idx ON eth_nft_transfer (blockchain, network, contract_address, token_id, from_address, to_address, txn_hash);

CREATE TABLE eth_nft_owner
(
	blockchain text NOT NULL,
	network text NOT NULL,
	contract_address text NOT NULL,
	token_id bigint NOT NULL,
	block_number bigint NOT NULL,
	owner_address text NOT NULL,
	txn_hash text NOT NULL,
	block_timestamp timestamp without time zone NOT NULL,
	created_at    timestamp without time zone NOT NULL,
	updated_at    timestamp without time zone NOT NULL
);

CREATE UNIQUE INDEX eth_nft_owner_unq_owner_idx ON eth_nft_owner (blockchain, network, contract_address, token_id);
CREATE INDEX eth_nft_owner_owned_idx ON eth_nft_owner (blockchain, network, owner_address);

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION eth_update_nft_owner() RETURNS TRIGGER AS $$
	BEGIN
		IF (TG_OP = 'INSERT') THEN
			-- Insert all new record to the db
			INSERT INTO eth_nft_owner (blockchain, network, contract_address, token_id, block_number, owner_address, txn_hash, block_timestamp, created_at, updated_at)
			VALUES (NEW.blockchain, NEW.network, NEW.contract_address, NEW.token_id, NEW.block_number, NEW.to_address, NEW.txn_hash, NEW.block_timestamp, NOW(), NOW())
			ON CONFLICT (blockchain, network, contract_address, token_id) DO NOTHING;

			-- Update owners
			UPDATE eth_nft_owner
			SET owner_address = NEW.to_address,
				block_number = NEW.block_number,
				txn_hash = NEW.txn_hash,
				updated_at = NOW()
			WHERE blockchain = NEW.blockchain AND network = NEW.network AND contract_address = NEW.contract_address AND token_id = NEW.token_id AND (
				-- Update the owner anyway if it's a new block
				NEW.block_number > eth_nft_owner.block_number OR
				-- Token is being transferred within the same trasaction
				(NEW.from_address = eth_nft_owner.owner_address)
			);

			-- Update collection db
			UPDATE eth_nft_collection
			-- Minus 1 to ensure the synchronization is not terminated in the middle of a block
			SET from_block_height = NEW.block_number - 1
			WHERE network = NEW.network AND contract_address = NEW.contract_address AND NEW.block_number - 1 >= 0;
			
		END IF;
		RETURN NEW;
	END;
$$ language 'plpgsql';
-- +migrate StatementEnd

CREATE TRIGGER eth_update_nft_owner_trigger
AFTER INSERT ON eth_nft_transfer
FOR EACH ROW EXECUTE FUNCTION eth_update_nft_owner();



-- +migrate Down
DROP INDEX eth_nft_owner_unq_owner_idx;
DROP INDEX eth_nft_owner_owned_idx;
DROP TABLE eth_nft_owner;
DROP INDEX eth_nft_transfer_unq_transfer_idx;
DROP TABLE eth_nft_transfer;
DROP INDEX eth_nft_collection_unq_collection_idx;
DROP TABLE eth_nft_collection;
DROP FUNCTION eth_update_nft_owner();
