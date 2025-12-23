package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"github.com/zemetia/en-indo-be/constants"
	"github.com/zemetia/en-indo-be/controller"
	"github.com/zemetia/en-indo-be/middleware"
	"github.com/zemetia/en-indo-be/service"
)

func Musik(route *gin.Engine, injector *do.Injector) {
	jwtService := do.MustInvokeNamed[service.JWTService](injector, constants.JWTService)
	userService := do.MustInvokeNamed[service.UserService](injector, constants.UserService)
	musikController := do.MustInvoke[controller.MusikController](injector)

	musik := route.Group("/api/musik")
	musik.Use(middleware.Authenticate(jwtService, userService))
	{
		// Pelayan (Musician) Management
		musik.GET("/pelayan", musikController.GetPelayanMusik)
		musik.GET("/pelayan/:person_id", musikController.GetPelayanMusikByID)
		musik.POST("/pelayan/:person_id/assign", musikController.AssignPelayanan)
		musik.DELETE("/pelayan/assignment/:assignment_id", musikController.RemovePelayanan)
		musik.PUT("/pelayan/assignment/:assignment_id/toggle", musikController.ToggleActive)

		// Available Pelayanan Roles
		musik.GET("/pelayanan-roles", musikController.GetAvailablePelayanan)

		// Available People (congregation members not in music yet)
		musik.GET("/available-people", musikController.GetAvailablePeople)

		// Lagu Management
		laguController := do.MustInvokeNamed[controller.LaguController](injector, constants.LaguController)
		musik.GET("/lagu", laguController.FindAll)
		musik.GET("/lagu/:id", laguController.FindByID)
		musik.POST("/lagu", laguController.Create)
		musik.PUT("/lagu/:id", laguController.Update)
		musik.DELETE("/lagu/:id", laguController.Delete)
	}
}
