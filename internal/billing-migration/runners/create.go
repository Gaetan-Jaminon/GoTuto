package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createCmd = &cobra.Command{
	Use:   "create [migration_name]",
	Short: "Create a new migration file",
	Long:  `Create a new migration file pair (up and down) with a timestamp prefix.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		
		// Validate migration name
		if !isValidMigrationName(name) {
			log.Fatal("Invalid migration name. Use only letters, numbers, and underscores.")
		}
		
		// Get migrations path
		migrationsPath := viper.GetString("migrations.path")
		if err := os.MkdirAll(migrationsPath, 0755); err != nil {
			log.Fatal("Failed to create migrations directory:", err)
		}
		
		// Generate timestamp
		timestamp := time.Now().Format("20060102150405")
		
		// Create file names
		upFile := filepath.Join(migrationsPath, fmt.Sprintf("%s_%s.up.sql", timestamp, name))
		downFile := filepath.Join(migrationsPath, fmt.Sprintf("%s_%s.down.sql", timestamp, name))
		
		// Create up migration file
		if err := createMigrationFile(upFile, "up", name); err != nil {
			log.Fatal("Failed to create up migration:", err)
		}
		
		// Create down migration file
		if err := createMigrationFile(downFile, "down", name); err != nil {
			log.Fatal("Failed to create down migration:", err)
		}
		
		fmt.Printf("Created migration files:\n")
		fmt.Printf("  - %s\n", upFile)
		fmt.Printf("  - %s\n", downFile)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}

func isValidMigrationName(name string) bool {
	match, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", name)
	return match
}

func createMigrationFile(path, direction, name string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	
	// Write template content
	template := fmt.Sprintf(`-- Migration: %s
-- Direction: %s
-- Created: %s

`, name, direction, time.Now().Format("2006-01-02 15:04:05"))
	
	if direction == "up" {
		template += fmt.Sprintf(`-- TODO: Add your %s migration SQL here
-- Example:
-- CREATE TABLE example (
--     id SERIAL PRIMARY KEY,
--     name VARCHAR(255) NOT NULL,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
-- );
`, strings.ToUpper(direction))
	} else {
		template += fmt.Sprintf(`-- TODO: Add your %s migration SQL here
-- This should undo the changes made in the corresponding up migration
-- Example:
-- DROP TABLE IF EXISTS example;
`, strings.ToUpper(direction))
	}
	
	_, err = file.WriteString(template)
	return err
}