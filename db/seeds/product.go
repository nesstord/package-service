package seeds

import (
	"gorm.io/gorm"
	"package-service/models"
)

func CreateProducts(db *gorm.DB, products []models.Product) error {
	for _, product := range products {
		product := product
		if err := db.FirstOrCreate(&product).Error; err != nil {
			return err
		}
	}

	return nil
}
