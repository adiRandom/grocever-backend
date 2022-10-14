package database

import (
	"dealScraper/lib/helpers"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var db *gorm.DB = nil

const migrationError = "Error while migrating the database"
const cannotGetDb = "Cannot get database"
const cannotConnectToDb = "Cannot connect to database"

func getDbUrl() string {
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	databaseName := os.Getenv("MYSQL_DATABASE")
	options := os.Getenv("MYSQL_OPTIONS")
	return fmt.Sprintf("%s:%s@%s/%s?%s", user, password, host, databaseName, options)
}

func InitDatabase(entities ...interface{}) error {
	dsn := getDbUrl()
	_db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return helpers.Error{Msg: cannotConnectToDb, Reason: err.Error()}
	}

	db = _db
	err = db.AutoMigrate(entities...)
	if err != nil {
		return helpers.Error{Msg: migrationError, Reason: err.Error()}
	}

	return nil
}

func GetDb() (*gorm.DB, error) {
	if db == nil {
		return nil, helpers.Error{Msg: cannotGetDb, Reason: "Database not initialized"}
	}
	return db, nil
}
