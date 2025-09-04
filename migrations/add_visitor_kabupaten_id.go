package migrations

import (
	"gorm.io/gorm"
)

func AddVisitorKabupatenID(db *gorm.DB) error {
	// Check if visitors table exists
	if !db.Migrator().HasTable("visitors") {
		return nil // Table doesn't exist yet, let AutoMigrate handle it
	}

	// Check if column already exists
	if db.Migrator().HasColumn("visitors", "kabupaten_id") {
		return nil // Column already exists
	}

	// Check if we have any kabupatens
	var kabupatenCount int64
	if err := db.Table("kabupatens").Count(&kabupatenCount).Error; err != nil {
		return err
	}

	if kabupatenCount == 0 {
		// No kabupatens exist, we can't add the constraint
		// Just add the column as nullable for now
		if err := db.Exec("ALTER TABLE visitors ADD COLUMN kabupaten_id INT NULL").Error; err != nil {
			return err
		}
		return nil
	}

	// Get the first available kabupaten_id to set as default
	var firstKabupatenID uint
	if err := db.Table("kabupatens").Select("id").Order("id ASC").Limit(1).Scan(&firstKabupatenID).Error; err != nil {
		return err
	}

	// Add kabupaten_id column as nullable first
	if err := db.Exec("ALTER TABLE visitors ADD COLUMN kabupaten_id INT NULL").Error; err != nil {
		return err
	}

	// Set default value for existing records
	if err := db.Exec("UPDATE visitors SET kabupaten_id = ? WHERE kabupaten_id IS NULL", firstKabupatenID).Error; err != nil {
		return err
	}

	return nil
}