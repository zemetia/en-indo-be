package migrations

import (
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

func AddTagsToLagu(db *gorm.DB) error {
	if db.Migrator().HasColumn(&entity.Lagu{}, "Tags") {
		return nil
	}
	return db.Migrator().AddColumn(&entity.Lagu{}, "Tags")
}
