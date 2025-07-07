package main

import (
	"log"
	"os"

	"github.com/msyaifudin/pos/cmd/api"
	"github.com/msyaifudin/pos/cmd/migrate"
	"github.com/msyaifudin/pos/cmd/seed"
	"github.com/msyaifudin/pos/cmd/resetdb"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal(`Usage: go run main.go <command>

Commands:
  api    - Runs the API server
  migrate - Runs database migrations
  seed   - Seeds the database with initial data
  resetdb - Resets the database (drops all tables, migrates, and seeds)`)
	}

	command := os.Args[1]

	switch command {
	case "api":
		api.Run()
	case "migrate":
		migrate.Run()
	case "seed":
		seed.Run()
	case "resetdb":
		resetdb.Run()
	default:
		log.Fatalf("Unknown command: %s", command)
	}
}