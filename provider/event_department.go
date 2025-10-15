package provider

import (
	"github.com/samber/do"
	"github.com/zemetia/en-indo-be/constants"
	"github.com/zemetia/en-indo-be/controller"
	"github.com/zemetia/en-indo-be/repository"
	"github.com/zemetia/en-indo-be/service"
	"gorm.io/gorm"
)

func ProvideEventDepartmentDependencies(injector *do.Injector) {
	do.Provide(injector, func(i *do.Injector) (repository.EventDepartmentRepository, error) {
		db := do.MustInvokeNamed[*gorm.DB](i, constants.DB)
		return repository.NewEventDepartmentRepository(db), nil
	})

	do.Provide(injector, func(i *do.Injector) (service.EventDepartmentService, error) {
		eventDepartmentRepo := do.MustInvoke[repository.EventDepartmentRepository](i)
		return service.NewEventDepartmentService(eventDepartmentRepo), nil
	})

	do.Provide(injector, func(i *do.Injector) (controller.EventDepartmentController, error) {
		eventDepartmentService := do.MustInvoke[service.EventDepartmentService](i)
		return controller.NewEventDepartmentController(eventDepartmentService), nil
	})
}
