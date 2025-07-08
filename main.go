package main

import (
	"log"
	"os"

	"github.com/msyaifudin/pos/cmd/api"
	"github.com/msyaifudin/pos/cmd/migrate"
	"github.com/msyaifudin/pos/cmd/seed"
	"github.com/msyaifudin/pos/cmd/resetdb"
	"github.com/msyaifudin/pos/cmd/resetcasbin"
	"github.com/msyaifudin/pos/pkg/casbin"
	"github.com/msyaifudin/pos/internal/database"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal(`Usage: go run main.go <command>

Commands:
  api    - Runs the API server
  migrate - Runs database migrations
  seed   - Seeds the database with initial data
  resetdb - Resets the database (drops all tables, migrates, and seeds)
  resetcasbin - Resets Casbin policies to default from policy-default.csv
  checkcasbinpolicy - Checks a specific Casbin policy (e.g., owner, user_payments, read)`)
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
	case "resetcasbin":
		resetcasbin.Run()
	case "checkcasbinpolicy":
		// Initialize database and Casbin before checking policy
		database.InitDB()
		casbin.InitCasbin()
		// Force reload policy from DB to ensure latest policies are used
		if err := casbin.Enforcer.LoadPolicy(); err != nil {
			log.Fatalf("Failed to load policy for check: %v", err)
		}
		log.Printf("Policy check for owner, user_payments, read: %t", casbin.CheckPolicy("owner", "user_payments", "read"))
	default:
		log.Fatalf("Unknown command: %s", command)
	}
}