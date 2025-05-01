package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"github.com/zemetia/en-indo-be/constants"
	"github.com/zemetia/en-indo-be/controller"
	"github.com/zemetia/en-indo-be/middleware"
	"github.com/zemetia/en-indo-be/service"
)

// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/samber/do"
// 	"github.com/zemetia/en-indo-be/constants"
// 	"github.com/zemetia/en-indo-be/controller"
// 	"github.com/zemetia/en-indo-be/middleware"
// 	"github.com/zemetia/en-indo-be/service"
// )

func LifeGroup(route *gin.Engine, injector *do.Injector) {
	jwtService := do.MustInvokeNamed[service.JWTService](injector, constants.JWTService)
	lifeGroupController := do.MustInvoke[controller.LifeGroupController](injector)

	lifeGroup := route.Group("/api/lifegroup")
	lifeGroup.Use(middleware.Authenticate(jwtService))
	{
		lifeGroup.POST("", lifeGroupController.Create)
		lifeGroup.GET("", lifeGroupController.GetAll)
		lifeGroup.GET("/:id", lifeGroupController.GetByID)
		lifeGroup.PUT("/:id", lifeGroupController.Update)
		lifeGroup.DELETE("/:id", lifeGroupController.Delete)
		lifeGroup.PUT("/:id/leader", lifeGroupController.UpdateLeader)
		lifeGroup.PUT("/:id/members", lifeGroupController.UpdateMembers)

		// Manajemen LifeGroup
		lifeGroup.POST("/:id/life-group/:life_group_id", lifeGroupController.AddToLifeGroup)
		lifeGroup.DELETE("/:id/life-group/:life_group_id", lifeGroupController.RemoveFromLifeGroup)
	}
}
