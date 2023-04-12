package models

import "time"

type Box struct {
	Id   int    `json:"id" gorm:"primaryKey"`
	Sscc string `json:"sscc" gorm:"uniqueIndex,size:18"`

	Packages []Package `json:"-"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime:false"`
	UpdatedAt time.Time `json:"updated_at"`
}
