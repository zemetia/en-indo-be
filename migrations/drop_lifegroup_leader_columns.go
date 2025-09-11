package migrations

import (
	"gorm.io/gorm"
)

// DropLifeGroupLeaderColumns removes leader_id and co_leader_id columns from life_groups table
// This is called after removing these fields from the LifeGroup entity
func DropLifeGroupLeaderColumns(db *gorm.DB) error {
	// Check if life_groups table exists
	if !db.Migrator().HasTable("life_groups") {
		return nil // Table doesn't exist yet, nothing to drop
	}

	// Drop leader_id column if it exists
	if db.Migrator().HasColumn("life_groups", "leader_id") {
		if err := db.Migrator().DropColumn("life_groups", "leader_id"); err != nil {
			return err
		}
	}

	// Drop co_leader_id column if it exists
	if db.Migrator().HasColumn("life_groups", "co_leader_id") {
		if err := db.Migrator().DropColumn("life_groups", "co_leader_id"); err != nil {
			return err
		}
	}

	return nil
}