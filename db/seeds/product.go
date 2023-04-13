package seeds

import (
	"gorm.io/gorm"
	"package-service/models"
)

func CreateProducts(db *gorm.DB, products []models.Product) error {
	for _, product := range products {
		if err := db.Where("gtin = ?", product.Gtin).FirstOrCreate(&product).Error; err != nil {
			return err
		}
	}

	return nil
}
