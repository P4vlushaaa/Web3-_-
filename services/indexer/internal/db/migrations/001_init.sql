CREATE TABLE IF NOT EXISTS blocks (
                                      height BIGINT PRIMARY KEY,
                                      time TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS nft_transfers (
                                             id SERIAL PRIMARY KEY,
                                             block_height BIGINT NOT NULL,
                                             tx_hash VARCHAR(66) NOT NULL,
    token_id TEXT NOT NULL,
    from_addr VARCHAR(42),
    to_addr VARCHAR(42),
    timestamp TIMESTAMP NOT NULL
    );