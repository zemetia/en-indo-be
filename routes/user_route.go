package routes

// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/samber/do"
// 	"github.com/zemetia/en-indo-be/constants"
// 	"github.com/zemetia/en-indo-be/controller"
// 	"github.com/zemetia/en-indo-be/middleware"
// 	"github.com/zemetia/en-indo-be/service"
// )

// func User(route *gin.Engine, injector *do.Injector) {
// 	jwtService := do.MustInvokeNamed[service.JWTService](injector, constants.JWTService)
// 	userController := do.MustInvoke[controller.UserController](injector)

// 	routes := route.Group("/api/user")
// 	{
// 		// User
// 		routes.POST("", userController.Register)
// 		routes.GET("", userController.GetAllUser)
// 		routes.POST("/login", userController.Login)
// 		routes.DELETE("", middleware.Authenticate(jwtService), userController.Delete)
// 		routes.PATCH("", middleware.Authenticate(jwtService), userController.Update)
// 		routes.GET("/me", middleware.Authenticate(jwtService), userController.Me)
// 		routes.POST("/verify_email", userController.VerifyEmail)
// 		routes.POST("/send_verification_email", userController.SendVerificationEmail)
// 	}
// }
