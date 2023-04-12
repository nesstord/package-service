package models

import "time"

type Product struct {
	Id    int    `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Gtin  string `json:"gtin" gorm:"uniqueIndex,size:14"`
	Packs int    `json:"packs"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
