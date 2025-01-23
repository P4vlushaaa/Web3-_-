package models

type MarketItem struct {
	TokenID string `json:"token_id"`
	OnSale  bool   `json:"on_sale"`
}
