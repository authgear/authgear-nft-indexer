-- +migrate Up
ALTER TABLE eth_nft_collection ALTER COLUMN total_supply DROP NOT NULL;

-- +migrate Down
ALTER TABLE eth_nft_collection ALTER COLUMN total_supply SET NOT NULL;