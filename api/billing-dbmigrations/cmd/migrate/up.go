package main

import (
	"fmt"
	"log"
	
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all pending migrations",
	Long:  `Apply all pending database migrations to bring the database schema up to date.`,
	Run: func(cmd *cobra.Command, args []string) {
		m, err := getMigrate()
		if err != nil {
			log.Fatal("Failed to initialize migrate:", err)
		}
		defer m.Close()
		
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Failed to apply migrations:", err)
		}
		
		version, dirty, _ := m.Version()
		if err == migrate.ErrNoChange {
			fmt.Printf("Database is already up to date at version %d\n", version)
		} else {
			fmt.Printf("Successfully applied migrations to version %d\n", version)
		}
		
		if dirty {
			fmt.Println("WARNING: Database is in dirty state")
		}
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}

func getMigrate() (*migrate.Migrate, error) {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.name"),
		viper.GetString("database.sslmode"),
	)
	
	migrationsPath := fmt.Sprintf("file://%s", viper.GetString("migrations.path"))
	
	return migrate.New(migrationsPath, dbURL)
}