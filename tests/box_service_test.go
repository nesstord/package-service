package tests

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	_db "package-service/db"
	"package-service/http/requests"
	"package-service/services"
	"package-service/services/exceptions"
	"testing"
)

func TestBoxAggregationSuccess(t *testing.T) {
	request := requests.AggregateRequest{
		Sscc:    "000000000000000001",
		Created: "2022-12-01T06:45:15+07:00",
		Sgtins:  []string{"04603988000001IE9HALA4IBIH1", "04603988000001JE91ALA4H5J18", "04603988000001KE91ALA517K1J", "046039880000015E9FALA4L95F8"},
	}

	//init mock db
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      db,
		SkipInitializeWithVersion: true,
	})

	conn, err := gorm.Open(dialector, &gorm.Config{})
	_db.DB = conn

	if err != nil {
		t.Fatalf("ошибка '%s' не ожидалась при открытии mock-бд", err)
	}
	defer db.Close()

	productRows := sqlmock.NewRows([]string{"id", "name", "gtin", "packs"}).AddRow(1, "Продукт1", "04603988000001", 4)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT * FROM `boxes` WHERE sscc = ? ORDER BY `boxes`.`id` LIMIT 1").
		WithArgs(request.Sscc).
		WillReturnRows(sqlmock.NewRows([]string{}))
	mock.ExpectQuery("SELECT * FROM `products` WHERE gtin = ? ORDER BY `products`.`id` LIMIT 1").
		WithArgs("04603988000001").
		WillReturnRows(productRows)
	for _, sgtin := range request.Sgtins {
		mock.ExpectQuery("SELECT * FROM `packages` WHERE sgtin = ? ORDER BY `packages`.`id` LIMIT 1").
			WithArgs(sgtin).
			WillReturnRows(sqlmock.NewRows([]string{}))
	}
	mock.ExpectExec("INSERT INTO `boxes` (`sscc`,`created_at`,`updated_at`) VALUES (?,?,?)").
		WithArgs(request.Sscc, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO `packages` (`sgtin`,`box_id`,`product_id`,`created_at`,`updated_at`)"+
		" VALUES (?,?,?,?,?),(?,?,?,?,?),(?,?,?,?,?),(?,?,?,?,?) ON DUPLICATE KEY UPDATE `box_id`=VALUES(`box_id`)").
		WithArgs(
			request.Sgtins[0], 1, 1, sqlmock.AnyArg(), sqlmock.AnyArg(),
			request.Sgtins[1], 1, 1, sqlmock.AnyArg(), sqlmock.AnyArg(),
			request.Sgtins[2], 1, 1, sqlmock.AnyArg(), sqlmock.AnyArg(),
			request.Sgtins[3], 1, 1, sqlmock.AnyArg(), sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(4, 4))
	mock.ExpectCommit()

	response := services.BoxAggregate(context.Background(), request)

	if !response.Ok {
		t.Errorf("ошибка не ожидалась при выполнении аггрегации: %s", response.Error)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("ожидания не совпали: %s", err)
	}
}

// 1. Не может быть нескольких коробок с одинаковым SSCC
func TestBoxAggregationAlreadyUsedSscc(t *testing.T) {
	request := requests.AggregateRequest{
		Sscc:    "000000000000000001",
		Created: "2022-12-01T06:45:15+07:00",
		Sgtins:  []string{"04603988000001IE9HALA4IBIH1", "04603988000001JE91ALA4H5J18", "04603988000001KE91ALA517K1J", "046039880000015E9FALA4L95F8"},
	}

	//init mock db
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      db,
		SkipInitializeWithVersion: true,
	})

	conn, err := gorm.Open(dialector, &gorm.Config{})
	_db.DB = conn

	if err != nil {
		t.Fatalf("ошибка '%s' не ожидалась при открытии mock-бд", err)
	}
	defer db.Close()

	boxRows := sqlmock.NewRows([]string{"id", "sscc"}).AddRow(1, "000000000000000001")

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT * FROM `boxes` WHERE sscc = ? ORDER BY `boxes`.`id` LIMIT 1").
		WithArgs(request.Sscc).
		WillReturnRows(boxRows)
	mock.ExpectRollback()

	response := services.BoxAggregate(context.Background(), request)

	if response.Ok {
		t.Errorf("успех не ожидался при выполнении аггрегации: %s", response.Error)
	}

	if response.Error != exceptions.BoxErrorAlreadyUsedSscc {
		t.Errorf("данная ошибка не ожидалась при выполнении аггрегации: %s", response.Error)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("ожидания не совпали: %s", err)
	}
}

