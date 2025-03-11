package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

type GormDB struct{}

func (g *GormDB) InitDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("kubero.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return nil, err
	}
	return db, nil
}

func (g *GormDB) AutoMigrateDB(db *gorm.DB, models ...interface{}) error {
	return db.AutoMigrate(models...)
}

type Instance struct {
	gorm.Model
	Name   string
	ApiURL string
}

var DB *gorm.DB

func NewGormDB() Database {
	return &GormDB{}
}
