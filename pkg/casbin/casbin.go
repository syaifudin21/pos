package casbin

import (
	"log"

	"github.com/casbin/casbin/v2"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

var Enforcer *casbin.Enforcer

func InitCasbin() {
	var err error
	// Use a file adapter for now. Later, we can switch to a GORM adapter.
	a := fileadapter.NewAdapter("pkg/casbin/policy.csv")
	Enforcer, err = casbin.NewEnforcer("pkg/casbin/model.conf", a)
	if err != nil {
		log.Fatalf("Failed to create Casbin enforcer: %v", err)
	}

	log.Println("Casbin enforcer initialized")
}
