package models

type Product struct {
	Name  string `json:"name"`
	Gtin  string `json:"gtin"`
	Packs int    `json:"packs"`
}
