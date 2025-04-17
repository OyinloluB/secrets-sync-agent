package cmd

import (
	"fmt"
	"os"

	"github.com/OyinloluB/secrets-sync-agent/internal/db"
	"github.com/OyinloluB/secrets-sync-agent/internal/encryption"
	"github.com/spf13/cobra"
)

var rotateKey string
var newSecretValue string
var rotateEncryptionKey string

// rotateCmd represents the rotate command
var rotateCmd = &cobra.Command{
	Use:   "rotate",
	Short: "Rotate (update) the value of an existing secret",
	Run: func(cmd *cobra.Command, args []string) {
		if rotateKey == "" || newSecretValue == "" || rotateEncryptionKey == "" {
			fmt.Println("Error: secret key, new value, and encryption key must be provided")
			os.Exit(1)
		}

		encryptedValue, err := encryption.Encrypt(newSecretValue, rotateEncryptionKey)
		if err != nil {
			fmt.Printf("Encryption failed: %v\n", err)
			os.Exit(1)
		}

		updateQuery := `UPDATE secrets SET value = ?, created_at = CURRENT_TIMESTAMP WHERE key = ?;`
		result, err := db.DB.Exec(updateQuery, encryptedValue, rotateKey)
		if err != nil {
			fmt.Printf("Failed to update secret: %v\n", err)
			os.Exit(1)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			fmt.Printf("Failed to confirm update: %v\n", err)
			os.Exit(1)
		}

		if rowsAffected == 0 {
			fmt.Println("No secret found with the specified key.")
			return
		}

		fmt.Println("Secret rotated successfully!")
	},
}

func init() {
	rootCmd.AddCommand(rotateCmd)

	rotateCmd.Flags().StringVarP(&rotateKey, "key", "k", "", "Key for the secret to rotate")
	rotateCmd.Flags().StringVarP(&newSecretValue, "new-value", "n", "", "New value for the secret")
	rotateCmd.Flags().StringVarP(&rotateEncryptionKey, "encryption-key", "e", "", "Encryption key (must match original encryption)")
}
