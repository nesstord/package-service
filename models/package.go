package models

import "time"

type Package struct {
	Id    int    `json:"id" gorm:"primaryKey"`
	Sgtin string `json:"sgtin" gorm:"uniqueIndex,size:27"`

	BoxId     int     `json:"box_id" gorm:"foreignKey"`
	Box       Box     `json:"-"`
	ProductId int     `json:"product_id" gorm:"foreignKey"`
	Product   Product `json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
