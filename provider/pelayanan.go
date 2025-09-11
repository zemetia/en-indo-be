package provider

import (
	"github.com/samber/do"
	"github.com/zemetia/en-indo-be/constants"
	"github.com/zemetia/en-indo-be/controller"
	"github.com/zemetia/en-indo-be/repository"
	"github.com/zemetia/en-indo-be/service"
	"gorm.io/gorm"
)

func ProvidePelayananDependencies(injector *do.Injector) {
	db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB)

	// Repository
	pelayananRepository := repository.NewPelayananRepository(db)
	personRepository := repository.NewPersonRepository(db)
	churchRepository := repository.NewChurchRepository(db)
	departmentRepository := repository.NewDepartmentRepository(db)

	// Service
	do.ProvideNamed(injector, constants.PelayananService, func(i *do.Injector) (service.PelayananService, error) {
		userService := do.MustInvokeNamed[service.UserService](i, constants.UserService)
		return service.NewPelayananService(pelayananRepository, personRepository, churchRepository, departmentRepository, userService), nil
	})

	// Controller
	do.Provide(injector, func(i *do.Injector) (controller.PelayananController, error) {
		pelayananService := do.MustInvokeNamed[service.PelayananService](i, constants.PelayananService)
		return controller.NewPelayananController(pelayananService), nil
	})
}
