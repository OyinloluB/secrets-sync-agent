package cmd

import (
	"fmt"
	"os"

	"github.com/OyinloluB/secrets-sync-agent/internal/db"
	"github.com/spf13/cobra"
)

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Delete all secrets and reset ID counters",
	Run: func(cmd *cobra.Command, args []string) {
		confirm := ""
		fmt.Println("⚠️  WARNING: This will delete all secrets permanently. Type 'yes' to confirm:")
		fmt.Scanln(&confirm)

		if confirm != "yes" {
			fmt.Println("Reset cancelled.")
			return
		}

		_, err := db.DB.Exec(`DELETE FROM secrets;`)
		if err != nil {
			fmt.Printf("Failed to delete secrets: %v\n", err)
			os.Exit(1)
		}

		_, err = db.DB.Exec(`DELETE FROM sqlite_sequence WHERE name='secrets';`)
		if err != nil {
			fmt.Printf("Failed to reset ID sequence: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("All secrets deleted; agent reset successfully!")
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
