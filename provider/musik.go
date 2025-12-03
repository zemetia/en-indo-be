package provider

import (
	"github.com/samber/do"
	"github.com/zemetia/en-indo-be/constants"
	"github.com/zemetia/en-indo-be/controller"
	"github.com/zemetia/en-indo-be/repository"
	"github.com/zemetia/en-indo-be/service"
	"gorm.io/gorm"
)

func ProvideMusikDependencies(injector *do.Injector) {
	db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB)

	// Repositories
	pelayananRepository := repository.NewPelayananRepository(db)
	departmentRepository := repository.NewDepartmentRepository(db)
	userRepository := repository.NewUserRepository(db)
	personRepository := repository.NewPersonRepository(db)

	// Service
	musikService := service.NewMusikService(pelayananRepository, departmentRepository, userRepository, personRepository)

	// Register MusikService in the injector
	do.ProvideNamed(injector, constants.MusikService, func(i *do.Injector) (service.MusikService, error) {
		return musikService, nil
	})

	// Controller
	do.Provide(injector, func(i *do.Injector) (controller.MusikController, error) {
		return controller.NewMusikController(musikService), nil
	})
}
