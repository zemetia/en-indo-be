package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"github.com/zemetia/en-indo-be/controller"
)

func Visitor(server *gin.Engine, injector *do.Injector) {
	visitorController := do.MustInvoke[controller.VisitorController](injector)
	visitorInfoController := do.MustInvoke[controller.VisitorInformationController](injector)

	// Visitor routes
	visitorRoutes := server.Group("/api/visitor")
	{
		visitorRoutes.POST("", visitorController.Create)
		visitorRoutes.GET("", visitorController.GetAll)
		visitorRoutes.GET("/search", visitorController.Search)
		visitorRoutes.GET("/:id", visitorController.GetByID)
		visitorRoutes.PUT("/:id", visitorController.Update)
		visitorRoutes.DELETE("/:id", visitorController.Delete)
	}

	// Visitor Information routes - separate to avoid route conflicts
	visitorInfoRoutes := server.Group("/api/visitor-information")
	{
		visitorInfoRoutes.GET("", visitorInfoController.GetAll)
		visitorInfoRoutes.GET("/:id", visitorInfoController.GetByID)
		visitorInfoRoutes.PUT("/:id", visitorInfoController.Update)
		visitorInfoRoutes.DELETE("/:id", visitorInfoController.Delete)
		visitorInfoRoutes.GET("/visitor/:visitor_id", visitorInfoController.GetByVisitorID)
		visitorInfoRoutes.POST("/visitor/:visitor_id", visitorInfoController.Create)
	}
}