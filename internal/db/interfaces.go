package db

import "gorm.io/gorm"

// Database defines the interface for database operations
type Database interface {
	InitDB() (*gorm.DB, error)
	AutoMigrateDB(db *gorm.DB, models ...interface{}) error
}

func NewDatabaseService() Database {
	return &GormDB{}
}
