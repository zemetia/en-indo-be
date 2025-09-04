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
	userService := do.MustInvokeNamed[service.UserService](injector, constants.UserService)
	lifeGroupController := do.MustInvoke[controller.LifeGroupController](injector)
	personMemberController := do.MustInvoke[controller.LifeGroupPersonMemberController](injector)
	visitorMemberController := do.MustInvoke[controller.LifeGroupVisitorMemberController](injector)

	lifeGroup := route.Group("/api/lifegroup")
	lifeGroup.Use(middleware.Authenticate(jwtService, userService))
	{
		lifeGroup.POST("", lifeGroupController.Create)
		lifeGroup.GET("", lifeGroupController.GetAll)
		lifeGroup.GET("/:id", lifeGroupController.GetByID)
		lifeGroup.PUT("/:id", lifeGroupController.Update)
		lifeGroup.DELETE("/:id", lifeGroupController.Delete)
		lifeGroup.PUT("/:id/leader", lifeGroupController.UpdateLeader)

		// Church and user endpoints
		lifeGroup.GET("/church/:church_id", lifeGroupController.GetByChurch)
		lifeGroup.GET("/user/:user_id", lifeGroupController.GetByUser)

		// Batch endpoints
		lifeGroup.POST("/batch/churches", lifeGroupController.GetByMultipleChurches)

		// Person Member Management
		lifeGroup.POST("/:id/person-members", personMemberController.AddPersonMember)
		lifeGroup.GET("/:id/person-members", personMemberController.GetPersonMembers)
		lifeGroup.PUT("/:id/person-members/position", personMemberController.UpdatePersonMemberPosition)
		lifeGroup.DELETE("/:id/person-members", personMemberController.RemovePersonMember)
		lifeGroup.GET("/:id/leadership-structure", personMemberController.GetLeadershipStructure)
		
		// Visitor Member Management
		lifeGroup.POST("/:id/visitor-members", visitorMemberController.AddVisitorMember)
		lifeGroup.GET("/:id/visitor-members", visitorMemberController.GetVisitorMembers)
		lifeGroup.DELETE("/:id/visitor-members", visitorMemberController.RemoveVisitorMember)

	}
}
