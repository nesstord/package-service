package services

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"io"
	"net/http"
	"package-service/db"
	"package-service/http/requests"
	"package-service/http/responses"
	"package-service/http/validators"
	"package-service/models"
	"package-service/services/exceptions"
	"time"
)

func BoxAggregate(ctx context.Context, request requests.AggregateRequest) responses.BasicResponse {
	tx := db.DB.Begin()

	ch := make(chan responses.BasicResponse)
	go func(response chan<- responses.BasicResponse) {
		//1. Проверяем, есть ли уже такая коробка

		created, _ := time.Parse(time.RFC3339, request.Created)
		box := models.Box{
			Sscc:      request.Sscc,
			CreatedAt: created,
		}
		err := tx.Debug().Where("sscc = ?", request.Sscc).First(&box).Error

		if err == nil {
			response <- responses.Error(exceptions.BoxErrorAlreadyUsedSscc, 400)
			return
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			response <- responses.Error(err.Error(), 500)
			return
		}

		firstGtin := validators.GtinRegexp.FindStringSubmatch(request.Sgtins[0])[0]

		//5. Проверяем товар
		product := models.Product{
			Gtin: firstGtin,
		}
		if err := tx.Debug().Where("gtin = ?", firstGtin).First(&product).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				response <- responses.Error(exceptions.BoxErrorUnknownGtin, 400)
				return
			} else {
				response <- responses.Error(err.Error(), 500)
				return
			}
		}

		//4. Проверяем количество упаковок, доступных для коробки
		if product.Packs != len(request.Sgtins) {
			response <- responses.Error(exceptions.BoxErrorInvalidPackagesNumber, 400)
			return
		}

		for _, sgtin := range request.Sgtins {
			//3. Проверяем, одинаковые ли GTIN
			if gtin := validators.GtinRegexp.FindStringSubmatch(sgtin); gtin[0] != firstGtin {
				response <- responses.Error(exceptions.BoxErrorDifferentGtins, 400)
				return
			}

			//2. Проверяем, была ли такая упаковка в коробке
			pack := models.Package{}
			err := tx.Debug().Where("sgtin = ?", sgtin).First(&pack).Error

			if err == nil {
				response <- responses.Error(exceptions.BoxErrorAlreadyUsedSgtin, 400)
				return
			} else if !errors.Is(err, gorm.ErrRecordNotFound) {
				response <- responses.Error(err.Error(), 500)
				return
			}

			box.Packages = append(box.Packages, models.Package{
				Sgtin:     sgtin,
				ProductId: product.Id,
			})
		}

		tx.Create(&box)

		response <- responses.Ok()
	}(ch)

	select {
	case <-ctx.Done():
		tx.Rollback()
		return responses.Error("отмена контекстом", http.StatusBadRequest)
	case res := <-ch:
		if res.Ok {
			tx.Commit()
		} else {
			tx.Rollback()
		}

		return res
	}
}

func BoxGetBySgtins(request requests.GetBoxesBySgtinsRequest) (map[string]interface{}, error) {
	m := make(map[string]interface{}, len(request.Sgtins))
	for _, sgtin := range request.Sgtins {
		m[sgtin] = nil
	}

	packages := make([]models.Package, 0)
	if err := db.DB.Debug().Preload("Box").Where("sgtin IN ?", request.Sgtins).Find(&packages).Error; err != nil {
		return nil, err
	}

	for _, pack := range packages {
		m[pack.Sgtin] = pack.Box.Sscc
	}

	return m, nil
}

func BoxGetByGtin(request requests.GetBoxesByGtinRequest, writer io.Writer) error {
	packages := make([]models.Package, 0)
	if err := db.DB.Debug().InnerJoins("Product").Preload("Box").Where("gtin = ?", request.Gtin).Find(&packages).Error; err != nil {
		return fmt.Errorf("ошибка получения упаковок: %s", err.Error())
	}

	records := make([][]string, 0)
	records = append(records, []string{"sgtin", "sscc"})

	for _, pack := range packages {
		records = append(records, []string{pack.Sgtin, pack.Box.Sscc})
	}

	w := csv.NewWriter(writer)
	if err := w.WriteAll(records); err != nil {
		return fmt.Errorf("ошибка записи в буффер: %s", err.Error())
	}

	return nil
}
