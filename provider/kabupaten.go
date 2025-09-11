package provider

import (
	"github.com/samber/do"
	"github.com/zemetia/en-indo-be/controller"
	"github.com/zemetia/en-indo-be/repository"
	"github.com/zemetia/en-indo-be/service"
)

func ProvideKabupatenDependencies(injector *do.Injector) {
	// Kabupaten Service
	do.Provide(injector, func(i *do.Injector) (service.KabupatenService, error) {
		kabupatenRepository := do.MustInvoke[repository.KabupatenRepository](i)
		return service.NewKabupatenService(kabupatenRepository), nil
	})

	// Kabupaten Controller
	do.Provide(injector, func(i *do.Injector) (controller.KabupatenController, error) {
		kabupatenService := do.MustInvoke[service.KabupatenService](i)
		return controller.NewKabupatenController(kabupatenService), nil
	})
}
