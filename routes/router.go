package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

func RegisterRoutes(server *gin.Engine, injector *do.Injector) {
	// Create API v1 group
	v1 := server.Group("/api/v1")
	
	// Register routes
	User(server, injector)
	Person(server, injector)
	Church(server, injector)
	// Provinsi(server, injector)
	// Kabupaten(server, injector)
	LifeGroup(server, injector)
	// Department(server, injector)
	Pelayanan(server, injector)
	
	// Register event routes
	EventRoutes(v1, injector)
}
