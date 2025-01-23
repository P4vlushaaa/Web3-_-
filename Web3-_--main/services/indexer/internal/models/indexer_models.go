package models

type BlockRecord struct {
	Height int64  `json:"height"`
	Time   string `json:"time"`
}

type NftTransferRecord struct {
	ID          int64  `json:"id"`
	BlockHeight int64  `json:"block_height"`
	TxHash      string `json:"tx_hash"`
	TokenId     string `json:"token_id"`
	FromAddr    string `json:"from_addr"`
	ToAddr      string `json:"to_addr"`
	Timestamp   string `json:"timestamp"`
}
