package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(dsn string) (*Repository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}

func (r *Repository) GetLastIndexedBlock() (int64, error) {
	var height int64
	err := r.db.QueryRow(`SELECT height FROM blocks ORDER BY height DESC LIMIT 1`).Scan(&height)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return height, err
}

func (r *Repository) SaveBlock(height int64) error {
	_, err := r.db.Exec(`INSERT INTO blocks(height,time) VALUES($1,NOW())`, height)
	return err
}

func (r *Repository) SaveNFTTransfer(txHash string, blockHeight int64, tokenId, from, to string) error {
	_, err := r.db.Exec(`INSERT INTO nft_transfers(block_height, tx_hash, token_id, from_addr, to_addr, timestamp)
                         VALUES($1, $2, $3, $4, $5, NOW())`,
		blockHeight, txHash, tokenId, from, to)
	return err
}
