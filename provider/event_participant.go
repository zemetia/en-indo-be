package provider

import (
	"github.com/samber/do"
	"github.com/zemetia/en-indo-be/constants"
	"github.com/zemetia/en-indo-be/controller"
	"github.com/zemetia/en-indo-be/repository"
	"github.com/zemetia/en-indo-be/service"
	"gorm.io/gorm"
)

func ProvideEventParticipantDependencies(injector *do.Injector) {
	do.Provide(injector, func(i *do.Injector) (repository.EventParticipantRepository, error) {
		db := do.MustInvokeNamed[*gorm.DB](i, constants.DB)
		return repository.NewEventParticipantRepository(db), nil
	})

	do.Provide(injector, func(i *do.Injector) (service.EventParticipantService, error) {
		participantRepo := do.MustInvoke[repository.EventParticipantRepository](i)
		eventRepo := do.MustInvoke[repository.EventRepository](i)
		return service.NewEventParticipantService(participantRepo, eventRepo), nil
	})

	do.Provide(injector, func(i *do.Injector) (*controller.EventParticipantController, error) {
		participantService := do.MustInvoke[service.EventParticipantService](i)
		return controller.NewEventParticipantController(participantService), nil
	})
}
