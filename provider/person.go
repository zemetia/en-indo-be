package provider

import (
	"github.com/samber/do"
	"github.com/zemetia/en-indo-be/constants"
	"github.com/zemetia/en-indo-be/controller"
	"github.com/zemetia/en-indo-be/repository"
	"github.com/zemetia/en-indo-be/service"
	"gorm.io/gorm"
)

func ProvidePersonDependencies(injector *do.Injector) {
	db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB)
	// jwtService := do.MustInvokeNamed[service.JWTService](injector, constants.JWTService)

	// Repository
	personRepository := repository.NewPersonRepository(db)
	churchRepository := repository.NewChurchRepository(db)
	kabupatenRepository := repository.NewKabupatenRepository(db)
	lifeGroupRepository := repository.NewLifeGroupRepository(db)

	// Service
	personService := service.NewPersonService(personRepository, churchRepository, kabupatenRepository, lifeGroupRepository)

	// Controller
	do.Provide(injector, func(i *do.Injector) (controller.PersonController, error) {
		return controller.NewPersonController(personService), nil
	})
}
