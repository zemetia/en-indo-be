package provider

import (
	"github.com/samber/do"
	"github.com/zemetia/en-indo-be/constants"
	"github.com/zemetia/en-indo-be/controller"
	"github.com/zemetia/en-indo-be/repository"
	"github.com/zemetia/en-indo-be/service"
	"gorm.io/gorm"
)

func ProvideEventDependencies(injector *do.Injector) {
	// Register EventRepository
	do.Provide(injector, func(i *do.Injector) (repository.EventRepository, error) {
		db := do.MustInvokeNamed[*gorm.DB](i, constants.DB)
		return repository.NewEventRepository(db), nil
	})

	// Register EventPICRepository
	do.Provide(injector, func(i *do.Injector) (repository.EventPICRepository, error) {
		db := do.MustInvokeNamed[*gorm.DB](i, constants.DB)
		return repository.NewEventPICRepository(db), nil
	})

	// Register EventService
	do.Provide(injector, func(i *do.Injector) (service.EventService, error) {
		eventRepo := do.MustInvoke[repository.EventRepository](i)
		eventPICRepo := do.MustInvoke[repository.EventPICRepository](i)
		return service.NewEventService(eventRepo, eventPICRepo), nil
	})

	// Register EventPICService
	do.Provide(injector, func(i *do.Injector) (service.EventPICService, error) {
		eventPICRepo := do.MustInvoke[repository.EventPICRepository](i)
		eventRepo := do.MustInvoke[repository.EventRepository](i)
		return service.NewEventPICService(eventPICRepo, eventRepo), nil
	})

	// Register EventController
	do.Provide(injector, func(i *do.Injector) (*controller.EventController, error) {
		eventService := do.MustInvoke[service.EventService](i)
		return controller.NewEventController(eventService), nil
	})

	// Register EventPICController
	do.Provide(injector, func(i *do.Injector) (*controller.EventPICController, error) {
		eventPICService := do.MustInvoke[service.EventPICService](i)
		return controller.NewEventPICController(eventPICService), nil
	})

	// Register EventPICRoleController
	do.Provide(injector, func(i *do.Injector) (*controller.EventPICRoleController, error) {
		eventPICService := do.MustInvoke[service.EventPICService](i)
		return controller.NewEventPICRoleController(eventPICService), nil
	})
}
