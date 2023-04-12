package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"package-service/db"
	"package-service/http/requests"
	"package-service/http/responses"
	"package-service/http/validators"
	"package-service/models"
	"regexp"
	"strconv"
	"time"
)

func Aggregate(c *gin.Context) {
	request := requests.AggregateRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error(err.Error(), 400))
		return
	}

	//1. Проверяем, есть ли уже такая коробка
	created, _ := time.Parse(time.RFC3339, request.Created)
	box := models.Box{
		Sscc:      request.Sscc,
		CreatedAt: created,
	}
	err := db.DB.Debug().Where("sscc = ?", request.Sscc).First(&box).Error

	if err == nil {
		c.JSON(http.StatusBadRequest, responses.Error("SSCC уже использован", 400))
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error(), 500))
		return
	}

	re := regexp.MustCompile(validators.GtinPattern)
	firstGtin := re.FindStringSubmatch(request.Sgtins[0])[0]

	//5. Проверяем товар
	product := models.Product{
		Gtin: firstGtin,
	}
	if err := db.DB.Debug().Where("gtin = ?", firstGtin).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, responses.Error("Неизвестный GTIN", 400))
			return
		} else {
			c.JSON(http.StatusInternalServerError, responses.Error(err.Error(), 500))
			return
		}
	}

	//4. Проверяем количество упаковок, доступных для коробки
	if product.Packs != len(request.Sgtins) {
		n := strconv.Itoa(product.Packs)
		c.JSON(http.StatusBadRequest, responses.Error("Для данного товара в коробке может быть только "+n+" пачек", 400))
		return
	}

	for _, sgtin := range request.Sgtins {
		//3. Проверяем, одинаковые ли GTIN
		if gtin := re.FindStringSubmatch(sgtin); gtin[0] != firstGtin {
			c.JSON(http.StatusBadRequest, responses.Error("Запрос содержит разные GTIN", 400))
			return
		}

		//2. Проверяем, была ли такая упаковка в коробке
		pack := models.Package{}
		err := db.DB.Debug().Where("sgtin = ?", sgtin).First(&pack).Error

		if err == nil {
			c.JSON(http.StatusBadRequest, responses.Error("SGTIN уже использован", 400))
			return
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, responses.Error(err.Error(), 500))
			return
		}

		box.Packages = append(box.Packages, models.Package{
			Sgtin:     sgtin,
			ProductId: product.Id,
		})
	}

	db.DB.Create(&box)

	c.JSON(http.StatusOK, responses.Ok())
}
