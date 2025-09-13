package provider

import (
	"github.com/samber/do"
	"github.com/zemetia/en-indo-be/controller"
	"github.com/zemetia/en-indo-be/repository"
	"github.com/zemetia/en-indo-be/service"
)

func ProvideProvinsiDependencies(injector *do.Injector) {
	// Provinsi Service
	do.Provide(injector, func(i *do.Injector) (service.ProvinsiService, error) {
		provinsiRepository := do.MustInvoke[repository.ProvinsiRepository](i)
		return service.NewProvinsiService(provinsiRepository), nil
	})

	// Provinsi Controller
	do.Provide(injector, func(i *do.Injector) (controller.ProvinsiController, error) {
		provinsiService := do.MustInvoke[service.ProvinsiService](i)
		return controller.NewProvinsiController(provinsiService), nil
	})
}
