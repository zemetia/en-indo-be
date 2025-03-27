package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zemetia/en-indo-be/controller"
	"github.com/zemetia/en-indo-be/middleware"
)

func Role(router *gin.Engine, roleController *controller.RoleController) {
	role := router.Group("/api/role")
	{
		role.POST("", middleware.Authenticate(jwtService), roleController.Create)
		role.GET("", middleware.Authenticate(jwtService), roleController.GetAll)
		role.GET("/:id", middleware.Authenticate(jwtService), roleController.GetByID)
		role.PUT("/:id", middleware.Authenticate(jwtService), roleController.Update)
		role.DELETE("/:id", middleware.Authenticate(jwtService), roleController.Delete)
		role.POST("/:id/permissions", middleware.Authenticate(jwtService), roleController.AddPermissions)
		role.DELETE("/:id/permissions", middleware.Authenticate(jwtService), roleController.RemovePermissions)
		role.POST("/user/:id/assign", middleware.Authenticate(jwtService), roleController.AssignToUser)
		role.POST("/user/:id/remove", middleware.Authenticate(jwtService), roleController.RemoveFromUser)
	}
}
