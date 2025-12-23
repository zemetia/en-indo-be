package migrations

import (
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

func AddOriginalLaguIDToLagu(db *gorm.DB) error {
	if db.Migrator().HasColumn(&entity.Lagu{}, "OriginalLaguID") {
		return nil
	}
	return db.Migrator().AddColumn(&entity.Lagu{}, "OriginalLaguID")
}
