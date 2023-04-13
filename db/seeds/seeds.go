package seeds

import (
	"gorm.io/gorm"
	"package-service/models"
)

func All() []Seed {
	return []Seed{
		{
			Name: "CreateProducts",
			Run: func(db *gorm.DB) error {
				return CreateProducts(db, []models.Product{
					{
						Name:  "Продукт1",
						Gtin:  "04603988000001",
						Packs: 4,
					},
					{
						Name:  "Продукт2",
						Gtin:  "04603988000002",
						Packs: 48,
					},
					{
						Name:  "Продукт3",
						Gtin:  "04603988000003",
						Packs: 180,
					},
					{
						Name:  "Продукт4",
						Gtin:  "04603988000004",
						Packs: 360,
					},
					{
						Name:  "Продукт5",
						Gtin:  "04603988000005",
						Packs: 18,
					},
					{
						Name:  "Продукт6",
						Gtin:  "04603988000006",
						Packs: 36,
					},
					{
						Name:  "Продукт7",
						Gtin:  "04603988000007",
						Packs: 15,
					},
					{
						Name:  "Продукт8",
						Gtin:  "04603988000008",
						Packs: 90,
					},
					{
						Name:  "Продукт9",
						Gtin:  "04603988000009",
						Packs: 144,
					},
				})
			},
		},
	}
}
