package provider

import (
	"github.com/samber/do"
	"github.com/zemetia/en-indo-be/config"
	"github.com/zemetia/en-indo-be/constants"
	"github.com/zemetia/en-indo-be/service"
	"gorm.io/gorm"
)

func InitDatabase(injector *do.Injector) {
	do.ProvideNamed(injector, constants.DB, func(i *do.Injector) (*gorm.DB, error) {
		return config.SetUpDatabaseConnection(), nil
	})
}

func RegisterDependencies(injector *do.Injector) {
	InitDatabase(injector)

	do.ProvideNamed(injector, constants.JWTService, func(i *do.Injector) (service.JWTService, error) {
		return service.NewJWTService(), nil
	})

	do.ProvideNamed(injector, constants.DocumentService, func(i *do.Injector) (service.DocumentService, error) {
		return service.NewDocumentService(constants.BASE_URL, constants.UPLOAD_PATH), nil
	})

	ProvideUserDependencies(injector)
	ProvideLifeGroupDependencies(injector)
	ProvideDepartmentDependencies(injector)
	ProvideChurchDependencies(injector)
	ProvideProvinsiDependencies(injector)
	ProvideKabupatenDependencies(injector)
	ProvidePelayananDependencies(injector)
	ProvidePersonDependencies(injector)
	ProvideVisitorDependencies(injector)
}
