package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"kubero/pkg/kuberoApi"
	"log"
)

func main() {
	// Connect to the database
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Create the database tables
	autoMigrateErr := db.AutoMigrate(&kuberoApi.Metadata{}, &kuberoApi.GitKeys{}, &kuberoApi.GitRepository{}, &kuberoApi.GitWebhook{}, &kuberoApi.Git{}, &kuberoApi.Build{}, &kuberoApi.Fetch{}, &kuberoApi.Run{}, &kuberoApi.Buildpack{}, &kuberoApi.Phase{}, &kuberoApi.PipelineSpec{}, &kuberoApi.PipelineCRD{}, &kuberoApi.AppSpec{}, &kuberoApi.AppCRD{})
	if autoMigrateErr != nil {
		log.Fatal("Failed to create database tables:", autoMigrateErr)
		return
	}
}
