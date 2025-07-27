package main

import (
	"fmt"
	"log"
	
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show current migration version",
	Long:  `Display the current migration version and dirty state of the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		m, err := getMigrate()
		if err != nil {
			log.Fatal("Failed to initialize migrate:", err)
		}
		defer m.Close()
		
		version, dirty, err := m.Version()
		if err != nil {
			fmt.Println("No migrations have been applied yet")
			return
		}
		
		fmt.Printf("Current version: %d\n", version)
		if dirty {
			fmt.Println("Status: DIRTY - Manual intervention required!")
			fmt.Println("The database is in an inconsistent state. You may need to:")
			fmt.Println("  1. Fix the issue manually")
			fmt.Println("  2. Use 'force' command to set a specific version")
		} else {
			fmt.Println("Status: Clean")
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}