package main

import (
	"github.com/OyinloluB/secrets-sync-agent/cmd"
	"github.com/OyinloluB/secrets-sync-agent/internal/db"
)

func main() {
	db.InitDB("secrets.db")

	cmd.Execute()
}
