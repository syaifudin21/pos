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

	// Initialize the GORM adapter
	adapter, err := gormadapter.NewAdapterByDB(database.DB)
	if err != nil {
		log.Fatalf("Failed to create Casbin GORM adapter: %v", err)
	}

	// Create the enforcer with the GORM adapter and the model
	Enforcer, err = casbin.NewEnforcer("pkg/casbin/model.conf", adapter)
	if err != nil {
		log.Fatalf("Failed to create Casbin enforcer: %v", err)
	}

	// Initialize the Redis watcher
	// Use the Redis connection string from your configuration
	// For example: "redis://localhost:6379"
	redisWatcherURL := os.Getenv("REDIS_ADDR")
	if redisWatcherURL == "" {
		redisWatcherURL = "localhost:6379"
	}
	Watcher, err = rediswatcher.NewWatcher(redisWatcherURL, rediswatcher.WatcherOptions{})
	if err != nil {
		log.Fatalf("Failed to create Redis watcher: %v", err)
	}

	// Set the watcher for the enforcer
	Enforcer.SetWatcher(Watcher)

	// Load the policy from DB.
	if err := Enforcer.LoadPolicy(); err != nil {
		log.Fatalf("Failed to load policy from DB: %v", err)
	}

	log.Println("Casbin enforcer initialized with GORM adapter and Redis watcher")
}

// UpdateCasbinPolicy saves the policy to DB and publishes the change to other instances
func UpdateCasbinPolicy() error {
	if err := Enforcer.SavePolicy(); err != nil {
		return err
	}
	if Watcher != nil {
		return Watcher.Update()
	}
	return nil
}
