package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

func RegisterRoutes(server *gin.Engine, injector *do.Injector) {
	User(server, injector)
	Person(server, injector)
	// Church(server, injector)
	// Provinsi(server, injector)
	// Kabupaten(server, injector)
	LifeGroup(server, injector)
	// Department(server, injector)
}
