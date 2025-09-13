package provider

import (
	"github.com/samber/do"
	"github.com/zemetia/en-indo-be/constants"
	"github.com/zemetia/en-indo-be/controller"
	"github.com/zemetia/en-indo-be/repository"
	"github.com/zemetia/en-indo-be/service"
	"gorm.io/gorm"
)

func ProvideDepartmentDependencies(injector *do.Injector) {
	db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB)

	// Repository
	departmentRepository := repository.NewDepartmentRepository(db)
	pelayananRepository := repository.NewPelayananRepository(db)

	// Service
	departmentService := service.NewDepartmentService(departmentRepository, pelayananRepository)

	// Controller
	do.Provide(injector, func(i *do.Injector) (controller.DepartmentController, error) {
		return controller.NewDepartmentController(departmentService), nil
	})
}
