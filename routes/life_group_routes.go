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
	lifeGroupService := do.MustInvokeNamed[service.LifeGroupService](injector, constants.LifeGroupService)
	lifeGroupController := do.MustInvoke[controller.LifeGroupController](injector)
	personMemberController := do.MustInvoke[controller.LifeGroupPersonMemberController](injector)
	visitorMemberController := do.MustInvoke[controller.LifeGroupVisitorMemberController](injector)

	lifeGroup := route.Group("/api/lifegroup")
	lifeGroup.Use(middleware.Authenticate(jwtService, userService))
	{
		// Open endpoints (no PIC restriction)
		lifeGroup.POST("", lifeGroupController.Create)
		lifeGroup.GET("", lifeGroupController.GetAll)

		// Member and PIC specific endpoints
		lifeGroup.GET("/my-lifegroup", lifeGroupController.GetMyLifeGroup)
		lifeGroup.GET("/daftar", lifeGroupController.GetDaftarLifeGroup)
		lifeGroup.GET("/pic-role", lifeGroupController.GetLifeGroupsByPICRole)

		// Church and user endpoints
		lifeGroup.GET("/church/:church_id", lifeGroupController.GetByChurch)
		lifeGroup.GET("/user/:user_id", lifeGroupController.GetByUser)

		// Batch endpoints
		lifeGroup.POST("/batch/churches", lifeGroupController.GetByMultipleChurches)

		// Edit endpoints (require PIC or leader/co-leader access)
		editGroup := lifeGroup.Group("")
		editGroup.Use(middleware.RequireLifeGroupEditAccess(lifeGroupService))
		{
			editGroup.PUT("/:id", lifeGroupController.Update)
			editGroup.PUT("/:id/leader", lifeGroupController.UpdateLeader)
		}

		// Delete endpoints (require PIC or leader access only)
		deleteGroup := lifeGroup.Group("")
		deleteGroup.Use(middleware.RequireLifeGroupDeleteAccess(lifeGroupService))
		{
			deleteGroup.DELETE("/:id", lifeGroupController.Delete)
		}

		// Member management endpoints (require PIC or leader/co-leader access)
		manageGroup := lifeGroup.Group("")
		manageGroup.Use(middleware.RequireLifeGroupManageAccess(lifeGroupService))
		{

			// Person Member Management
			manageGroup.POST("/:id/person-members", personMemberController.AddPersonMember)
			manageGroup.POST("/:id/person-members/batch", personMemberController.AddPersonMembersBatch)
			manageGroup.PUT("/:id/person-members/position", personMemberController.UpdatePersonMemberPosition)
			manageGroup.DELETE("/:id/person-members", personMemberController.RemovePersonMember)

			// Visitor Member Management
			manageGroup.POST("/:id/visitor-members", visitorMemberController.AddVisitorMember)
			manageGroup.POST("/:id/visitor-members/batch", visitorMemberController.AddVisitorMembersBatch)
			manageGroup.DELETE("/:id/visitor-members", visitorMemberController.RemoveVisitorMember)
		}

		// View-only endpoints for members (require view access - PIC, leader, co-leader, or member)
		viewGroup := lifeGroup.Group("")
		viewGroup.Use(middleware.RequireLifeGroupViewAccess(lifeGroupService))
		{
			viewGroup.GET("/:id", lifeGroupController.GetByID)
			viewGroup.GET("/:id/person-members", personMemberController.GetPersonMembers)
			viewGroup.GET("/:id/leadership-structure", personMemberController.GetLeadershipStructure)
			viewGroup.GET("/:id/visitor-members", visitorMemberController.GetVisitorMembers)
		}

	}
}
