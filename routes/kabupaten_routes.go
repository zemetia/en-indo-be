package routes

// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/samber/do"
// 	"github.com/zemetia/en-indo-be/controller"
// )

// func Kabupaten(route *gin.Engine, injector *do.Injector) {
// 	kabupatenController := do.MustInvoke[controller.KabupatenController](injector)

// 	kabupaten := route.Group("/api/kabupaten")
// 	{
// 		kabupaten.GET("", kabupatenController.GetAll)
// 		kabupaten.GET("/:id", kabupatenController.GetByID)
// 		kabupaten.GET("/provinsi/:id", kabupatenController.GetByProvinsiID)
// 	}
// }
