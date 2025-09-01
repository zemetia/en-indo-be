package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"github.com/zemetia/en-indo-be/constants"
	"github.com/zemetia/en-indo-be/controller"
	"github.com/zemetia/en-indo-be/middleware"
	"github.com/zemetia/en-indo-be/service"
)

func Pelayanan(route *gin.Engine, injector *do.Injector) {
	jwtService := do.MustInvokeNamed[service.JWTService](injector, constants.JWTService)
	pelayananController := do.MustInvoke[controller.PelayananController](injector)

	routes := route.Group("/api/pelayanan")
	routes.Use(middleware.Authenticate(jwtService))
	{
		// User's own pelayanan assignments
		routes.GET("/my", pelayananController.GetMyPelayanan)

		// Get all available pelayanan (for dropdowns, etc.)
		routes.GET("/list", pelayananController.GetAllPelayanan)

		// Admin-only routes for managing assignments
		adminRoutes := routes.Group("")
		// TODO: Add admin middleware when available
		// adminRoutes.Use(middleware.AdminOnly())
		{
			adminRoutes.GET("/assignments", pelayananController.GetAllAssignments)
			adminRoutes.POST("/assign", pelayananController.AssignPelayanan)
			adminRoutes.PUT("/assignments/:id", pelayananController.UpdatePelayananAssignment)
			adminRoutes.DELETE("/assignments/:id", pelayananController.UnassignPelayanan)
			adminRoutes.GET("/assignments/:id", pelayananController.GetAssignmentByID)
		}
	}
}