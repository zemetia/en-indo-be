package provider

import (
	"github.com/samber/do"
	"github.com/zemetia/en-indo-be/constants"
	"github.com/zemetia/en-indo-be/controller"
	"github.com/zemetia/en-indo-be/repository"
	"github.com/zemetia/en-indo-be/service"
	"gorm.io/gorm"
)

func ProvideVisitorDependencies(injector *do.Injector) {
	db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB)

	// Repository
	visitorRepository := repository.NewVisitorRepository(db)
	visitorInfoRepository := repository.NewVisitorInformationRepository(db)

	// Service
	visitorService := service.NewVisitorService(visitorRepository)
	visitorInfoService := service.NewVisitorInformationService(visitorInfoRepository, visitorRepository)

	// Controller
	do.Provide(injector, func(i *do.Injector) (controller.VisitorController, error) {
		return controller.NewVisitorController(visitorService), nil
	})

	do.Provide(injector, func(i *do.Injector) (controller.VisitorInformationController, error) {
		return controller.NewVisitorInformationController(visitorInfoService), nil
	})
}
