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
	personMemberRepository := repository.NewLifeGroupPersonMemberRepository(db)
	visitorMemberRepository := repository.NewLifeGroupVisitorMemberRepository(db)
	personRepository := repository.NewPersonRepository(db)
	visitorRepository := repository.NewVisitorRepository(db)

	// Service
	lifeGroupService := service.NewLifeGroupService(lifeGroupRepository)
	personMemberService := service.NewLifeGroupPersonMemberService(personMemberRepository, personRepository, lifeGroupRepository)
	visitorMemberService := service.NewLifeGroupVisitorMemberService(visitorMemberRepository, visitorRepository, lifeGroupRepository)

	// Controller
	do.Provide(injector, func(i *do.Injector) (controller.LifeGroupController, error) {
		return controller.NewLifeGroupController(lifeGroupService), nil
	})

	do.Provide(injector, func(i *do.Injector) (controller.LifeGroupPersonMemberController, error) {
		return controller.NewLifeGroupPersonMemberController(personMemberService), nil
	})

	do.Provide(injector, func(i *do.Injector) (controller.LifeGroupVisitorMemberController, error) {
		return controller.NewLifeGroupVisitorMemberController(visitorMemberService), nil
	})
}
