package casbin

import (
	"log"
	"os"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/persist"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	rediswatcher "github.com/casbin/redis-watcher/v2"
	"github.com/msyaifudin/pos/internal/database"
)

var Enforcer *casbin.Enforcer
var Watcher persist.Watcher

func InitCasbin() {
	var err error

	log.Println("Casbin: Initializing GORM adapter...")
	// Initialize the GORM adapter
	adapter, err := gormadapter.NewAdapterByDB(database.DB)
	if err != nil {
		log.Fatalf("Casbin: Failed to create Casbin GORM adapter: %v", err)
	}
	log.Println("Casbin: GORM adapter initialized.")

	log.Println("Casbin: Creating enforcer...")
	// Create the enforcer with the GORM adapter and the model
	Enforcer, err = casbin.NewEnforcer("pkg/casbin/model.conf", adapter)
	if err != nil {
		log.Fatalf("Casbin: Failed to create Casbin enforcer: %v", err)
	}
	log.Println("Casbin: Enforcer created.")

	log.Println("Casbin: Initializing Redis watcher...")
	// Initialize the Redis watcher
	// Use the Redis connection string from your configuration
	// For example: "redis://localhost:6379"
	redisWatcherURL := os.Getenv("REDIS_ADDR")
	if redisWatcherURL == "" {
		redisWatcherURL = "localhost:6379"
	}
	Watcher, err = rediswatcher.NewWatcher(redisWatcherURL, rediswatcher.WatcherOptions{})
	if err != nil {
		log.Fatalf("Casbin: Failed to create Redis watcher: %v", err)
	}
	log.Println("Casbin: Redis watcher initialized.")

	// Set the watcher for the enforcer
	Enforcer.SetWatcher(Watcher)

	log.Println("Casbin: Loading policy from DB...")
	// Load the policy from DB.
	if err := Enforcer.LoadPolicy(); err != nil {
		log.Fatalf("Casbin: Failed to load policy from DB: %v", err)
	}
	log.Println("Casbin: Policy loaded from DB.")

	log.Println("Casbin: Enforcer initialized with GORM adapter and Redis watcher")
}

// UpdateCasbinPolicy saves the policy to DB and publishes the change to other instances
func UpdateCasbinPolicy() error {
	log.Println("Casbin: Saving policy to DB...")
	if err := Enforcer.SavePolicy(); err != nil {
		return err
	}
	log.Println("Casbin: Policy saved to DB.")

	if Watcher != nil {
		log.Println("Casbin: Publishing policy update via Redis watcher...")
		return Watcher.Update()
	}
	return nil
}

// CheckPolicy checks if a subject has permission to access an object with an action.
func CheckPolicy(sub, obj, act string) bool {
	log.Printf("Casbin: Checking policy for sub=%s, obj=%s, act=%s", sub, obj, act)
	ok, err := Enforcer.Enforce(sub, obj, act)
	if err != nil {
		log.Printf("Casbin: Error during policy enforcement: %v", err)
		return false
	}
	log.Printf("Casbin: Policy check result: %t", ok)
	return ok
}