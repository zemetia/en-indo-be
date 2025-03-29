package routes

// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/samber/do"
// 	"github.com/zemetia/en-indo-be/constants"
// 	"github.com/zemetia/en-indo-be/controller"
// 	"github.com/zemetia/en-indo-be/middleware"
// 	"github.com/zemetia/en-indo-be/service"
// )

// func Department(route *gin.Engine, injector *do.Injector) {
// 	jwtService := do.MustInvokeNamed[service.JWTService](injector, constants.JWTService)
// 	departmentController := do.MustInvoke[controller.DepartemenController](injector)

// 	department := route.Group("/api/department")
// 	{
// 		department.POST("", middleware.Authenticate(jwtService), departmentController.Create)
// 		department.GET("", middleware.Authenticate(jwtService), departmentController.GetAll)
// 		department.GET("/:id", middleware.Authenticate(jwtService), departmentController.GetByID)
// 		department.GET("/church/:id", middleware.Authenticate(jwtService), departmentController.GetByChurchID)
// 		department.PUT("/:id", middleware.Authenticate(jwtService), departmentController.Update)
// 		department.DELETE("/:id", middleware.Authenticate(jwtService), departmentController.Delete)
// 	}
// }