// 5. Нельзя агрегировать пачки с неизвестным GTIN
func TestBoxAggregationUnknownProduct(t *testing.T) {
	request := requests.AggregateRequest{
		Sscc:    "000000000000000001",
		Created: "2022-12-01T06:45:15+07:00",
		Sgtins:  []string{"04603988000001IE9HALA4IBIH1", "04603988000001JE91ALA4H5J18", "04603988000001KE91ALA517K1J", "046039880000015E9FALA4L95F8"},
	}

	//init mock db
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      db,
		SkipInitializeWithVersion: true,
	})

	conn, err := gorm.Open(dialector, &gorm.Config{})
	_db.DB = conn

	if err != nil {
		t.Fatalf("ошибка '%s' не ожидалась при открытии mock-бд", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT * FROM `boxes` WHERE sscc = ? ORDER BY `boxes`.`id` LIMIT 1").
		WithArgs(request.Sscc).
		WillReturnRows(sqlmock.NewRows([]string{}))
	mock.ExpectQuery("SELECT * FROM `products` WHERE gtin = ? ORDER BY `products`.`id` LIMIT 1").
		WithArgs("04603988000001").
		WillReturnRows(sqlmock.NewRows([]string{}))
	mock.ExpectRollback()

	response := services.BoxAggregate(context.Background(), request)

	if response.Ok {
		t.Errorf("успех не ожидался при выполнении аггрегации: %s", response.Error)
	}

	if response.Error != exceptions.BoxErrorUnknownGtin {
		t.Errorf("данная ошибка не ожидалась при выполнении аггрегации: %s", response.Error)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("ожидания не совпали: %s", err)
	}
}

// 2. Одна и та же пачка может быть агрегирована только в одну коробку
func TestBoxAggregationAlreadyUsedPackage(t *testing.T) {
	request := requests.AggregateRequest{
		Sscc:    "000000000000000001",
		Created: "2022-12-01T06:45:15+07:00",
		Sgtins:  []string{"04603988000001IE9HALA4IBIH1", "04603988000001JE91ALA4H5J18", "04603988000001KE91ALA517K1J", "046039880000015E9FALA4L95F8"},
	}

	//init mock db
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      db,
		SkipInitializeWithVersion: true,
	})

	conn, err := gorm.Open(dialector, &gorm.Config{})
	_db.DB = conn

	if err != nil {
		t.Fatalf("ошибка '%s' не ожидалась при открытии mock-бд", err)
	}
	defer db.Close()

	productRows := sqlmock.NewRows([]string{"id", "name", "gtin", "packs"}).AddRow(1, "Продукт1", "04603988000001", 4)
	packagesRows := sqlmock.NewRows([]string{"id", "sgtin"}).AddRow(1, "04603988000001IE9HALA4IBIH1")

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT * FROM `boxes` WHERE sscc = ? ORDER BY `boxes`.`id` LIMIT 1").
		WithArgs(request.Sscc).
		WillReturnRows(sqlmock.NewRows([]string{}))
	mock.ExpectQuery("SELECT * FROM `products` WHERE gtin = ? ORDER BY `products`.`id` LIMIT 1").
		WithArgs("04603988000001").
		WillReturnRows(productRows)
	mock.ExpectQuery("SELECT * FROM `packages` WHERE sgtin = ? ORDER BY `packages`.`id` LIMIT 1").
		WithArgs(request.Sgtins[0]).
		WillReturnRows(packagesRows)
	mock.ExpectRollback()

	response := services.BoxAggregate(context.Background(), request)

	if response.Ok {
		t.Errorf("успех не ожидался при выполнении аггрегации: %s", response.Error)
	}

	if response.Error != exceptions.BoxErrorAlreadyUsedSgtin {
		t.Errorf("данная ошибка не ожидалась при выполнении аггрегации: %s", response.Error)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("ожидания не совпали: %s", err)
	}
}

