/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/OyinloluB/secrets-sync-agent/internal/db"
	"github.com/OyinloluB/secrets-sync-agent/internal/encryption"
	"github.com/spf13/cobra"
)

var retrieveKey string
var retrieveEncryptionKey string

// retrieveCmd represents the retrieve command
var retrieveCmd = &cobra.Command{
	Use:   "retrieve",
	Short: "Retrieve and decrypt a secret",
	Run: func(cmd *cobra.Command, args []string) {
		if retrieveKey == "" || retrieveEncryptionKey == "" {
			fmt.Println("Error: secret key and encryption key must be provided")
			os.Exit(1)
		}

		var encryptedValue string
		var expiresAt sql.NullString

		query := `SELECT value, expires_at FROM secrets WHERE key = ?;`
		err := db.DB.QueryRow(query, retrieveKey).Scan(&encryptedValue, &expiresAt)

		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("Secret not found")
				return
			}
			fmt.Printf("Failed to retrieve secret: %v\n", err)
			os.Exit(1)
		}

		if expiresAt.Valid && expiresAt.String != "" {
			expiryTime, parseErr := time.Parse(time.RFC3339, expiresAt.String)
			if parseErr != nil {
				fmt.Printf("Invalid expiration date format: %v\n", parseErr)
			} else {
				if time.Now().After(expiryTime) {
					fmt.Println("Secret has expired and cannot be retrieved.")
					return
				}
			}
		}

		plaintext, err := encryption.Decrypt(encryptedValue, retrieveEncryptionKey)
		if err != nil {
			fmt.Printf("Decryption failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Secret value: %s\n", plaintext)
	},
}

func init() {
	rootCmd.AddCommand(retrieveCmd)

	retrieveCmd.Flags().StringVarP(&retrieveKey, "key", "k", "", "Key for the secret to retrieve")
	retrieveCmd.Flags().StringVarP(&retrieveEncryptionKey, "encryption-key", "e", "", "Encryption key (must match what was used during storage)")
}
