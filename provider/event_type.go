package provider

import (
	"github.com/samber/do"
	"github.com/zemetia/en-indo-be/constants"
	"github.com/zemetia/en-indo-be/controller"
	"github.com/zemetia/en-indo-be/repository"
	"github.com/zemetia/en-indo-be/service"
	"gorm.io/gorm"
)

func ProvideEventTypeDependencies(injector *do.Injector) {
	do.ProvideNamed(injector, constants.EventTypeRepository, func(i *do.Injector) (repository.EventTypeRepository, error) {
		db := do.MustInvokeNamed[*gorm.DB](i, constants.DB)
		return repository.NewEventTypeRepository(db), nil
	})

	do.ProvideNamed(injector, constants.EventTypeService, func(i *do.Injector) (*service.EventTypeService, error) {
		repo := do.MustInvokeNamed[repository.EventTypeRepository](i, constants.EventTypeRepository)
		return service.NewEventTypeService(repo), nil
	})

	do.ProvideNamed(injector, constants.EventTypeController, func(i *do.Injector) (controller.EventTypeController, error) {
		svc := do.MustInvokeNamed[*service.EventTypeService](i, constants.EventTypeService)
		return controller.NewEventTypeController(svc), nil
	})
}
