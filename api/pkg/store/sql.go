package store

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func NewSqlStore(driver string, connString string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	switch driver {
	case "sqlite", "sqlite3":
		db, err = gorm.Open(sqlite.Open(connString), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{},
		})
		if err != nil {
			return nil, err
		}
		if err = db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
			return nil, err
		}
	case "mysql":
		db, err = gorm.Open(mysql.Open(connString), &gorm.Config{})
		if err != nil {
			return nil, err
		}
	case "postgres":
		db, err = gorm.Open(postgres.Open(connString), &gorm.Config{})
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid driver provided: %s", driver)
	}

	return db, nil
}
