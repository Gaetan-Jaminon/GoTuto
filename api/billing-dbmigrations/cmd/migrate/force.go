package main

import (
	"fmt"
	"log"
	"strconv"
	
	"github.com/spf13/cobra"
)

var forceCmd = &cobra.Command{
	Use:   "force [version]",
	Short: "Force database migration version",
	Long: `Force the database migration version without running migrations. 
This is useful when the database is in a dirty state and needs manual intervention.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		version, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal("Invalid version number:", err)
		}
		
		m, err := getMigrate()
		if err != nil {
			log.Fatal("Failed to initialize migrate:", err)
		}
		defer m.Close()
		
		if err := m.Force(version); err != nil {
			log.Fatal("Failed to force version:", err)
		}
		
		fmt.Printf("Successfully forced database version to: %d\n", version)
		fmt.Println("WARNING: This operation bypasses migrations. Ensure database schema matches this version!")
	},
}

func init() {
	rootCmd.AddCommand(forceCmd)
}