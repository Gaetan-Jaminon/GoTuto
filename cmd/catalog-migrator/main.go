package main

import (
	"fmt"
	"log"
	"os"
	
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "catalog-migrator",
		Short: "Database migration tool for catalog service",
		Long:  `A CLI tool to manage database migrations for the catalog service.`,
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yaml)")
	rootCmd.PersistentFlags().String("db-host", "localhost", "Database host")
	rootCmd.PersistentFlags().Int("db-port", 5432, "Database port")
	rootCmd.PersistentFlags().String("db-user", "catalog_migrator", "Database user")
	rootCmd.PersistentFlags().String("db-password", "", "Database password")
	rootCmd.PersistentFlags().String("db-name", "gotuto", "Database name")
	rootCmd.PersistentFlags().String("db-sslmode", "disable", "Database SSL mode")
	rootCmd.PersistentFlags().String("migrations-path", "./internal/catalog/migrations", "Path to catalog migrations directory")
	
	// Bind flags to viper
	viper.BindPFlag("database.host", rootCmd.PersistentFlags().Lookup("db-host"))
	viper.BindPFlag("database.port", rootCmd.PersistentFlags().Lookup("db-port"))
	viper.BindPFlag("database.user", rootCmd.PersistentFlags().Lookup("db-user"))
	viper.BindPFlag("database.password", rootCmd.PersistentFlags().Lookup("db-password"))
	viper.BindPFlag("database.name", rootCmd.PersistentFlags().Lookup("db-name"))
	viper.BindPFlag("database.sslmode", rootCmd.PersistentFlags().Lookup("db-sslmode"))
	viper.BindPFlag("migrations.path", rootCmd.PersistentFlags().Lookup("migrations-path"))
	
	// Add commands
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(downCmd)
	rootCmd.AddCommand(versionCmd)
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all pending migrations",
	Long:  `Apply all pending database migrations to bring the catalog schema up to date.`,
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
			fmt.Printf("Catalog schema is already up to date at version %d\n", version)
		} else {
			fmt.Printf("Successfully applied catalog migrations to version %d\n", version)
		}
		
		if dirty {
			fmt.Println("WARNING: Catalog schema is in dirty state")
		}
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback migrations",
	Long:  `Rollback catalog database migrations. Use --steps to specify number of migrations to rollback.`,
	Run: func(cmd *cobra.Command, args []string) {
		steps, _ := cmd.Flags().GetInt("steps")
		
		m, err := getMigrate()
		if err != nil {
			log.Fatal("Failed to initialize migrate:", err)
		}
		defer m.Close()
		
		if steps > 0 {
			if err := m.Steps(-steps); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("Failed to rollback %d catalog migrations: %v", steps, err)
			}
			fmt.Printf("Successfully rolled back %d catalog migrations\n", steps)
		} else {
			// Default: rollback 1 migration
			if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
				log.Fatal("Failed to rollback catalog migration:", err)
			}
			fmt.Println("Successfully rolled back 1 catalog migration")
		}
		
		version, dirty, _ := m.Version()
		if version > 0 {
			fmt.Printf("Current catalog version: %d\n", version)
		} else {
			fmt.Println("Catalog schema has no migrations applied")
		}
		
		if dirty {
			fmt.Println("WARNING: Catalog schema is in dirty state")
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show current migration version",
	Long:  `Display the current migration version and dirty state of the catalog schema.`,
	Run: func(cmd *cobra.Command, args []string) {
		m, err := getMigrate()
		if err != nil {
			log.Fatal("Failed to initialize migrate:", err)
		}
		defer m.Close()
		
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatal("Failed to get catalog version:", err)
		}
		
		fmt.Printf("Current catalog version: %d\n", version)
		if dirty {
			fmt.Println("Status: DIRTY")
		} else {
			fmt.Println("Status: CLEAN")
		}
	},
}

func init() {
	downCmd.Flags().Int("steps", 1, "Number of migrations to rollback")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
	}
	
	// Environment variables with CATALOG prefix
	viper.SetEnvPrefix("CATALOG")
	viper.AutomaticEnv()
	
	// Read config file
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func getMigrate() (*migrate.Migrate, error) {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s&search_path=catalog",
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

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}