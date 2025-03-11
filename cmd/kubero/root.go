package kubero

import (
	"github.com/faelmori/kubero-cli/internal/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var rootCmd = &cobra.Command{
	Use:   "kubero",
	Short: "Kubero CLI for managing Kubernetes clusters and applications",
	Long: `Kubero is a platform as a service (PaaS) that enables developers to build, run,
and operate applications on Kubernetes.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	initDB()
	rootCmd.AddCommand()
}

func initConfig() {
	viper.AutomaticEnv()
}

func initDB() {
	var err error
	db.DB, err = gorm.Open(sqlite.Open("kubero.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	autoMigrateErr := db.DB.AutoMigrate(&db.Instance{})
	if autoMigrateErr != nil {
		log.Fatal("Failed to migrate database:", autoMigrateErr)
	}
}
