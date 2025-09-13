package provider

import (
	"github.com/samber/do"
	"github.com/zemetia/en-indo-be/constants"
	"github.com/zemetia/en-indo-be/controller"
	"github.com/zemetia/en-indo-be/repository"
	"github.com/zemetia/en-indo-be/service"
	"gorm.io/gorm"
)

func ProvideChurchDependencies(injector *do.Injector) {
	// Church Repository
	do.Provide(injector, func(i *do.Injector) (repository.ChurchRepository, error) {
		db := do.MustInvokeNamed[*gorm.DB](i, constants.DB)
		return repository.NewChurchRepository(db), nil
	})

	// Kabupaten Repository (needed by ChurchService)
	do.Provide(injector, func(i *do.Injector) (repository.KabupatenRepository, error) {
		db := do.MustInvokeNamed[*gorm.DB](i, constants.DB)
		return repository.NewKabupatenRepository(db), nil
	})

	// Provinsi Repository (needed by ChurchService)
	do.Provide(injector, func(i *do.Injector) (repository.ProvinsiRepository, error) {
		db := do.MustInvokeNamed[*gorm.DB](i, constants.DB)
		return repository.NewProvinsiRepository(db), nil
	})

	// Church Service
	do.Provide(injector, func(i *do.Injector) (*service.ChurchService, error) {
		churchRepository := do.MustInvoke[repository.ChurchRepository](i)
		kabupatenRepository := do.MustInvoke[repository.KabupatenRepository](i)
		provinsiRepository := do.MustInvoke[repository.ProvinsiRepository](i)
		return service.NewChurchService(churchRepository, kabupatenRepository, provinsiRepository), nil
	})

	// Church Controller
	do.Provide(injector, func(i *do.Injector) (controller.ChurchController, error) {
		churchService := do.MustInvoke[*service.ChurchService](i)
		return controller.NewChurchController(churchService), nil
	})
}
