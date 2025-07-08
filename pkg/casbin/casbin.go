package casbin

import (
	"log"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/msyaifudin/pos/internal/database"
)

var Enforcer *casbin.Enforcer

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

	// Load the policy from DB.
	if err := Enforcer.LoadPolicy(); err != nil {
		log.Fatalf("Failed to load policy from DB: %v", err)
	}

	log.Println("Casbin enforcer initialized with GORM adapter")
}