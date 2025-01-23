package models

type NFTProperties struct {
	Owner      string `json:"owner"`
	Name       string `json:"name"`
	BlurredRef string `json:"blurred_ref"`
	FullRef    string `json:"full_ref"`
	Price      int    `json:"price"`
}
