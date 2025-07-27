package main

import (
	"fmt"
	"log"
	"os"
	
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "billing-migrator",
		Short: "Database migration tool for billing service",
		Long:  `A CLI tool to manage database migrations for the billing service.`,
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yaml)")
	rootCmd.PersistentFlags().String("db-host", "localhost", "Database host")
	rootCmd.PersistentFlags().Int("db-port", 5432, "Database port")
	rootCmd.PersistentFlags().String("db-user", "postgres", "Database user")
	rootCmd.PersistentFlags().String("db-password", "", "Database password")
	rootCmd.PersistentFlags().String("db-name", "billing", "Database name")
	rootCmd.PersistentFlags().String("db-sslmode", "disable", "Database SSL mode")
	rootCmd.PersistentFlags().String("migrations-path", "./migrations", "Path to migrations directory")
	
	// Bind flags to viper
	viper.BindPFlag("database.host", rootCmd.PersistentFlags().Lookup("db-host"))
	viper.BindPFlag("database.port", rootCmd.PersistentFlags().Lookup("db-port"))
	viper.BindPFlag("database.user", rootCmd.PersistentFlags().Lookup("db-user"))
	viper.BindPFlag("database.password", rootCmd.PersistentFlags().Lookup("db-password"))
	viper.BindPFlag("database.name", rootCmd.PersistentFlags().Lookup("db-name"))
	viper.BindPFlag("database.sslmode", rootCmd.PersistentFlags().Lookup("db-sslmode"))
	viper.BindPFlag("migrations.path", rootCmd.PersistentFlags().Lookup("migrations-path"))
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
	
	// Environment variables
	viper.SetEnvPrefix("BILLING_MIGRATE")
	viper.AutomaticEnv()
	
	// Read config file
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}