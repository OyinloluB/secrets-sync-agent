/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/OyinloluB/secrets-sync-agent/internal/db"
	"github.com/OyinloluB/secrets-sync-agent/internal/encryption"
	"github.com/spf13/cobra"
)

var secretKey string
var secretValue string
var encryptionKey string

// storeCmd represents the store command
var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "Store a new secret",
	Run: func(cmd *cobra.Command, args []string) {
		if secretKey == "" || secretValue == "" || encryptionKey == "" {
			fmt.Println("Error: secret key, value, and encryption key must be provided")
			os.Exit(1)
		}

		encryptedValue, err := encryption.Encrypt(secretValue, encryptionKey)
		if err != nil {
			fmt.Printf("Encryption failed: %v\n", err)
			os.Exit(1)
		}

		insertQuery := `INSERT INTO secrets (key, value) VALUES (?, ?);`
		_, err = db.DB.Exec(insertQuery, secretKey, encryptedValue)
		if err != nil {
			fmt.Printf("Failed to store secret: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Secret stored successfully!")
	},
}

func init() {
	rootCmd.AddCommand(storeCmd)

	storeCmd.Flags().StringVarP(&secretKey, "key", "k", "", "Key for the secret")
	storeCmd.Flags().StringVarP(&secretValue, "value", "v", "", "Value of the secret")
	storeCmd.Flags().StringVarP(&encryptionKey, "encryption-key", "e", "", "Encryption key (must be 32 characters)")
}
