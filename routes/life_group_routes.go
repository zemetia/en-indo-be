package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zemetia/en-indo-be/controller"
	"github.com/zemetia/en-indo-be/middleware"
	"github.com/zemetia/en-indo-be/service"
)

func LifeGroup(r *gin.Engine, lifeGroupService *service.LifeGroupService) {
	lifeGroupController := controller.NewLifeGroupController(lifeGroupService)

	lifeGroup := r.Group("/api/life-group")
	lifeGroup.Use(middleware.AuthMiddleware())
	{
		lifeGroup.POST("", lifeGroupController.Create)
		lifeGroup.GET("", lifeGroupController.GetAll)
		lifeGroup.GET("/:id", lifeGroupController.GetByID)
		lifeGroup.PUT("/:id", lifeGroupController.Update)
		lifeGroup.DELETE("/:id", lifeGroupController.Delete)
		lifeGroup.PUT("/:id/leader", lifeGroupController.UpdateLeader)
		lifeGroup.PUT("/:id/members", lifeGroupController.UpdateMembers)
		lifeGroup.PUT("/:id/persons", lifeGroupController.UpdatePersons)
		lifeGroup.GET("/church/:church_id", lifeGroupController.GetByChurchID)
	}
}
