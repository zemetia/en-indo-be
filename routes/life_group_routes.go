package routes

// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/samber/do"
// 	"github.com/zemetia/en-indo-be/constants"
// 	"github.com/zemetia/en-indo-be/controller"
// 	"github.com/zemetia/en-indo-be/middleware"
// 	"github.com/zemetia/en-indo-be/service"
// )

// func LifeGroup(route *gin.Engine, injector *do.Injector) {
// 	jwtService := do.MustInvokeNamed[service.JWTService](injector, constants.JWTService)
// 	lifeGroupController := do.MustInvoke[controller.LifeGroupController](injector)

// 	lifeGroup := route.Group("/api/life-group")
// 	lifeGroup.Use(middleware.Authenticate(jwtService))
// 	{
// 		lifeGroup.POST("", lifeGroupController.Create)
// 		lifeGroup.GET("", lifeGroupController.GetAll)
// 		lifeGroup.GET("/:id", lifeGroupController.GetByID)
// 		lifeGroup.PUT("/:id", lifeGroupController.Update)
// 		lifeGroup.DELETE("/:id", lifeGroupController.Delete)
// 		lifeGroup.PUT("/:id/leader", lifeGroupController.UpdateLeader)
// 		lifeGroup.PUT("/:id/members", lifeGroupController.UpdateMembers)
// 		lifeGroup.PUT("/:id/persons", lifeGroupController.UpdatePersons)
// 		lifeGroup.GET("/church/:church_id", lifeGroupController.GetByChurchID)

// 		// Manajemen LifeGroup
// 		routes.POST("/:id/life-group/:life_group_id", personController.AddToLifeGroup)
// 		routes.DELETE("/:id/life-group/:life_group_id", personController.RemoveFromLifeGroup)
// 	}
// }
