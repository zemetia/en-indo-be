package provider

import (
	"github.com/samber/do"
	"github.com/zemetia/en-indo-be/constants"
	"github.com/zemetia/en-indo-be/controller"
	"github.com/zemetia/en-indo-be/repository"
	"github.com/zemetia/en-indo-be/service"
	"gorm.io/gorm"
)

func ProvideLifeGroupDependencies(injector *do.Injector) {
	db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB)

	// Repository
	lifeGroupRepository := repository.NewLifeGroupRepository(db)

	// Service
	lifeGroupService := service.NewLifeGroupService(lifeGroupRepository)

	// Controller
	do.Provide(injector, func(i *do.Injector) (controller.LifeGroupController, error) {
		return controller.NewLifeGroupController(lifeGroupService), nil
	})
}
