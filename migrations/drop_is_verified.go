package migrations

import (
	"gorm.io/gorm"
)

// DropIsVerifiedColumn drops the is_verified column from users table
func DropIsVerifiedColumn(db *gorm.DB) error {
	// Check if column exists before trying to drop it
	if db.Migrator().HasColumn(&UserForMigration{}, "is_verified") {
		if err := db.Migrator().DropColumn(&UserForMigration{}, "is_verified"); err != nil {
			return err
		}
	}
	return nil
}

// UserForMigration is a temporary struct for migration purposes
type UserForMigration struct {
	IsVerified bool `gorm:"column:is_verified"`
}

func (UserForMigration) TableName() string {
	return "users"
}