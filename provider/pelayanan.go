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

	// Service
	pelayananService := service.NewPelayananService(pelayananRepository, personRepository, churchRepository)

	// Controller
	do.Provide(injector, func(i *do.Injector) (controller.PelayananController, error) {
		return controller.NewPelayananController(pelayananService), nil
	})
}