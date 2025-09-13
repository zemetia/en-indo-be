package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

func RegisterRoutes(server *gin.Engine, injector *do.Injector) {
	// Create API group (without version)
	api := server.Group("/api")

	// Register routes
	User(server, injector)
	Person(server, injector)
	Church(server, injector)
	Provinsi(server, injector)
	Kabupaten(server, injector)
	LifeGroup(server, injector)
	Department(server, injector)
	Pelayanan(server, injector)
	Visitor(server, injector)

	// Register event routes with /api prefix
	EventRoutes(api, injector)
}
