package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
	"os"
	"package-service/db/seeds"
	"package-service/models"
)

var DB *gorm.DB

func Connect() error {
	dbHost, hostExists := os.LookupEnv("DB_HOST")
	if !hostExists {
		panic("DB host address not exists")
	}

	dbPort, portExists := os.LookupEnv("DB_PORT")
	if !portExists {
		panic("DB port not exists")
	}

	dbName, nameExists := os.LookupEnv("DB_DATABASE")
	if !nameExists {
		panic("DB database name not exists")
	}

	dbUsername, usernameExists := os.LookupEnv("DB_USERNAME")
	if !usernameExists {
		panic("DB username not exists")
	}

	dbPassword, passwordExists := os.LookupEnv("DB_PASSWORD")
	if !passwordExists {
		panic("DB password not exists")
	}

	dsn := dbUsername + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=true"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	database.Use(prometheus.New(prometheus.Config{
		DBName:          dbName,
		RefreshInterval: 15,
		StartServer:     true,
		HTTPServerPort:  8080,
	}))

	if err != nil {
		return err
	}

	if err := database.AutoMigrate(&models.Box{}, &models.Product{}, &models.Package{}); err != nil {
		return err
	}

	for _, seed := range seeds.All() {
		if err := seed.Run(database); err != nil {
			return fmt.Errorf("running seed '%s', failed with error: %s", seed.Name, err)
		}
	}

	DB = database

	return nil
}
