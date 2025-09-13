package provider

import (
	"github.com/samber/do"
	"github.com/zemetia/en-indo-be/constants"
	"github.com/zemetia/en-indo-be/controller"
	"github.com/zemetia/en-indo-be/repository"
	"github.com/zemetia/en-indo-be/service"
	"gorm.io/gorm"
)

func ProvideUserDependencies(injector *do.Injector) {
	db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB)
	jwtService := do.MustInvokeNamed[service.JWTService](injector, constants.JWTService)
	documentService := do.MustInvokeNamed[service.DocumentService](injector, constants.DocumentService)

	// Repository
	userRepository := repository.NewUserRepository(db)
	personRepository := repository.NewPersonRepository(db)

	// Service
	do.ProvideNamed(injector, constants.UserService, func(i *do.Injector) (service.UserService, error) {
		return service.NewUserService(userRepository, personRepository, documentService, jwtService), nil
	})

	// Controller
	do.Provide(injector, func(i *do.Injector) (controller.UserController, error) {
		userService := do.MustInvokeNamed[service.UserService](i, constants.UserService)
		return controller.NewUserController(userService), nil
	})
}
