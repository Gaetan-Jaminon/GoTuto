package main

import (
	"fmt"
	"log"
	
	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"
)

var (
	steps int
	all   bool
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback migrations",
	Long:  `Rollback database migrations. Use --steps to specify number of migrations to rollback, or --all to rollback all migrations.`,
	Run: func(cmd *cobra.Command, args []string) {
		m, err := getMigrate()
		if err != nil {
			log.Fatal("Failed to initialize migrate:", err)
		}
		defer m.Close()
		
		if all {
			if err := m.Down(); err != nil && err != migrate.ErrNoChange {
				log.Fatal("Failed to rollback all migrations:", err)
			}
			fmt.Println("Successfully rolled back all migrations")
		} else if steps > 0 {
			if err := m.Steps(-steps); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("Failed to rollback %d migrations: %v", steps, err)
			}
			fmt.Printf("Successfully rolled back %d migrations\n", steps)
		} else {
			// Default: rollback 1 migration
			if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
				log.Fatal("Failed to rollback migration:", err)
			}
			fmt.Println("Successfully rolled back 1 migration")
		}
		
		version, dirty, _ := m.Version()
		if version > 0 {
			fmt.Printf("Current version: %d\n", version)
		} else {
			fmt.Println("Database has no migrations applied")
		}
		
		if dirty {
			fmt.Println("WARNING: Database is in dirty state")
		}
	},
}

func init() {
	downCmd.Flags().IntVar(&steps, "steps", 0, "Number of migrations to rollback")
	downCmd.Flags().BoolVar(&all, "all", false, "Rollback all migrations")
	rootCmd.AddCommand(downCmd)
}