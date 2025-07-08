package resetcasbin

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/msyaifudin/pos/internal/database"
	"github.com/msyaifudin/pos/pkg/casbin"
)

func Run() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database
	database.InitDB()

	log.Println("Starting Casbin policy reset...")

	// Initialize Casbin enforcer with GORM adapter
	casbin.InitCasbin()

	// Clear existing policies in DB to ensure a clean slate
	log.Println("Clearing all existing Casbin policies from DB...")
	casbin.Enforcer.ClearPolicy()

	// Read policies from policy-default.csv
	file, err := os.Open("pkg/casbin/policy-default.csv")
	if err != nil {
		log.Fatalf("Failed to open pkg/casbin/policy-default.csv: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		log.Printf("Processing policy line: %s", line) // Added logging
		if line == "" || strings.HasPrefix(line, "#") {
			continue // Skip empty lines and comments
		}

		parts := strings.Split(line, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}

		policyType := parts[0]
		policyArgs := parts[1:]

		switch policyType {
		case "p":
			// 'p' policies expect 3 arguments: sub, obj, act
			if len(policyArgs) != 3 {
				log.Printf("Skipping invalid 'p' policy line (expected 3 arguments, got %d): %s", len(policyArgs), line)
				continue
			}
			args := make([]interface{}, len(policyArgs))
			for i, v := range policyArgs {
				args[i] = v
			}
			if _, err := casbin.Enforcer.AddPolicy(args...); err != nil {
				log.Fatalf("Failed to add policy %v: %v", policyArgs, err)
			} else {
				log.Printf("Successfully added policy: %v", policyArgs) // Added logging
			}
		case "g":
			// 'g' policies expect 2 arguments: user, role
			if len(policyArgs) != 2 {
				log.Printf("Skipping invalid 'g' policy line (expected 2 arguments, got %d): %s", len(policyArgs), line)
				continue
			}
			args := make([]interface{}, len(policyArgs))
			for i, v := range policyArgs {
				args[i] = v
			}
			if _, err := casbin.Enforcer.AddGroupingPolicy(args...); err != nil {
				log.Fatalf("Failed to add grouping policy %v: %v", policyArgs, err)
			} else {
				log.Printf("Successfully added grouping policy: %v", policyArgs) // Added logging
			}
		default:
			log.Printf("Unknown policy type: %s in line: %s", policyType, line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading policy-default.csv: %v", err)
	}

	log.Println("Casbin policies loaded from policy-default.csv and saved to database.")

	// Update Casbin policy to notify other instances via Redis watcher
	if err := casbin.UpdateCasbinPolicy(); err != nil {
		log.Fatalf("Failed to update Casbin policy via Redis watcher: %v", err)
	}

	log.Println("Casbin policy reset completed.")
}