// 3. В одну коробку можно упаковывать только пачки с одинаковым GTIN
func TestBoxAggregationDifferentGtins(t *testing.T) {
	request := requests.AggregateRequest{
		Sscc:    "000000000000000001",
		Created: "2022-12-01T06:45:15+07:00",
		Sgtins:  []string{"04603988000001IE9HALA4IBIH1", "04603988000002JE91ALA4H5J18", "04603988000001KE91ALA517K1J", "046039880000015E9FALA4L95F8"},
	}

	//init mock db
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      db,
		SkipInitializeWithVersion: true,
	})

	conn, err := gorm.Open(dialector, &gorm.Config{})
	_db.DB = conn

	if err != nil {
		t.Fatalf("ошибка '%s' не ожидалась при открытии mock-бд", err)
	}
	defer db.Close()

	productRows := sqlmock.NewRows([]string{"id", "name", "gtin", "packs"}).AddRow(1, "Продукт1", "04603988000001", 4)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT * FROM `boxes` WHERE sscc = ? ORDER BY `boxes`.`id` LIMIT 1").
		WithArgs(request.Sscc).
		WillReturnRows(sqlmock.NewRows([]string{}))
	mock.ExpectQuery("SELECT * FROM `products` WHERE gtin = ? ORDER BY `products`.`id` LIMIT 1").
		WithArgs("04603988000001").
		WillReturnRows(productRows)
	mock.ExpectQuery("SELECT * FROM `packages` WHERE sgtin = ? ORDER BY `packages`.`id` LIMIT 1").
		WithArgs(request.Sgtins[0]).
		WillReturnRows(sqlmock.NewRows([]string{}))
	mock.ExpectRollback()

	response := services.BoxAggregate(context.Background(), request)

	if response.Ok {
		t.Errorf("успех не ожидался при выполнении аггрегации: %s", response.Error)
	}

	if response.Error != exceptions.BoxErrorDifferentGtins {
		t.Errorf("данная ошибка не ожидалась при выполнении аггрегации: %s", response.Error)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("ожидания не совпали: %s", err)
	}
}

// 4. В одну коробку можно упаковать только N штук пачек, где N задается в справочнике продуктов
func TestBoxAggregationInvalidPackagesAmount(t *testing.T) {
	request := requests.AggregateRequest{
		Sscc:    "000000000000000001",
		Created: "2022-12-01T06:45:15+07:00",
		Sgtins:  []string{"04603988000001IE9HALA4IBIH1", "04603988000002JE91ALA4H5J18", "04603988000001KE91ALA517K1J"},
	}

	//init mock db
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      db,
		SkipInitializeWithVersion: true,
	})

	conn, err := gorm.Open(dialector, &gorm.Config{})
	_db.DB = conn

	if err != nil {
		t.Fatalf("ошибка '%s' не ожидалась при открытии mock-бд", err)
	}
	defer db.Close()

	productRows := sqlmock.NewRows([]string{"id", "name", "gtin", "packs"}).AddRow(1, "Продукт1", "04603988000001", 4)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT * FROM `boxes` WHERE sscc = ? ORDER BY `boxes`.`id` LIMIT 1").
		WithArgs(request.Sscc).
		WillReturnRows(sqlmock.NewRows([]string{}))
	mock.ExpectQuery("SELECT * FROM `products` WHERE gtin = ? ORDER BY `products`.`id` LIMIT 1").
		WithArgs("04603988000001").
		WillReturnRows(productRows)
	mock.ExpectRollback()

	response := services.BoxAggregate(context.Background(), request)

	if response.Ok {
		t.Errorf("успех не ожидался при выполнении аггрегации: %s", response.Error)
	}

	if response.Error != exceptions.BoxErrorInvalidPackagesNumber {
		t.Errorf("данная ошибка не ожидалась при выполнении аггрегации: %s", response.Error)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("ожидания не совпали: %s", err)
	}
}
